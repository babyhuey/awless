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

package cloud

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrFetchAccessDenied = errors.New("access denied to cloud resource")

// Resources
const (
	Region string = "region"
	//infra
	Vpc              string = "vpc"
	Subnet           string = "subnet"
	Image            string = "image"
	ImportImageTask  string = "importimagetask"
	SecurityGroup    string = "securitygroup"
	AvailabilityZone string = "availabilityzone"
	Keypair          string = "keypair"
	Volume           string = "volume"
	Instance         string = "instance"
	InstanceProfile  string = "instanceprofile"
	InternetGateway  string = "internetgateway"
	NatGateway       string = "natgateway"
	RouteTable       string = "routetable"
	ElasticIP        string = "elasticip"
	Snapshot         string = "snapshot"
	NetworkInterface string = "networkinterface"
	Certificate      string = "certificate"
	//loadbalancer
	ClassicLoadBalancer string = "classicloadbalancer"
	LoadBalancer        string = "loadbalancer"
	TargetGroup         string = "targetgroup"
	Listener            string = "listener"
	//database
	Database      string = "database"
	DBSubnetGroup string = "dbsubnetgroup"
	//access
	User         string = "user"
	Role         string = "role"
	Group        string = "group"
	Policy       string = "policy"
	AccessKey    string = "accesskey"
	LoginProfile string = "loginprofile"
	MFADevice    string = "mfadevice"
	//storage
	Bucket   string = "bucket"
	S3Object string = "s3object"
	ACL      string = "storageacl"
	//notification
	Subscription string = "subscription"
	Topic        string = "topic"
	//queue
	Queue string = "queue"
	//dns
	Zone   string = "zone"
	Record string = "record"
	//lambda
	Function string = "function"
	//autoscaling
	LaunchConfiguration string = "launchconfiguration"
	ScalingGroup        string = "scalinggroup"
	ScalingPolicy       string = "scalingpolicy"
	//monitoring
	Metric string = "metric"
	Alarm  string = "alarm"
	//cdn
	Distribution string = "distribution"
	//cloudformation
	Stack string = "stack"
	//container
	Repository        string = "repository"
	Registry          string = "registry"
	ContainerCluster  string = "containercluster"
	ContainerService  string = "containerservice"
	ContainerTask     string = "containertask"
	Container         string = "container"
	ContainerInstance string = "containerinstance"
	//application autoscaling
	AppScalingTarget string = "appscalingtarget"
	AppScalingPolicy string = "appscalingpolicy"
	//eks
	EKSCluster   string = "ekscluster"
	EKSNodeGroup string = "eksnodegroup"
	//dynamodb
	DynamoDBTable string = "dynamodbtable"
	//secrets & encryption
	Secret string = "secret"
	Key    string = "key"
	//api gateway
	APIGateway      string = "apigateway"
	APIGatewayRoute string = "apigatewayroute"
	APIGatewayStage string = "apigatewaystage"
	//ssm
	SSMParameter string = "ssmparameter"
	//efs
	FileSystem  string = "filesystem"
	MountTarget string = "mounttarget"
	//cloudtrail
	Trail string = "trail"
	//cloudwatchlogs
	LogGroup string = "loggroup"
	//elasticache
	CacheCluster     string = "cachecluster"
	ReplicationGroup string = "replicationgroup"
	CacheSubnetGroup string = "cachesubnetgroup"
	//eventbridge
	EventBus  string = "eventbus"
	EventRule string = "eventrule"
	//stepfunctions
	StateMachine string = "statemachine"
	//wafv2
	WebACL    string = "webacl"
	IPSet     string = "ipset"
	RuleGroup string = "rulegroup"
)

type Service interface {
	Region() string
	Profile() string
	Name() string
	ResourceTypes() []string
	IsSyncDisabled() bool
	Fetch(context.Context) (GraphAPI, error)
	FetchByType(context.Context, string) (GraphAPI, error)
}

type Services []Service

func (srvs Services) Names() (names []string) {
	for _, srv := range srvs {
		names = append(names, srv.Name())
	}
	return
}

var ServiceRegistry = make(map[string]Service)

func AllServices() (out []Service) {
	for _, srv := range ServiceRegistry {
		out = append(out, srv)
	}
	return
}

func GetServiceForType(t string) (Service, error) {
	for _, srv := range ServiceRegistry {
		for _, typ := range srv.ResourceTypes() {
			if typ == t {
				return srv, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot find cloud service for resource type %s", t)
}

// irregularPlurals holds the resource types whose plural is not formed by the rules
// below. A name ending in "s" needs "es", but that cannot be inverted by suffix alone:
// "eventbuses" and "databases" have the same ending and different singulars
// ("eventbus" and "database"), so the mapping is listed rather than inferred.
var irregularPlurals = map[string]string{
	EventBus: "eventbuses",
}

var irregularSingulars = func() map[string]string {
	m := make(map[string]string, len(irregularPlurals))
	for singular, plural := range irregularPlurals {
		m[plural] = singular
	}
	return m
}()

func PluralizeResource(singular string) string {
	if plural, ok := irregularPlurals[singular]; ok {
		return plural
	}
	if strings.HasSuffix(singular, "cy") || strings.HasSuffix(singular, "ry") {
		return strings.TrimSuffix(singular, "y") + "ies"
	}
	return singular + "s"
}

func SingularizeResource(plural string) string {
	if singular, ok := irregularSingulars[plural]; ok {
		return singular
	}
	if strings.HasSuffix(plural, "ies") {
		return strings.TrimSuffix(plural, "ies") + "y"
	}
	return strings.TrimSuffix(plural, "s")
}
