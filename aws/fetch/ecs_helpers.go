package awsfetch

import (
	"context"
	"sync"

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

	type listTasksOutput struct {
		err     error
		output  *ecs.ListTasksOutput
		cluster string
	}
	tasksNamesc := make(chan listTasksOutput)
	var wg sync.WaitGroup

	for _, cluster := range clusterArns {
		wg.Add(1)
		go func(cl string) {
			defer wg.Done()
			paginator := ecs.NewListTasksPaginator(api, &ecs.ListTasksInput{Cluster: &cl, DesiredStatus: ecstypes.DesiredStatusRunning})
			for paginator.HasMorePages() {
				out, er := paginator.NextPage(ctx)
				if er != nil {
					tasksNamesc <- listTasksOutput{err: er}
					return
				}
				tasksNamesc <- listTasksOutput{output: out, cluster: cl}
			}
		}(cluster)

		wg.Add(1)
		go func(cl string) {
			defer wg.Done()
			paginator := ecs.NewListTasksPaginator(api, &ecs.ListTasksInput{Cluster: &cl, DesiredStatus: ecstypes.DesiredStatusStopped})
			for paginator.HasMorePages() {
				out, er := paginator.NextPage(ctx)
				if er != nil {
					tasksNamesc <- listTasksOutput{err: er}
					return
				}
				tasksNamesc <- listTasksOutput{output: out, cluster: cl}
			}
		}(cluster)
	}

	type describeTasksOutput struct {
		err    error
		output *ecs.DescribeTasksOutput
	}

	tasksc := make(chan describeTasksOutput)
	var tasksWG sync.WaitGroup

	tasksWG.Add(1)
	go func() {
		defer tasksWG.Done()
		for r := range tasksNamesc {
			if r.err != nil {
				tasksc <- describeTasksOutput{err: r.err}
				return
			}
			if len(r.output.TaskArns) == 0 {
				continue
			}

			tasksWG.Add(1)
			go func(taskArns []string, cluster string) {
				defer tasksWG.Done()
				tasksOut, er := api.DescribeTasks(ctx, &ecs.DescribeTasksInput{Cluster: &cluster, Tasks: taskArns})
				tasksc <- describeTasksOutput{err: er, output: tasksOut}
			}(r.output.TaskArns, r.cluster)
		}
	}()

	go func() {
		wg.Wait()
		close(tasksNamesc)
		tasksWG.Wait()
		close(tasksc)
	}()

	for r := range tasksc {
		if err = r.err; err != nil {
			return
		}
		res = append(res, r.output.Tasks...)
	}

	return
}
