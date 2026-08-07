package awsfetch

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

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

	awsconv "github.com/bootswithdefer/awless/aws/conv"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/cloud/properties"
	"github.com/bootswithdefer/awless/cloud/rdf"
	"github.com/bootswithdefer/awless/fetch"
	"github.com/bootswithdefer/awless/graph"
)

func addManualInfraFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["containerinstance"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []ecstypes.ContainerInstance
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.containerinstance.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[containerinstance]")
			return resources, objects, nil
		}

		clusterArns, err := getClusterArns(ctx, cache, conf.APIs.Ecs)
		if err != nil {
			return resources, objects, err
		}

		for _, cluster := range clusterArns {
			paginator := ecs.NewListContainerInstancesPaginator(conf.APIs.Ecs, &ecs.ListContainerInstancesInput{Cluster: &cluster})
			for paginator.HasMorePages() {
				out, err := paginator.NextPage(ctx)
				if err != nil {
					return resources, objects, err
				}
				if len(out.ContainerInstanceArns) == 0 {
					continue
				}

				containerInstancesOut, err := conf.APIs.Ecs.DescribeContainerInstances(ctx, &ecs.DescribeContainerInstancesInput{Cluster: &cluster, ContainerInstances: out.ContainerInstanceArns})
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

	funcs["container"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []ecstypes.Container
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.container.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[container]")
			return resources, objects, nil
		}

		var tasks []ecstypes.Task

		if val, e := cache.Get("getAllTasks", func() (interface{}, error) {
			return getAllTasks(ctx, cache, conf.APIs.Ecs)
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

	funcs["containertask"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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

		var wg sync.WaitGroup
		resc := make(chan resStruct)

		fetchDefinitionsInput := &ecs.ListTaskDefinitionsInput{}
		if givenFamilyPrefix, hasFilter := getUserFiltersFromContext(ctx)["name"]; hasFilter {
			fetchDefinitionsInput.FamilyPrefix = &givenFamilyPrefix
		}

		paginator := ecs.NewListTaskDefinitionsPaginator(conf.APIs.Ecs, fetchDefinitionsInput)
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, arn := range out.TaskDefinitionArns {
				wg.Add(1)
				go func(taskDefArn string) {
					defer wg.Done()
					tasksOut, err := conf.APIs.Ecs.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{TaskDefinition: &taskDefArn})
					if err != nil {
						resc <- resStruct{err: err}
						return
					}
					resc <- resStruct{res: tasksOut.TaskDefinition}
				}(arn)
			}
		}

		go func() {
			wg.Wait()
			close(resc)
		}()

		var tasks []ecstypes.Task
		if val, e := cache.Get("getAllTasks", func() (interface{}, error) {
			return getAllTasks(ctx, cache, conf.APIs.Ecs)
		}); e != nil {
			return resources, objects, e
		} else if v, ok := val.([]ecstypes.Task); ok {
			tasks = v
		}

		var errs []string
		var err error

		for res := range resc {
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
			switch {
			case runningServicesCount+stoppedServicesCount+runningTasksCount+stoppedTasksCount == 0:
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

	funcs["containercluster"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []ecstypes.Cluster

		if !conf.getBoolDefaultTrue("aws.infra.containercluster.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[containercluster]")
			return resources, objects, nil
		}

		clusterNames, err := getClusterArns(ctx, cache, conf.APIs.Ecs)
		if err != nil {
			return resources, objects, nil
		}

		for _, clusterArns := range sliceOfSlice(clusterNames, 100) {
			clustersOut, err := conf.APIs.Ecs.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: clusterArns})
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

	funcs["listener"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []elbv2types.Listener
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.infra.listener.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource infra[listener]")
			return resources, objects, nil
		}

		errc := make(chan error)
		resultc := make(chan elbv2types.Listener)
		var wg sync.WaitGroup

		lbPaginator := elbv2.NewDescribeLoadBalancersPaginator(conf.APIs.Elbv2, &elbv2.DescribeLoadBalancersInput{})
		for lbPaginator.HasMorePages() {
			out, err := lbPaginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, lb := range out.LoadBalancers {
				wg.Add(1)
				go func(lb elbv2types.LoadBalancer) {
					defer wg.Done()
					listenerPaginator := elbv2.NewDescribeListenersPaginator(conf.APIs.Elbv2, &elbv2.DescribeListenersInput{LoadBalancerArn: lb.LoadBalancerArn})
					for listenerPaginator.HasMorePages() {
						lout, lerr := listenerPaginator.NextPage(ctx)
						if lerr != nil {
							errc <- lerr
							return
						}
						for _, listen := range lout.Listeners {
							resultc <- listen
						}
					}
				}(lb)
			}
		}

		go func() {
			wg.Wait()
			close(resultc)
		}()

		for {
			select {
			case err := <-errc:
				if err != nil {
					return resources, objects, err
				}
			case listener, ok := <-resultc:
				if !ok {
					return resources, objects, nil
				}
				objects = append(objects, listener)
				res, err := awsconv.NewResource(listener)
				if err != nil {
					return resources, objects, err
				}
				resources = append(resources, res)
			}
		}
	}
}

func addManualAccessFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["user"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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
			accountDetails, err := getAccountAuthorizationDetails(ctx, cache, conf.APIs.Iam)
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
			paginator := iam.NewListUsersPaginator(conf.APIs.Iam, &iam.ListUsersInput{})
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

	funcs["group"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []iamtypes.GroupDetail

		if !conf.getBoolDefaultTrue("aws.access.group.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[group]")
			return resources, objects, nil
		}

		accountDetails, err := getAccountAuthorizationDetails(ctx, cache, conf.APIs.Iam)
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

	funcs["role"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []iamtypes.RoleDetail

		if !conf.getBoolDefaultTrue("aws.access.role.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[role]")
			return resources, objects, nil
		}

		accountDetails, err := getAccountAuthorizationDetails(ctx, cache, conf.APIs.Iam)
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

	funcs["policy"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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

			accountDetails, err := getAccountAuthorizationDetails(ctx, cache, conf.APIs.Iam)
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
	funcs["accesskey"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []iamtypes.AccessKeyMetadata

		if !conf.getBoolDefaultTrue("aws.access.accesskey.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource access[accesskey]")
			return resources, objects, nil
		}

		var wg sync.WaitGroup
		resourcesC := make(chan *graph.Resource)
		objectsC := make(chan iamtypes.AccessKeyMetadata)
		errC := make(chan error)
		var hasError bool

		usersPaginator := iam.NewListUsersPaginator(conf.APIs.Iam, &iam.ListUsersInput{})
		for usersPaginator.HasMorePages() && !hasError {
			outUsers, err := usersPaginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}

			for _, user := range outUsers.Users {
				wg.Add(1)
				go func(u iamtypes.User) {
					userRes, err := awsconv.InitResource(u)
					if err != nil {
						hasError = true
						errC <- err
						return
					}
					defer wg.Done()

					akPaginator := iam.NewListAccessKeysPaginator(conf.APIs.Iam, &iam.ListAccessKeysInput{UserName: u.UserName})
					for akPaginator.HasMorePages() {
						out, err := akPaginator.NextPage(ctx)
						if err != nil {
							hasError = true
							errC <- err
							return
						}
						for _, output := range out.AccessKeyMetadata {
							objectsC <- output
							res, e := awsconv.NewResource(output)
							if e != nil {
								errC <- e
								hasError = true
								return
							}
							res.AddRelation(rdf.ChildrenOfRel, userRes)
							resourcesC <- res
						}
					}
				}(user)
			}
		}

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
}
func addManualStorageFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["bucket"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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

	funcs["s3object"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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
	funcs["queue"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []string
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.messaging.queue.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource messaging[queue]")
			return resources, objects, nil
		}

		out, err := conf.APIs.Sqs.ListQueues(ctx, &sqs.ListQueuesInput{})
		if err != nil {
			return nil, objects, err
		}

		errC := make(chan error)
		objectsC := make(chan string)
		resourcesC := make(chan *graph.Resource)
		var wg sync.WaitGroup

		for _, output := range out.QueueUrls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				objectsC <- url
				res := graph.InitResource(cloud.Queue, url)
				res.Properties()[properties.ID] = url
				attrs, err := conf.APIs.Sqs.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{AttributeNames: []sqstypes.QueueAttributeName{sqstypes.QueueAttributeNameAll}, QueueUrl: &url})
				if err != nil {
					var apiErr smithy.APIError
					if errors.As(err, &apiErr) && (apiErr.ErrorCode() == "AWS.SimpleQueueService.NonExistentQueue" || apiErr.ErrorCode() == "AWS.SimpleQueueService.QueueDeletedRecently") {
						return
					}
					errC <- err
					return
				}
				for k, v := range attrs.Attributes {
					switch k {
					case "ApproximateNumberOfMessages":
						count, err := strconv.Atoi(v)
						if err != nil {
							errC <- err
						}
						res.Properties()[properties.ApproximateMessageCount] = count
					case "CreatedTimestamp":
						if v != "" {
							timestamp, err := strconv.ParseInt(v, 10, 64)
							if err != nil {
								errC <- err
							}
							res.Properties()[properties.Created] = time.Unix(timestamp, 0)
						}
					case "LastModifiedTimestamp":
						if v != "" {
							timestamp, err := strconv.ParseInt(v, 10, 64)
							if err != nil {
								errC <- err
							}
							res.Properties()[properties.Modified] = time.Unix(timestamp, 0)
						}
					case "QueueArn":
						res.Properties()[properties.Arn] = v
					case "DelaySeconds":
						delay, err := strconv.Atoi(v)
						if err != nil {
							errC <- err
						}
						res.Properties()[properties.Delay] = delay
					}

				}
				resourcesC <- res
			}(output)

		}

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
}
func addManualDnsFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["record"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var objects []route53types.ResourceRecordSet
		var resources []*graph.Resource

		if !conf.getBoolDefaultTrue("aws.dns.record.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource dns[record]")
			return resources, objects, nil
		}

		zoneName, hasZoneFilter := getUserFiltersFromContext(ctx)["zone"]

		errC := make(chan error)
		zoneC := make(chan route53types.HostedZone)
		objectsC := make(chan route53types.ResourceRecordSet)
		resourcesC := make(chan *graph.Resource)

		go func() {
			paginator := route53.NewListHostedZonesPaginator(conf.APIs.Route53, &route53.ListHostedZonesInput{})
			for paginator.HasMorePages() {
				out, err := paginator.NextPage(ctx)
				if err != nil {
					errC <- err
					break
				}
				for _, output := range out.HostedZones {
					if hasZoneFilter {
						if strings.Contains(strings.ToLower(awssdk.ToString(output.Name)), strings.ToLower(zoneName)) {
							zoneC <- output
						}
					} else {
						zoneC <- output
					}
				}
			}
			close(zoneC)
		}()

		go func() {
			var wg sync.WaitGroup

			for zone := range zoneC {
				wg.Add(1)
				go func(z route53types.HostedZone) {
					defer wg.Done()
					input := &route53.ListResourceRecordSetsInput{HostedZoneId: z.Id}
					for {
						out, err := conf.APIs.Route53.ListResourceRecordSets(ctx, input)
						if err != nil {
							errC <- err
							return
						}
						for _, output := range out.ResourceRecordSets {
							objectsC <- output
							res, err := awsconv.NewResource(output)
							if err != nil {
								errC <- err
								return
							}
							res.Properties()[properties.Zone] = awssdk.ToString(z.Name)

							parent, err := awsconv.InitResource(z)
							if err != nil {
								errC <- err
								return
							}
							res.AddRelation(rdf.ChildrenOfRel, parent)
							resourcesC <- res
						}
						if !out.IsTruncated {
							break
						}
						input.StartRecordName = out.NextRecordName
						input.StartRecordType = out.NextRecordType
						input.StartRecordIdentifier = out.NextRecordIdentifier
					}
				}(zone)
			}

			go func() {
				wg.Wait()
				close(objectsC)
				close(resourcesC)
			}()
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
}
func addManualLambdaFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualMonitoringFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualCdnFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}
func addManualCloudformationFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

func addManualEksFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["ekscluster"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []ekstypes.Cluster

		if !conf.getBoolDefaultTrue("aws.eks.ekscluster.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource eks[ekscluster]")
			return resources, objects, nil
		}

		paginator := eks.NewListClustersPaginator(conf.APIs.Eks, &eks.ListClustersInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, name := range out.Clusters {
				descOut, err := conf.APIs.Eks.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: &name})
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

	funcs["eksnodegroup"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []ekstypes.Nodegroup

		if !conf.getBoolDefaultTrue("aws.eks.eksnodegroup.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource eks[eksnodegroup]")
			return resources, objects, nil
		}

		clusterPaginator := eks.NewListClustersPaginator(conf.APIs.Eks, &eks.ListClustersInput{})
		for clusterPaginator.HasMorePages() {
			cOut, err := clusterPaginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, clusterName := range cOut.Clusters {
				ngPaginator := eks.NewListNodegroupsPaginator(conf.APIs.Eks, &eks.ListNodegroupsInput{ClusterName: &clusterName})
				for ngPaginator.HasMorePages() {
					ngOut, err := ngPaginator.NextPage(ctx)
					if err != nil {
						return resources, objects, err
					}
					for _, ngName := range ngOut.Nodegroups {
						descOut, err := conf.APIs.Eks.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{ClusterName: &clusterName, NodegroupName: &ngName})
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
	funcs["dynamodbtable"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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
	funcs["key"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []kmstypes.KeyMetadata

		if !conf.getBoolDefaultTrue("aws.secretsmanager.key.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource secretsmanager[key]")
			return resources, objects, nil
		}

		paginator := kms.NewListKeysPaginator(conf.APIs.Kms, &kms.ListKeysInput{})
		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				return resources, objects, err
			}
			for _, keyEntry := range out.Keys {
				descOut, err := conf.APIs.Kms.DescribeKey(ctx, &kms.DescribeKeyInput{KeyId: keyEntry.KeyId})
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
	funcs["apigateway"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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

	funcs["apigatewayroute"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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
				parent := graph.InitResource(cloud.ApiGateway, awssdk.ToString(api.ApiId))
				res.AddRelation(rdf.ChildrenOfRel, parent)
				resources = append(resources, res)
			}
		}
		return resources, objects, nil
	}

	funcs["apigatewaystage"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
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
				parent := graph.InitResource(cloud.ApiGateway, awssdk.ToString(api.ApiId))
				res.AddRelation(rdf.ChildrenOfRel, parent)
				resources = append(resources, res)
			}
		}
		return resources, objects, nil
	}
}

func addManualSsmFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

func addManualCloudtrailFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
}

func addManualEfsFetchFuncs(conf *Config, funcs map[string]fetch.Func) {
	funcs["mounttarget"] = func(ctx context.Context, cache fetch.Cache) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []efstypes.MountTargetDescription

		if !conf.getBoolDefaultTrue("aws.efs.mounttarget.sync") && !getBoolFromContext(ctx, "force") {
			conf.Log.Verbose("sync: *disabled* for resource efs[mounttarget]")
			return resources, objects, nil
		}

		fsOut, err := conf.APIs.Efs.DescribeFileSystems(ctx, &efs.DescribeFileSystemsInput{})
		if err != nil {
			return resources, objects, err
		}
		for _, fs := range fsOut.FileSystems {
			mtOut, err := conf.APIs.Efs.DescribeMountTargets(ctx, &efs.DescribeMountTargetsInput{FileSystemId: fs.FileSystemId})
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
