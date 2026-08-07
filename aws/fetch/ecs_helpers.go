package awsfetch

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/wallix/awless/fetch"
)

func getClusterArns(ctx context.Context, cache fetch.Cache, api *ecs.Client) ([]string, error) {
	var arns []string
	if clusterName, hasFilter := getUserFiltersFromContext(ctx)["cluster"]; hasFilter {
		out, err := api.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: []string{clusterName}})
		if err != nil {
			return arns, err
		}
		for _, c := range out.Clusters {
			arns = append(arns, awssdk.ToString(c.ClusterArn))
		}
	} else {
		if val, cerr := cache.Get("getClustersNames", func() (interface{}, error) {
			paginator := ecs.NewListClustersPaginator(api, &ecs.ListClustersInput{})
			for paginator.HasMorePages() {
				out, err := paginator.NextPage(ctx)
				if err != nil {
					return arns, err
				}
				arns = append(arns, out.ClusterArns...)
			}
			return arns, nil
		}); cerr != nil {
			return arns, cerr
		} else if v, ok := val.([]string); ok {
			arns = v
		}
	}
	return arns, nil
}

func getAllTasks(ctx context.Context, cache fetch.Cache, api *ecs.Client) (res []ecstypes.Task, err error) {
	clusterArns, cerr := getClusterArns(ctx, cache, api)
	if cerr != nil {
		return res, cerr
	}

	// Two bounded phases rather than a streaming pipeline.
	//
	// The previous implementation had three defects: it called tasksWG.Add from
	// inside a worker while another goroutine could already be in
	// tasksWG.Wait() (undefined behavior), its consumer returned on the first
	// error leaving senders blocked on unbuffered channels forever, and it
	// spawned two goroutines per cluster plus one per page with no limit.
	//
	// Phase 1 lists task ARNs, phase 2 describes them. Batches are kept
	// per-page because DescribeTasks accepts at most 100 tasks per call.
	type taskBatch struct {
		cluster string
		arns    []string
	}

	var (
		mu      sync.Mutex
		batches []taskBatch
	)

	listG := new(errgroup.Group)
	listG.SetLimit(maxParallelAWSCalls)

	for _, cluster := range clusterArns {
		for _, status := range []ecstypes.DesiredStatus{ecstypes.DesiredStatusRunning, ecstypes.DesiredStatusStopped} {
			cl, st := cluster, status
			listG.Go(func() error {
				paginator := ecs.NewListTasksPaginator(api, &ecs.ListTasksInput{Cluster: &cl, DesiredStatus: st})
				for paginator.HasMorePages() {
					out, er := paginator.NextPage(ctx)
					if er != nil {
						return er
					}
					if len(out.TaskArns) == 0 {
						continue
					}
					mu.Lock()
					batches = append(batches, taskBatch{cluster: cl, arns: out.TaskArns})
					mu.Unlock()
				}
				return nil
			})
		}
	}

	if err = listG.Wait(); err != nil {
		return res, err
	}

	describeG := new(errgroup.Group)
	describeG.SetLimit(maxParallelAWSCalls)

	for _, batch := range batches {
		b := batch
		describeG.Go(func() error {
			out, er := api.DescribeTasks(ctx, &ecs.DescribeTasksInput{Cluster: &b.cluster, Tasks: b.arns})
			if er != nil {
				return er
			}
			mu.Lock()
			res = append(res, out.Tasks...)
			mu.Unlock()
			return nil
		})
	}

	if err = describeG.Wait(); err != nil {
		return res, err
	}

	return res, nil
}
