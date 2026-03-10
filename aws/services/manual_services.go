/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsservices

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/wallix/awless/cloud"
)

func GetCloudServicesForAPIs(apis ...string) (services []cloud.Service) {
	unique := make(map[string]struct{})
	for _, api := range apis {
		if name, ok := ServicePerAPI[api]; ok {
			if service, exists := cloud.ServiceRegistry[name]; exists {
				if _, done := unique[name]; !done {
					unique[name] = struct{}{}
					services = append(services, service)
				}
			}
		}
	}
	return
}

func GetCloudServicesForTypes(types ...string) (services []cloud.Service) {
	unique := make(map[string]struct{})
	for _, typ := range types {
		if name, ok := ServicePerResourceType[typ]; ok {
			if service, exists := cloud.ServiceRegistry[name]; exists {
				if _, done := unique[name]; !done {
					unique[name] = struct{}{}
					services = append(services, service)
				}
			}
		}
	}
	return
}

func ResourceTypesPerServiceName() map[string][]string {
	out := make(map[string][]string)
	for rT, s := range ServicePerResourceType {
		out[s] = append(out[s], rT)
	}
	return out
}

var arnResourceInfoRegex = regexp.MustCompile(`(root)|([\w-.]*)/([\w-./]*)`)

type Identity struct {
	Account, Arn, UserId, ResourceType, ResourcePath, Resource string
}

func (i *Identity) IsRoot() bool {
	return i.Resource == "root"
}

func (i *Identity) IsUserType() bool {
	return i.ResourceType == "user"
}

func (s *Access) GetIdentity() (*Identity, error) {
	resp, err := s.StsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	ident := &Identity{
		Account: aws.ToString(resp.Account),
		Arn:     aws.ToString(resp.Arn),
		UserId:  aws.ToString(resp.UserId),
	}

	splits := strings.Split(ident.Arn, ":")
	if l := len(splits); l > 0 {
		ident.ResourcePath = splits[l-1]
		matches := arnResourceInfoRegex.FindStringSubmatch(ident.ResourcePath)
		if len(matches) == 4 {
			if matches[1] == "root" {
				ident.Resource = "root"
				ident.ResourceType = "user"
			} else {
				ident.ResourceType = matches[2]
				ident.Resource = matches[3]
			}
		}
	}

	return ident, nil
}

type UserPolicies struct {
	Username string
	Inlined  []string
	Attached []string
	ByGroup  map[string][]string
}

func (s *Access) GetUserPolicies(username string) (*UserPolicies, error) {
	var wg sync.WaitGroup

	all := &UserPolicies{
		Username: username,
		ByGroup:  make(map[string][]string),
	}

	errc := make(chan error, 4)

	wg.Add(1)
	go func() {
		defer wg.Done()
		policies, err := s.IamClient.ListUserPolicies(context.Background(), &iam.ListUserPoliciesInput{
			UserName: aws.String(username),
		})
		if err != nil {
			errc <- err
			return
		}

		all.Inlined = append(all.Inlined, policies.PolicyNames...)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		attached, err := s.IamClient.ListAttachedUserPolicies(context.Background(), &iam.ListAttachedUserPoliciesInput{
			UserName: aws.String(username),
		})
		if err != nil {
			errc <- err
			return
		}

		for _, pol := range attached.AttachedPolicies {
			all.Attached = append(all.Attached, aws.ToString(pol.PolicyName))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		groups, err := s.IamClient.ListGroupsForUser(context.Background(), &iam.ListGroupsForUserInput{
			UserName: aws.String(username),
		})
		if err != nil {
			errc <- err
			return
		}

		type result struct {
			group, policy string
		}
		resultC := make(chan result)
		var wgg sync.WaitGroup
		for _, group := range groups.Groups {
			wgg.Add(1)
			go func(name string) {
				defer wgg.Done()

				output, err := s.IamClient.ListAttachedGroupPolicies(context.Background(), &iam.ListAttachedGroupPoliciesInput{
					GroupName: aws.String(name),
				})
				if err != nil {
					errc <- err
					return
				}
				for _, pol := range output.AttachedPolicies {
					resultC <- result{group: name, policy: aws.ToString(pol.PolicyName)}
				}
			}(aws.ToString(group.GroupName))
		}

		go func() {
			wgg.Wait()
			close(resultC)
		}()

		for res := range resultC {
			all.ByGroup[res.group] = append(all.ByGroup[res.group], res.policy)
		}
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for e := range errc {
		return all, e
	}

	return all, nil
}
