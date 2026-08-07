package awsfetch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/bootswithdefer/awless/fetch"
)

type AccountAuthorizationDetails struct {
	Groups   []iamtypes.GroupDetail
	Policies []iamtypes.ManagedPolicyDetail
	Roles    []iamtypes.RoleDetail
	Users    []iamtypes.UserDetail
}

func getAccountAuthorizationDetails(ctx context.Context, cache fetch.Cache, api *iam.Client) (*AccountAuthorizationDetails, error) {
	var entities []iamtypes.EntityType
	var cacheKey string
	resourceType, ok := fetch.IsFetchingByType(ctx)
	if ok {
		switch resourceType {
		case "user":
			cacheKey = "usersDetails"
			entities = append(entities, iamtypes.EntityTypeUser)
		case "group":
			cacheKey = "groupsDetails"
			entities = append(entities, iamtypes.EntityTypeGroup)
		case "role":
			cacheKey = "rolesDetails"
			entities = append(entities, iamtypes.EntityTypeRole)
		case "policy":
			cacheKey = "policiesDetails"
			entities = append(entities, iamtypes.EntityTypeLocalManagedPolicy, iamtypes.EntityTypeAWSManagedPolicy)
		}
	} else {
		cacheKey = "accountDetails"
		entities = append(entities, iamtypes.EntityTypeUser, iamtypes.EntityTypeGroup, iamtypes.EntityTypeRole)
		entities = append(entities, iamtypes.EntityTypeLocalManagedPolicy, iamtypes.EntityTypeAWSManagedPolicy)
	}

	if val, err := cache.Get(cacheKey, func() (interface{}, error) {
		return fetchAccountAuthorizationDetails(ctx, entities, api)
	}); err != nil {
		return nil, err
	} else if v, ok := val.(*AccountAuthorizationDetails); ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("cannot get account details (val of type %T)", val)
	}
}

func fetchAccountAuthorizationDetails(ctx context.Context, entities []iamtypes.EntityType, api *iam.Client) (*AccountAuthorizationDetails, error) {
	details := new(AccountAuthorizationDetails)
	paginator := iam.NewGetAccountAuthorizationDetailsPaginator(api, &iam.GetAccountAuthorizationDetailsInput{
		Filter: entities,
	})
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return details, err
		}
		details.Users = append(details.Users, out.UserDetailList...)
		details.Groups = append(details.Groups, out.GroupDetailList...)
		details.Roles = append(details.Roles, out.RoleDetailList...)
		details.Policies = append(details.Policies, out.Policies...)
	}

	return details, nil
}
