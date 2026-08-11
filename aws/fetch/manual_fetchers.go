package awsfetch

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	configservicetypes "github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafv2types "github.com/aws/aws-sdk-go-v2/service/wafv2/types"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	apigatewayv2types "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go"
	"golang.org/x/sync/errgroup"

	awsconv "github.com/bootswithdefer/awless/aws/conv"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/cloud/rdf"
	"github.com/bootswithdefer/awless/fetch"
	"github.com/bootswithdefer/awless/graph"
)

func addManualInfraFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["containerinstance"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []ecstypes.ContainerInstance
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.containerinstance.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[containerinstance]")
			return resources, objects, nil
		}

		clusterArns, err := getClusterArns(ctx, cache, conf.APIs.ECS)
		if err != nil {
			return resources, objects, err
		}

		for _, cluster := range clusterArns {
			paginator := ecs.NewListContainerInstancesPaginator(conf.APIs.ECS, &ecs.ListContainerInstancesInput{Cluster: &cluster})
			for paginator.HasMorePages() {
				out, err := paginator.NextPage(ctx)
				if err != nil {
					return resources, objects, err
				}
				if len(out.ContainerInstanceArns) == 0 {
					continue
				}

				containerInstancesOut, err := conf.APIs.ECS.DescribeContainerInstances(ctx, &ecs.DescribeContainerInstancesInput{Cluster: &cluster, ContainerInstances: out.ContainerInstanceArns})
				if err != nil {
					return resources, objects, err
				}

				for _, inst := range containerInstancesOut.ContainerInstances {
					objects = append(objects, inst)
					var res *graph.Resource
					if res, err = awsconv.NewResource(inst); err != nil {
						return resources, objects, err
					}
					res.Properties()[properties.Cluster] = cluster
					resources = append(resources, res)
					parent := graph.InitResource(cloud.ContainerCluster, cluster)
					res.AddRelation(rdf.ChildrenOfRel, parent)
				}
			}
		}
		return resources, objects, nil
	}

	funcs["container"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []ecstypes.Container
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.container.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[container]")
			return resources, objects, nil
		}

		var tasks []ecstypes.Task

		if val, e := cache.Get("getAllTasks", func() (any, error) {
			return getAllTasks(ctx, cache, conf.APIs.ECS)
		}); e != nil {
			return resources, objects, e
		} else if v, ok := val.([]ecstypes.Task); ok {
			tasks = v
		}

		for _, task := range tasks {
			for _, container := range task.Containers {
				objects = append(objects, container)
				res, err := awsconv.NewResource(container)
				if err != nil {
					return nil, nil, err
				}
				if task.ClusterArn != nil {
					res.Properties()[properties.Cluster] = awssdk.ToString(task.ClusterArn)
				}
				if task.ContainerInstanceArn != nil {
					res.Properties()[properties.ContainerInstance] = awssdk.ToString(task.ContainerInstanceArn)
				}
				if task.CreatedAt != nil {
					res.Properties()[properties.Created] = *task.CreatedAt
				}
				if task.StartedAt != nil {
					res.Properties()[properties.Launched] = *task.StartedAt
				}
				if task.StoppedAt != nil {
					res.Properties()[properties.Stopped] = *task.StoppedAt
				}
				if task.TaskDefinitionArn != nil {
					res.Properties()[properties.ContainerTask] = awssdk.ToString(task.TaskDefinitionArn)
				}
				if task.Group != nil {
					res.Properties()[properties.DeploymentName] = awssdk.ToString(task.Group)
				}

				res.AddRelation(rdf.ChildrenOfRel, graph.InitResource(cloud.ContainerCluster, awssdk.ToString(task.ClusterArn)))
				res.AddRelation(rdf.DependingOnRel, graph.InitResource(cloud.ContainerTask, awssdk.ToString(task.TaskDefinitionArn)))
				res.AddRelation(rdf.DependingOnRel, graph.InitResource(cloud.ContainerInstance, awssdk.ToString(task.ContainerInstanceArn)))

				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["containertask"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []ecstypes.TaskDefinition
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.containertask.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[containertask]")
			return resources, objects, nil
		}

		type resStruct struct {
			res *ecstypes.TaskDefinition
			err error
		}

		var arns []string

		fetchDefinitionsInput := &ecs.ListTaskDefinitionsInput{}
		if givenFamilyPrefix, hasFilter := getUserFiltersFromContext(ctx)["name"]; hasFilter {
			fetchDefinitionsInput.FamilyPrefix = &givenFamilyPrefix
		}

		paginator := ecs.NewListTaskDefinitionsPaginator(conf.APIs.ECS, fetchDefinitionsInput)
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			arns = append(arns, out.TaskDefinitionArns...)
		}

		// Bounded: one goroutine per task definition ARN with no limit meant an
		// account with many definitions issued that many simultaneous
		// DescribeTaskDefinition calls. Results are collected rather than
		// streamed so the consumer below is unchanged; it already drains every
		// result and accumulates errors instead of returning on the first.
		var (
			mu        sync.Mutex
			collected []resStruct
		)
		descG := new(errgroup.Group)
		descG.SetLimit(maxParallelAWSCalls)

		for _, arn := range arns {
			taskDefArn := arn
			descG.Go(func() error {
				tasksOut, err := conf.APIs.ECS.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{TaskDefinition: &taskDefArn})
				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					collected = append(collected, resStruct{err: err})
					return nil
				}
				collected = append(collected, resStruct{res: tasksOut.TaskDefinition})
				return nil
			})
		}
		_ = descG.Wait() // per-item errors are carried in collected, as before

		var tasks []ecstypes.Task
		if val, e := cache.Get("getAllTasks", func() (any, error) {
			return getAllTasks(ctx, cache, conf.APIs.ECS)
		}); e != nil {
			return resources, objects, e
		} else if v, ok := val.([]ecstypes.Task); ok {
			tasks = v
		}

		var errs []string
		var err error

		for _, res := range collected {
			if res.err != nil {
				errs = appendIfNotInSlice(errs, res.err.Error())
				continue
			}
			objects = append(objects, *res.res)
			var graphres *graph.Resource
			if graphres, err = awsconv.NewResource(res.res); err != nil {
				errs = appendIfNotInSlice(errs, err.Error())
				continue
			}
			var deployments []*graph.KeyValue
			var runningServicesCount, stoppedServicesCount, runningTasksCount, stoppedTasksCount uint
			for _, t := range tasks {
				if awssdk.ToString(t.TaskDefinitionArn) == awssdk.ToString(res.res.TaskDefinitionArn) {
					group := awssdk.ToString(t.Group)
					state := strings.ToLower(awssdk.ToString(t.LastStatus))
					clusterArn := awssdk.ToString(t.ClusterArn)
					if strings.HasPrefix(group, "service:") {
						switch state {
						case "stopped":
							stoppedServicesCount++
							deployments = append(deployments, &graph.KeyValue{KeyName: arnToName(clusterArn), Value: group[len("service:"):] + " (stopped service)"})
						case "running":
							runningServicesCount++
							deployments = append(deployments, &graph.KeyValue{KeyName: arnToName(clusterArn), Value: group[len("service:"):] + " (running service)"})
						}
					}
					if strings.HasPrefix(group, "family:") {
						switch state {
						case "stopped":
							deployments = append(deployments, &graph.KeyValue{KeyName: arnToName(clusterArn), Value: group[len("family:"):] + " (stopped task)"})
							stoppedTasksCount++
						case "running":
							deployments = append(deployments, &graph.KeyValue{KeyName: arnToName(clusterArn), Value: group[len("family:"):] + " (running task)"})
							runningTasksCount++
						}
					}
				}
			}
			if len(deployments) > 0 {
				graphres.Properties()[properties.Deployments] = deployments
			}
			switch runningServicesCount + stoppedServicesCount + runningTasksCount + stoppedTasksCount {
			case 0:
				if state := strings.ToLower(string(res.res.Status)); state == "active" {
					graphres.Properties()[properties.State] = "ready"
				} else {
					graphres.Properties()[properties.State] = state
				}
			default:
				var stateSl []string
				if runningServicesCount > 0 {
					stateSl = append(stateSl, fmt.Sprintf("%d %s running", runningServicesCount, pluralizeIfNeeded("service", runningServicesCount)))
				}
				if stoppedServicesCount > 0 {
					stateSl = append(stateSl, fmt.Sprintf("%d %s stopped", stoppedServicesCount, pluralizeIfNeeded("service", runningServicesCount)))
				}
				if runningTasksCount > 0 {
					stateSl = append(stateSl, fmt.Sprintf("%d %s running", runningTasksCount, pluralizeIfNeeded("task", runningServicesCount)))
				}
				if stoppedTasksCount > 0 {
					stateSl = append(stateSl, fmt.Sprintf("%d %s stopped", stoppedTasksCount, pluralizeIfNeeded("task", runningServicesCount)))
				}
				if len(stateSl) > 0 {
					graphres.Properties()[properties.State] = strings.Join(stateSl, " ")
				}
			}

			resources = append(resources, graphres)
		}

		if len(errs) > 0 {
			err = fmt.Errorf("%s", strings.Join(errs, "; "))
		}

		return resources, objects, err
	}

	funcs["containercluster"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ecstypes.Cluster

		if !conf.getBoolDefaultTrue("aws.infra.containercluster.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[containercluster]")
			return resources, objects, nil
		}

		clusterNames, err := getClusterArns(ctx, cache, conf.APIs.ECS)
		if err != nil {
			return resources, objects, err
		}

		for _, clusterArns := range sliceOfSlice(clusterNames, 100) {
			clustersOut, err := conf.APIs.ECS.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: clusterArns})
			if err != nil {
				return resources, objects, err
			}

			for _, cluster := range clustersOut.Clusters {
				objects = append(objects, cluster)
				res, err := awsconv.NewResource(cluster)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}
		return resources, objects, nil
	}

	funcs["listener"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []elbv2types.Listener
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.listener.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[listener]")
			return resources, objects, nil
		}

		// Bounded and leak-free: previously one goroutine per load balancer wrote
		// to unbuffered channels while the consumer returned on the first error,
		// leaving the rest blocked on send forever.
		var lbs []elbv2types.LoadBalancer

		lbPaginator := elbv2.NewDescribeLoadBalancersPaginator(conf.APIs.Elbv2, &elbv2.DescribeLoadBalancersInput{})
		for lbPaginator.HasMorePages() {
			out, err := lbPaginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			lbs = append(lbs, out.LoadBalancers...)
		}

		var mu sync.Mutex
		g := new(errgroup.Group)
		g.SetLimit(maxParallelAWSCalls)

		for _, lb := range lbs {
			g.Go(func() error {
				listenerPaginator := elbv2.NewDescribeListenersPaginator(conf.APIs.Elbv2, &elbv2.DescribeListenersInput{LoadBalancerArn: lb.LoadBalancerArn})
				for listenerPaginator.HasMorePages() {
					lout, err := listenerPaginator.NextPage(ctx)
					if err != nil {
						return err
					}
					for _, listen := range lout.Listeners {
						res, err := awsconv.NewResource(listen)
						if err != nil {
							return err
						}
						mu.Lock()
						objects = append(objects, listen)
						resources = append(resources, res)
						mu.Unlock()
					}
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return resources, objects, err
		}
		return resources, objects, nil
	}
}

func addManualAccessFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["user"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []iamtypes.UserDetail

		if !conf.getBoolDefaultTrue("aws.access.user.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[user]")
			return resources, objects, nil
		}

		var wg sync.WaitGroup
		resourcesC := make(chan *graph.Resource)
		objectsC := make(chan iamtypes.UserDetail)
		errC := make(chan error)

		wg.Add(1)
		go func() {
			defer wg.Done()
			accountDetails, err := getAccountAuthorizationDetails(ctx, cache, conf.APIs.IAM)
			if err != nil {
				errC <- err
				return
			}
			for _, output := range accountDetails.Users {
				objectsC <- output
				if res, e := awsconv.NewResource(output); e != nil {
					errC <- e
					return
				} else {
					resourcesC <- res
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			paginator := iam.NewListUsersPaginator(conf.APIs.IAM, &iam.ListUsersInput{})
			for paginator.HasMorePages() {
				page, err := paginator.NextPage(ctx)
				if err != nil {
					errC <- err
					return
				}
				for _, user := range page.Users {
					res, e := awsconv.NewResource(user)
					if e != nil {
						errC <- e
						return
					}
					resourcesC <- res
				}
			}
		}()

		go func() {
			wg.Wait()
			close(errC)
			close(objectsC)
			close(resourcesC)
		}()

		for {
			select {
			case e := <-errC:
				if e != nil {
					return resources, objects, e
				}
			case r, ok := <-resourcesC:
				if !ok {
					return resources, objects, nil
				}
				if r != nil {
					resources = append(resources, r)
				}
			case o, ok := <-objectsC:
				if !ok {
					return resources, objects, nil
				}
				objects = append(objects, o)
			}
		}
	}

	funcs["group"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []iamtypes.GroupDetail

		if !conf.getBoolDefaultTrue("aws.access.group.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[group]")
			return resources, objects, nil
		}

		accountDetails, err := getAccountAuthorizationDetails(ctx, cache, conf.APIs.IAM)
		if err != nil {
			return resources, objects, err
		}

		for _, output := range accountDetails.Groups {
			objects = append(objects, output)
			if res, err := awsconv.NewResource(output); err != nil {
				return resources, objects, err
			} else {
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["role"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []iamtypes.RoleDetail

		if !conf.getBoolDefaultTrue("aws.access.role.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[role]")
			return resources, objects, nil
		}

		accountDetails, err := getAccountAuthorizationDetails(ctx, cache, conf.APIs.IAM)
		if err != nil {
			return resources, objects, err
		}

		for _, output := range accountDetails.Roles {
			objects = append(objects, output)
			if res, err := awsconv.NewResource(output); err != nil {
				return resources, objects, err
			} else {
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}

	funcs["policy"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []iamtypes.ManagedPolicyDetail

		if !conf.getBoolDefaultTrue("aws.access.policy.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[policy]")
			return resources, objects, nil
		}

		errC := make(chan error)
		objectsC := make(chan iamtypes.ManagedPolicyDetail)
		resourcesC := make(chan *graph.Resource)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()

			accountDetails, err := getAccountAuthorizationDetails(ctx, cache, conf.APIs.IAM)
			if err != nil {
				errC <- err
				return
			}
			for _, p := range accountDetails.Policies {
				res, e := awsconv.NewResource(p)
				if e != nil {
					errC <- e
					return
				}
				if strings.HasPrefix(awssdk.ToString(p.Arn), "arn:aws:iam::aws:policy") {
					res.Properties()[properties.Type] = "AWS Managed"
				} else {
					res.Properties()[properties.Type] = "Customer Managed"
				}
				res.Properties()[properties.Attached] = awssdk.ToInt32(p.AttachmentCount) > 0
				resourcesC <- res
			}
		}()

		go func() {
			wg.Wait()
			close(errC)
			close(objectsC)
			close(resourcesC)
		}()

		for {
			select {
			case err := <-errC:
				if err != nil {
					return resources, objects, err
				}
			case o, ok := <-objectsC:
				if !ok {
					return resources, objects, nil
				}
				objects = append(objects, o)
			case r, ok := <-resourcesC:
				if !ok {
					return resources, objects, nil
				}
				resources = append(resources, r)

			}
		}
	}
	funcs["accesskey"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []iamtypes.AccessKeyMetadata

		if !conf.getBoolDefaultTrue("aws.access.accesskey.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[accesskey]")
			return resources, objects, nil
		}

		// Previously this had three defects beyond the unbounded fan-out: the
		// `defer wg.Done()` was registered *after* an early return, so a failing
		// InitResource left the counter high and wg.Wait() blocked forever;
		// `hasError` was written from several goroutines without
		// synchronization; and the consumer returned on the first error, leaving
		// the remaining senders blocked on unbuffered channels.
		var (
			mu  sync.Mutex
			g   = new(errgroup.Group)
			all []iamtypes.User
		)
		g.SetLimit(maxParallelAWSCalls)

		usersPaginator := iam.NewListUsersPaginator(conf.APIs.IAM, &iam.ListUsersInput{})
		for usersPaginator.HasMorePages() {
			outUsers, err := usersPaginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			all = append(all, outUsers.Users...)
		}

		for _, user := range all {
			u := user
			g.Go(func() error {
				userRes, err := awsconv.InitResource(u)
				if err != nil {
					return err
				}

				akPaginator := iam.NewListAccessKeysPaginator(conf.APIs.IAM, &iam.ListAccessKeysInput{UserName: u.UserName})
				for akPaginator.HasMorePages() {
					out, err := akPaginator.NextPage(ctx)
					if err != nil {
						return err
					}
					for _, output := range out.AccessKeyMetadata {
						res, err := awsconv.NewResource(output)
						if err != nil {
							return err
						}
						res.AddRelation(rdf.ChildrenOfRel, userRes)

						mu.Lock()
						objects = append(objects, output)
						resources = append(resources, res)
						mu.Unlock()
					}
				}
				return nil
			})
		}

		// Partial results are still returned alongside an error, as before.
		if err := g.Wait(); err != nil {
			return resources, objects, err
		}
		return resources, objects, nil
	}
}
func addManualStorageFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["bucket"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []s3types.Bucket

		if !conf.getBoolDefaultTrue("aws.storage.bucket.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource storage[bucket]")
			return resources, objects, nil
		}

		bucketM := &sync.Mutex{}

		err := forEachBucketParallel(ctx, cache, conf.APIs.S3, func(b s3types.Bucket) error {
			bucketM.Lock()
			objects = append(objects, b)
			bucketM.Unlock()
			res, err := awsconv.NewResource(b)
			if err != nil {
				return fmt.Errorf("build resource for bucket `%s`: %w", awssdk.ToString(b.Name), err)
			}
			grants, err := fetchAndExtractGrantsFn(ctx, conf.APIs.S3, awssdk.ToString(b.Name))
			if err != nil {
				return fmt.Errorf("fetching grants for bucket %s: %w", awssdk.ToString(b.Name), err)
			}
			res.Properties()[properties.Grants] = grants
			bucketM.Lock()
			resources = append(resources, res)
			bucketM.Unlock()
			return nil
		})
		return resources, objects, err
	}

	funcs["s3object"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []s3types.Object
		var resources []*graph.Resource

		resourcesC := make(chan *graph.Resource)

		if !conf.getBoolDefaultTrue("aws.storage.s3object.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource storage[s3object]")
			return resources, objects, nil
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range resourcesC {
				resources = append(resources, r)
			}
		}()

		err := forEachBucketParallel(ctx, cache, conf.APIs.S3, func(b s3types.Bucket) error {
			return fetchObjectsForBucket(ctx, conf.APIs.S3, b, resourcesC)
		})

		close(resourcesC)

		wg.Wait()

		return resources, objects, err
	}
}
func addManualMessagingFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["queue"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []string
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.messaging.queue.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource messaging[queue]")
			return resources, objects, nil
		}

		out, err := conf.APIs.SQS.ListQueues(ctx, &sqs.ListQueuesInput{})
		if err != nil {
			return nil, objects, err
		}

		// Bounded and leak-free: previously one goroutine per queue wrote to
		// unbuffered channels while the consumer returned on the first error,
		// leaving the rest blocked on send forever.
		var mu sync.Mutex
		g := new(errgroup.Group)
		g.SetLimit(maxParallelAWSCalls)

		for _, output := range out.QueueUrls {
			url := output
			g.Go(func() error {
				mu.Lock()
				objects = append(objects, url)
				mu.Unlock()
				res := graph.InitResource(cloud.Queue, url)
				res.Properties()[properties.ID] = url
				attrs, err := conf.APIs.SQS.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{AttributeNames: []sqstypes.QueueAttributeName{sqstypes.QueueAttributeNameAll}, QueueUrl: &url})
				if err != nil {
					var apiErr smithy.APIError
					if errors.As(err, &apiErr) && (apiErr.ErrorCode() == "AWS.SimpleQueueService.NonExistentQueue" || apiErr.ErrorCode() == "AWS.SimpleQueueService.QueueDeletedRecently") {
						// Queue vanished between List and Get; not an error.
						return nil
					}
					return err
				}
				for k, v := range attrs.Attributes {
					switch k {
					case "ApproximateNumberOfMessages":
						count, err := strconv.Atoi(v)
						if err != nil {
							return err
						}
						res.Properties()[properties.ApproximateMessageCount] = count
					case "CreatedTimestamp":
						if v != "" {
							timestamp, err := strconv.ParseInt(v, 10, 64)
							if err != nil {
								return err
							}
							res.Properties()[properties.Created] = time.Unix(timestamp, 0)
						}
					case "LastModifiedTimestamp":
						if v != "" {
							timestamp, err := strconv.ParseInt(v, 10, 64)
							if err != nil {
								return err
							}
							res.Properties()[properties.Modified] = time.Unix(timestamp, 0)
						}
					case "QueueArn":
						res.Properties()[properties.Arn] = v
					case "DelaySeconds":
						delay, err := strconv.Atoi(v)
						if err != nil {
							return err
						}
						res.Properties()[properties.Delay] = delay
					}

				}
				mu.Lock()
				resources = append(resources, res)
				mu.Unlock()
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return resources, objects, err
		}
		return resources, objects, nil
	}
}
func addManualDNSFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["record"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []route53types.ResourceRecordSet
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.dns.record.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource dns[record]")
			return resources, objects, nil
		}

		zoneName, hasZoneFilter := getUserFiltersFromContext(ctx)["zone"]

		// Two bounded phases replacing a three-stage channel pipeline. The old
		// version leaked: one goroutine per hosted zone wrote to unbuffered
		// channels while the consumer returned on the first error, leaving the
		// rest blocked on send forever.
		var zones []route53types.HostedZone

		paginator := route53.NewListHostedZonesPaginator(conf.APIs.Route53, &route53.ListHostedZonesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, output := range out.HostedZones {
				if hasZoneFilter && !strings.Contains(strings.ToLower(awssdk.ToString(output.Name)), strings.ToLower(zoneName)) {
					continue
				}
				zones = append(zones, output)
			}
		}

		var mu sync.Mutex
		g := new(errgroup.Group)
		g.SetLimit(maxParallelAWSCalls)

		for _, zone := range zones {
			z := zone
			g.Go(func() error {
				parent, err := awsconv.InitResource(z)
				if err != nil {
					return err
				}

				input := &route53.ListResourceRecordSetsInput{HostedZoneId: z.Id}
				for {
					out, err := conf.APIs.Route53.ListResourceRecordSets(ctx, input)
					if err != nil {
						return err
					}
					for _, output := range out.ResourceRecordSets {
						res, err := awsconv.NewResource(output)
						if err != nil {
							return err
						}
						res.Properties()[properties.Zone] = awssdk.ToString(z.Name)
						res.AddRelation(rdf.ChildrenOfRel, parent)

						mu.Lock()
						objects = append(objects, output)
						resources = append(resources, res)
						mu.Unlock()
					}
					if !out.IsTruncated {
						break
					}
					input.StartRecordName = out.NextRecordName
					input.StartRecordType = out.NextRecordType
					input.StartRecordIdentifier = out.NextRecordIdentifier
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return resources, objects, err
		}
		return resources, objects, nil
	}
}
func addManualLambdaFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualMonitoringFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualCDNFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualCloudformationFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

func addManualEKSFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["ekscluster"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ekstypes.Cluster

		if !conf.getBoolDefaultTrue("aws.eks.ekscluster.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource eks[ekscluster]")
			return resources, objects, nil
		}

		paginator := eks.NewListClustersPaginator(conf.APIs.EKS, &eks.ListClustersInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, name := range out.Clusters {
				descOut, err := conf.APIs.EKS.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: &name})
				if err != nil {
					return resources, objects, err
				}
				if descOut.Cluster != nil {
					objects = append(objects, *descOut.Cluster)
					res, err := awsconv.NewResource(*descOut.Cluster)
					if err != nil {
						return resources, objects, err
					}
					resources = append(resources, res)
				}
			}
		}
		return resources, objects, nil
	}

	funcs["eksnodegroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []ekstypes.Nodegroup

		if !conf.getBoolDefaultTrue("aws.eks.eksnodegroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource eks[eksnodegroup]")
			return resources, objects, nil
		}

		clusterPaginator := eks.NewListClustersPaginator(conf.APIs.EKS, &eks.ListClustersInput{})
		for clusterPaginator.HasMorePages() {
			cOut, err := clusterPaginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, clusterName := range cOut.Clusters {
				ngPaginator := eks.NewListNodegroupsPaginator(conf.APIs.EKS, &eks.ListNodegroupsInput{ClusterName: &clusterName})
				for ngPaginator.HasMorePages() {
					ngOut, err := ngPaginator.NextPage(ctx)
					if err != nil {
						return resources, objects, err
					}
					for _, ngName := range ngOut.Nodegroups {
						descOut, err := conf.APIs.EKS.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{ClusterName: &clusterName, NodegroupName: &ngName})
						if err != nil {
							return resources, objects, err
						}
						if descOut.Nodegroup != nil {
							objects = append(objects, *descOut.Nodegroup)
							res, err := awsconv.NewResource(*descOut.Nodegroup)
							if err != nil {
								return resources, objects, err
							}
							parent := graph.InitResource(cloud.EKSCluster, clusterName)
							res.AddRelation(rdf.ChildrenOfRel, parent)
							resources = append(resources, res)
						}
					}
				}
			}
		}
		return resources, objects, nil
	}
}

func addManualDynamodbFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["dynamodbtable"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []dynamodbtypes.TableDescription

		if !conf.getBoolDefaultTrue("aws.dynamodb.dynamodbtable.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource dynamodb[dynamodbtable]")
			return resources, objects, nil
		}

		paginator := dynamodb.NewListTablesPaginator(conf.APIs.Dynamodb, &dynamodb.ListTablesInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, name := range out.TableNames {
				descOut, err := conf.APIs.Dynamodb.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: &name})
				if err != nil {
					return resources, objects, err
				}
				if descOut.Table != nil {
					objects = append(objects, *descOut.Table)
					res, err := awsconv.NewResource(*descOut.Table)
					if err != nil {
						return resources, objects, err
					}
					resources = append(resources, res)
				}
			}
		}
		return resources, objects, nil
	}
}

func addManualSecretsmanagerFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["key"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []kmstypes.KeyMetadata

		if !conf.getBoolDefaultTrue("aws.secretsmanager.key.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource secretsmanager[key]")
			return resources, objects, nil
		}

		paginator := kms.NewListKeysPaginator(conf.APIs.KMS, &kms.ListKeysInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, keyEntry := range out.Keys {
				descOut, err := conf.APIs.KMS.DescribeKey(ctx, &kms.DescribeKeyInput{KeyId: keyEntry.KeyId})
				if err != nil {
					return resources, objects, err
				}
				if descOut.KeyMetadata != nil {
					objects = append(objects, *descOut.KeyMetadata)
					res, err := awsconv.NewResource(*descOut.KeyMetadata)
					if err != nil {
						return resources, objects, err
					}
					resources = append(resources, res)
				}
			}
		}
		return resources, objects, nil
	}
}

func addManualApigatewayFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["apigateway"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []apigatewayv2types.Api

		if !conf.getBoolDefaultTrue("aws.apigateway.apigateway.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource apigateway[apigateway]")
			return resources, objects, nil
		}

		var nextToken *string
		for {
			out, err := conf.APIs.Apigatewayv2.GetApis(ctx, &apigatewayv2.GetApisInput{NextToken: nextToken})
			if err != nil {
				return resources, objects, err
			}
			for _, api := range out.Items {
				objects = append(objects, api)
				res, err := awsconv.NewResource(api)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
			if out.NextToken == nil {
				break
			}
			nextToken = out.NextToken
		}
		return resources, objects, nil
	}

	funcs["apigatewayroute"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []apigatewayv2types.Route

		if !conf.getBoolDefaultTrue("aws.apigateway.apigatewayroute.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource apigateway[apigatewayroute]")
			return resources, objects, nil
		}

		apisOut, err := conf.APIs.Apigatewayv2.GetApis(ctx, &apigatewayv2.GetApisInput{})
		if err != nil {
			return resources, objects, err
		}
		for _, api := range apisOut.Items {
			routesOut, err := conf.APIs.Apigatewayv2.GetRoutes(ctx, &apigatewayv2.GetRoutesInput{ApiId: api.ApiId})
			if err != nil {
				return resources, objects, err
			}
			for _, route := range routesOut.Items {
				objects = append(objects, route)
				res, err := awsconv.NewResource(route)
				if err != nil {
					return resources, objects, err
				}
				parent := graph.InitResource(cloud.APIGateway, awssdk.ToString(api.ApiId))
				res.AddRelation(rdf.ChildrenOfRel, parent)
				resources = append(resources, res)
			}
		}
		return resources, objects, nil
	}

	funcs["apigatewaystage"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []apigatewayv2types.Stage

		if !conf.getBoolDefaultTrue("aws.apigateway.apigatewaystage.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource apigateway[apigatewaystage]")
			return resources, objects, nil
		}

		apisOut, err := conf.APIs.Apigatewayv2.GetApis(ctx, &apigatewayv2.GetApisInput{})
		if err != nil {
			return resources, objects, err
		}
		for _, api := range apisOut.Items {
			stagesOut, err := conf.APIs.Apigatewayv2.GetStages(ctx, &apigatewayv2.GetStagesInput{ApiId: api.ApiId})
			if err != nil {
				return resources, objects, err
			}
			for _, stage := range stagesOut.Items {
				objects = append(objects, stage)
				res, err := awsconv.NewResource(stage)
				if err != nil {
					return resources, objects, err
				}
				parent := graph.InitResource(cloud.APIGateway, awssdk.ToString(api.ApiId))
				res.AddRelation(rdf.ChildrenOfRel, parent)
				resources = append(resources, res)
			}
		}
		return resources, objects, nil
	}
}

func addManualSSMFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

func addManualCloudtrailFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

func addManualEFSFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["mounttarget"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var resources []*graph.Resource
		var objects []efstypes.MountTargetDescription

		if !conf.getBoolDefaultTrue("aws.efs.mounttarget.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource efs[mounttarget]")
			return resources, objects, nil
		}

		fsOut, err := conf.APIs.EFS.DescribeFileSystems(ctx, &efs.DescribeFileSystemsInput{})
		if err != nil {
			return resources, objects, err
		}
		for _, fs := range fsOut.FileSystems {
			mtOut, err := conf.APIs.EFS.DescribeMountTargets(ctx, &efs.DescribeMountTargetsInput{FileSystemId: fs.FileSystemId})
			if err != nil {
				return resources, objects, err
			}
			for _, mt := range mtOut.MountTargets {
				objects = append(objects, mt)
				res, err := awsconv.NewResource(mt)
				if err != nil {
					return resources, objects, err
				}
				parent := graph.InitResource(cloud.FileSystem, awssdk.ToString(fs.FileSystemId))
				res.AddRelation(rdf.ChildrenOfRel, parent)
				resources = append(resources, res)
			}
		}
		return resources, objects, nil
	}
}

func addManualCloudwatchlogsFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

func addManualElasticacheFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

func addManualEventbridgeFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["eventbus"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []eventbridgetypes.EventBus
		var resources []*graph.Resource

		var next *string
		for {
			out, err := conf.APIs.Eventbridge.ListEventBuses(ctx, &eventbridge.ListEventBusesInput{NextToken: next})
			if err != nil {
				return resources, objects, err
			}
			for _, bus := range out.EventBuses {
				objects = append(objects, bus)
				res, err := awsconv.NewResource(bus)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
			if out.NextToken == nil || awssdk.ToString(out.NextToken) == "" {
				break
			}
			next = out.NextToken
		}

		return resources, objects, nil
	}

	funcs["eventrule"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []eventbridgetypes.Rule
		var resources []*graph.Resource

		var next *string
		for {
			out, err := conf.APIs.Eventbridge.ListRules(ctx, &eventbridge.ListRulesInput{NextToken: next})
			if err != nil {
				return resources, objects, err
			}
			for _, rule := range out.Rules {
				objects = append(objects, rule)
				res, err := awsconv.NewResource(rule)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
			if out.NextToken == nil || awssdk.ToString(out.NextToken) == "" {
				break
			}
			next = out.NextToken
		}

		return resources, objects, nil
	}
}

func addManualStepfunctionsFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

// wafScopes returns the scopes worth querying from the configured region. CLOUDFRONT
// resources live in a global namespace that only us-east-1 can see; asking for them
// anywhere else fails the whole fetch, so they are only requested where they exist.
func wafScopes(region string) []wafv2types.Scope {
	if region == "us-east-1" {
		return []wafv2types.Scope{wafv2types.ScopeRegional, wafv2types.ScopeCloudfront}
	}
	return []wafv2types.Scope{wafv2types.ScopeRegional}
}

func addManualWafFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["webacl"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []wafv2types.WebACLSummary
		var resources []*graph.Resource

		for _, scope := range wafScopes(conf.APIs.Wafv2.Options().Region) {
			var next *string
			for {
				out, err := conf.APIs.Wafv2.ListWebACLs(ctx, &wafv2.ListWebACLsInput{Scope: scope, NextMarker: next})
				if err != nil {
					return resources, objects, err
				}
				for _, acl := range out.WebACLs {
					objects = append(objects, acl)
					res, err := awsconv.NewResource(acl)
					if err != nil {
						return resources, objects, err
					}
					// The summary carries no scope, but it is the only thing
					// distinguishing two ACLs that may share a name.
					res.Properties()[properties.Scope] = string(scope)
					resources = append(resources, res)
				}
				if out.NextMarker == nil || awssdk.ToString(out.NextMarker) == "" {
					break
				}
				next = out.NextMarker
			}
		}

		return resources, objects, nil
	}

	funcs["ipset"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []wafv2types.IPSetSummary
		var resources []*graph.Resource

		for _, scope := range wafScopes(conf.APIs.Wafv2.Options().Region) {
			var next *string
			for {
				out, err := conf.APIs.Wafv2.ListIPSets(ctx, &wafv2.ListIPSetsInput{Scope: scope, NextMarker: next})
				if err != nil {
					return resources, objects, err
				}
				for _, set := range out.IPSets {
					objects = append(objects, set)
					res, err := awsconv.NewResource(set)
					if err != nil {
						return resources, objects, err
					}
					res.Properties()[properties.Scope] = string(scope)
					resources = append(resources, res)
				}
				if out.NextMarker == nil || awssdk.ToString(out.NextMarker) == "" {
					break
				}
				next = out.NextMarker
			}
		}

		return resources, objects, nil
	}

	funcs["rulegroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []wafv2types.RuleGroupSummary
		var resources []*graph.Resource

		for _, scope := range wafScopes(conf.APIs.Wafv2.Options().Region) {
			var next *string
			for {
				out, err := conf.APIs.Wafv2.ListRuleGroups(ctx, &wafv2.ListRuleGroupsInput{Scope: scope, NextMarker: next})
				if err != nil {
					return resources, objects, err
				}
				for _, group := range out.RuleGroups {
					objects = append(objects, group)
					res, err := awsconv.NewResource(group)
					if err != nil {
						return resources, objects, err
					}
					res.Properties()[properties.Scope] = string(scope)
					resources = append(resources, res)
				}
				if out.NextMarker == nil || awssdk.ToString(out.NextMarker) == "" {
					break
				}
				next = out.NextMarker
			}
		}

		return resources, objects, nil
	}
}

func addManualConfigserviceFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["configrule"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, any, error) {
		var objects []configservicetypes.ConfigRule
		var resources []*graph.Resource

		// Compliance first, so every rule can be annotated as it is built. Keyed by
		// rule name, which is what both APIs agree on.
		compliance := make(map[string]string)
		compPager := configservice.NewDescribeComplianceByConfigRulePaginator(conf.APIs.Configservice, &configservice.DescribeComplianceByConfigRuleInput{})
		for compPager.HasMorePages() {
			out, err := compPager.NextPage(ctx)
			if err != nil {
				// Compliance is an enrichment; a rule list without it is still
				// useful, so this does not fail the fetch.
				conf.Log.Verbosef("sync: cannot read Config rule compliance: %s", err)
				break
			}
			for _, c := range out.ComplianceByConfigRules {
				if c.Compliance != nil {
					compliance[awssdk.ToString(c.ConfigRuleName)] = string(c.Compliance.ComplianceType)
				}
			}
		}

		pager := configservice.NewDescribeConfigRulesPaginator(conf.APIs.Configservice, &configservice.DescribeConfigRulesInput{})
		for pager.HasMorePages() {
			out, err := pager.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, rule := range out.ConfigRules {
				objects = append(objects, rule)
				res, err := awsconv.NewResource(rule)
				if err != nil {
					return resources, objects, err
				}
				if c, ok := compliance[awssdk.ToString(rule.ConfigRuleName)]; ok {
					res.Properties()[properties.Compliance] = c
				}
				resources = append(resources, res)
			}
		}

		return resources, objects, nil
	}
}
