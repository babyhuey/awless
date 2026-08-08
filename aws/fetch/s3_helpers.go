package awsfetch

import (
	"context"
	"sync"

	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/errgroup"

	awsconv "github.com/bootswithdefer/awless/aws/conv"
	"github.com/bootswithdefer/awless/cloud/rdf"
	"github.com/bootswithdefer/awless/fetch"
	"github.com/bootswithdefer/awless/graph"
)

func forEachBucketParallel(ctx context.Context, cache fetch.Cache, api *s3.Client, f func(b s3types.Bucket) error) error {
	var buckets []s3types.Bucket

	if val, e := cache.Get("getBucketsPerRegion", func() (any, error) {
		return getBucketsPerRegion(ctx, api)
	}); e != nil {
		return e
	} else if v, ok := val.([]s3types.Bucket); ok {
		buckets = v
	}

	// Bounded: one goroutine per bucket previously meant thousands of
	// simultaneous AWS calls on large accounts, causing throttling. errgroup
	// also removes the leak the old unbuffered error channel had, where
	// returning on the first error left the remaining senders blocked forever.
	g := new(errgroup.Group)
	g.SetLimit(maxParallelAWSCalls)

	for _, output := range buckets {
		b := output
		g.Go(func() error { return f(b) })
	}

	return g.Wait()
}

func fetchObjectsForBucket(ctx context.Context, api *s3.Client, bucket s3types.Bucket, resourcesC chan<- *graph.Resource) error {
	objectc := make(chan []s3types.Object)
	errc := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		paginator := s3.NewListObjectsV2Paginator(api, &s3.ListObjectsV2Input{Bucket: bucket.Name})
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				errc <- err
				return
			}
			objectc <- page.Contents
		}
	}()

	processObjects := func(objs []s3types.Object) {
		for _, output := range objs {
			res, err := awsconv.NewResource(output)
			if err != nil {
				errc <- err
				return
			}
			res.SetProperty("Bucket", awssdk.ToString(bucket.Name))
			resourcesC <- res
			parent, err := awsconv.InitResource(bucket)
			if err != nil {
				errc <- err
				return
			}
			res.AddRelation(rdf.ChildrenOfRel, parent)
			resourcesC <- parent
		}
	}

	go func() {
		wg.Wait()
		close(objectc)
		close(errc)
	}()

	for {
		select {
		case err := <-errc:
			return err
		case objects, ok := <-objectc:
			if !ok {
				return nil
			}
			processObjects(objects)
		}
	}
}

func getBucketsPerRegion(ctx context.Context, api *s3.Client) ([]s3types.Bucket, error) {
	var buckets []s3types.Bucket

	out, err := api.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return buckets, err
	}

	var userBucketName string
	var hasBucketFilter bool
	if id, hasID := getUserFiltersFromContext(ctx)["id"]; hasID {
		userBucketName = id
		hasBucketFilter = true
	} else if buck, hasBucket := getUserFiltersFromContext(ctx)["bucket"]; hasBucket {
		userBucketName = buck
		hasBucketFilter = true
	}

	if hasBucketFilter {
		for _, b := range out.Buckets {
			if strings.Contains(strings.ToLower(awssdk.ToString(b.Name)), strings.ToLower(userBucketName)) {
				buckets = append(buckets, b)
			}
		}
	} else {
		buckets = out.Buckets
	}

	// Bounded, and leak-free: the previous version used unbuffered bucketc and
	// errc channels with a consumer that returned on the first error, leaving
	// every other in-flight goroutine blocked on send forever.
	region, _ := ctx.Value("region").(string)

	var (
		mu              sync.Mutex
		bucketsInRegion []s3types.Bucket
	)

	g := new(errgroup.Group)
	g.SetLimit(maxParallelAWSCalls)

	for _, bucket := range buckets {
		b := bucket
		g.Go(func() error {
			loc, err := api.GetBucketLocation(ctx, &s3.GetBucketLocationInput{Bucket: b.Name})
			if err != nil {
				return err
			}

			locConstraint := string(loc.LocationConstraint)
			// An empty location constraint means us-east-1.
			inRegion := locConstraint == region || (locConstraint == "" && region == "us-east-1")
			if inRegion {
				mu.Lock()
				bucketsInRegion = append(bucketsInRegion, b)
				mu.Unlock()
			}
			return nil
		})
	}

	// Partial results are returned alongside an error, as before.
	err = g.Wait()
	return bucketsInRegion, err
}

func fetchAndExtractGrantsFn(ctx context.Context, api *s3.Client, bucketName string) ([]*graph.Grant, error) {
	acls, err := api.GetBucketAcl(ctx, &s3.GetBucketAclInput{Bucket: awssdk.String(bucketName)})
	if err != nil {
		return nil, err
	}
	var grants []*graph.Grant
	for _, acl := range acls.Grants {
		displayName := awssdk.ToString(acl.Grantee.DisplayName)
		granteeType := string(acl.Grantee.Type)
		granteeID := awssdk.ToString(acl.Grantee.ID)

		if awssdk.ToString(acl.Grantee.EmailAddress) != "" {
			displayName += "<" + awssdk.ToString(acl.Grantee.EmailAddress) + ">"
		}
		if granteeType == "Group" {
			granteeID += awssdk.ToString(acl.Grantee.URI)
		}
		grant := &graph.Grant{
			Permission: string(acl.Permission),
			Grantee: graph.Grantee{
				GranteeID:          granteeID,
				GranteeType:        granteeType,
				GranteeDisplayName: displayName,
			},
		}
		grants = append(grants, grant)
	}
	return grants, nil
}
