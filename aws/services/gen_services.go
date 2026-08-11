// Auto generated implementation for the AWS cloud service

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

// DO NOT EDIT - This file was automatically generated with go generate

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	acm "github.com/aws/aws-sdk-go-v2/service/acm"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	apigatewayv2 "github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	apigatewayv2types "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	applicationautoscaling "github.com/aws/aws-sdk-go-v2/service/applicationautoscaling"
	autoscaling "github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cloudformation "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	cloudfront "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	cloudtrail "github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailtypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	cloudwatch "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	cloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ecr "github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
	ecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	efs "github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	eks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	elasticache "github.com/aws/aws-sdk-go-v2/service/elasticache"
	elasticachetypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	eventbridge "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	iam "github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	kms "github.com/aws/aws-sdk-go-v2/service/kms"
	kmstypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	lambda "github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	rds "github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	route53 "github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	secretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	sfn "github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	sns "github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	ssm "github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	sts "github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	tstore "github.com/bootswithdefer/triplestore"

	awsfetch "github.com/bootswithdefer/awless/aws/fetch"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/fetch"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/logger"
)

const accessDenied = "Access Denied"

var ServiceNames = []string{
	"infra",
	"access",
	"storage",
	"messaging",
	"dns",
	"lambda",
	"monitoring",
	"cdn",
	"cloudformation",
	"eks",
	"dynamodb",
	"secretsmanager",
	"apigateway",
	"ssm",
	"efs",
	"cloudtrail",
	"cloudwatchlogs",
	"elasticache",
	"eventbridge",
	"stepfunctions",
}

var ResourceTypes = []string{
	"instance",
	"subnet",
	"vpc",
	"keypair",
	"securitygroup",
	"volume",
	"internetgateway",
	"natgateway",
	"routetable",
	"availabilityzone",
	"image",
	"importimagetask",
	"elasticip",
	"snapshot",
	"networkinterface",
	"classicloadbalancer",
	"loadbalancer",
	"targetgroup",
	"listener",
	"database",
	"dbsubnetgroup",
	"launchconfiguration",
	"scalinggroup",
	"scalingpolicy",
	"repository",
	"containercluster",
	"containertask",
	"container",
	"containerinstance",
	"certificate",
	"user",
	"group",
	"role",
	"policy",
	"accesskey",
	"instanceprofile",
	"mfadevice",
	"bucket",
	"s3object",
	"subscription",
	"topic",
	"queue",
	"zone",
	"record",
	"function",
	"metric",
	"alarm",
	"distribution",
	"stack",
	"ekscluster",
	"eksnodegroup",
	"dynamodbtable",
	"secret",
	"key",
	"apigateway",
	"apigatewayroute",
	"apigatewaystage",
	"ssmparameter",
	"filesystem",
	"mounttarget",
	"trail",
	"loggroup",
	"cachecluster",
	"replicationgroup",
	"cachesubnetgroup",
	"eventbus",
	"eventrule",
	"statemachine",
}

var ServicePerAPI = map[string]string{
	"ec2":                    "infra",
	"elbv2":                  "infra",
	"elb":                    "infra",
	"rds":                    "infra",
	"autoscaling":            "infra",
	"ecr":                    "infra",
	"ecs":                    "infra",
	"applicationautoscaling": "infra",
	"acm":                    "infra",
	"iam":                    "access",
	"sts":                    "access",
	"s3":                     "storage",
	"sns":                    "messaging",
	"sqs":                    "messaging",
	"route53":                "dns",
	"lambda":                 "lambda",
	"cloudwatch":             "monitoring",
	"cloudfront":             "cdn",
	"cloudformation":         "cloudformation",
	"eks":                    "eks",
	"dynamodb":               "dynamodb",
	"secretsmanager":         "secretsmanager",
	"kms":                    "secretsmanager",
	"apigatewayv2":           "apigateway",
	"ssm":                    "ssm",
	"efs":                    "efs",
	"cloudtrail":             "cloudtrail",
	"cloudwatchlogs":         "cloudwatchlogs",
	"elasticache":            "elasticache",
	"eventbridge":            "eventbridge",
	"sfn":                    "stepfunctions",
}

var ServicePerResourceType = map[string]string{
	"instance":            "infra",
	"subnet":              "infra",
	"vpc":                 "infra",
	"keypair":             "infra",
	"securitygroup":       "infra",
	"volume":              "infra",
	"internetgateway":     "infra",
	"natgateway":          "infra",
	"routetable":          "infra",
	"availabilityzone":    "infra",
	"image":               "infra",
	"importimagetask":     "infra",
	"elasticip":           "infra",
	"snapshot":            "infra",
	"networkinterface":    "infra",
	"classicloadbalancer": "infra",
	"loadbalancer":        "infra",
	"targetgroup":         "infra",
	"listener":            "infra",
	"database":            "infra",
	"dbsubnetgroup":       "infra",
	"launchconfiguration": "infra",
	"scalinggroup":        "infra",
	"scalingpolicy":       "infra",
	"repository":          "infra",
	"containercluster":    "infra",
	"containertask":       "infra",
	"container":           "infra",
	"containerinstance":   "infra",
	"certificate":         "infra",
	"user":                "access",
	"group":               "access",
	"role":                "access",
	"policy":              "access",
	"accesskey":           "access",
	"instanceprofile":     "access",
	"mfadevice":           "access",
	"bucket":              "storage",
	"s3object":            "storage",
	"subscription":        "messaging",
	"topic":               "messaging",
	"queue":               "messaging",
	"zone":                "dns",
	"record":              "dns",
	"function":            "lambda",
	"metric":              "monitoring",
	"alarm":               "monitoring",
	"distribution":        "cdn",
	"stack":               "cloudformation",
	"ekscluster":          "eks",
	"eksnodegroup":        "eks",
	"dynamodbtable":       "dynamodb",
	"secret":              "secretsmanager",
	"key":                 "secretsmanager",
	"apigateway":          "apigateway",
	"apigatewayroute":     "apigateway",
	"apigatewaystage":     "apigateway",
	"ssmparameter":        "ssm",
	"filesystem":          "efs",
	"mounttarget":         "efs",
	"trail":               "cloudtrail",
	"loggroup":            "cloudwatchlogs",
	"cachecluster":        "elasticache",
	"replicationgroup":    "elasticache",
	"cachesubnetgroup":    "elasticache",
	"eventbus":            "eventbridge",
	"eventrule":           "eventbridge",
	"statemachine":        "stepfunctions",
}

var APIPerResourceType = map[string]string{
	"instance":            "ec2",
	"subnet":              "ec2",
	"vpc":                 "ec2",
	"keypair":             "ec2",
	"securitygroup":       "ec2",
	"volume":              "ec2",
	"internetgateway":     "ec2",
	"natgateway":          "ec2",
	"routetable":          "ec2",
	"availabilityzone":    "ec2",
	"image":               "ec2",
	"importimagetask":     "ec2",
	"elasticip":           "ec2",
	"snapshot":            "ec2",
	"networkinterface":    "ec2",
	"classicloadbalancer": "elb",
	"loadbalancer":        "elbv2",
	"targetgroup":         "elbv2",
	"listener":            "elbv2",
	"database":            "rds",
	"dbsubnetgroup":       "rds",
	"launchconfiguration": "autoscaling",
	"scalinggroup":        "autoscaling",
	"scalingpolicy":       "autoscaling",
	"repository":          "ecr",
	"containercluster":    "ecs",
	"containertask":       "ecs",
	"container":           "ecs",
	"containerinstance":   "ecs",
	"certificate":         "acm",
	"user":                "iam",
	"group":               "iam",
	"role":                "iam",
	"policy":              "iam",
	"accesskey":           "iam",
	"instanceprofile":     "iam",
	"mfadevice":           "iam",
	"bucket":              "s3",
	"s3object":            "s3",
	"subscription":        "sns",
	"topic":               "sns",
	"queue":               "sqs",
	"zone":                "route53",
	"record":              "route53",
	"function":            "lambda",
	"metric":              "cloudwatch",
	"alarm":               "cloudwatch",
	"distribution":        "cloudfront",
	"stack":               "cloudformation",
	"ekscluster":          "eks",
	"eksnodegroup":        "eks",
	"dynamodbtable":       "dynamodb",
	"secret":              "secretsmanager",
	"key":                 "kms",
	"apigateway":          "apigatewayv2",
	"apigatewayroute":     "apigatewayv2",
	"apigatewaystage":     "apigatewayv2",
	"ssmparameter":        "ssm",
	"filesystem":          "efs",
	"mounttarget":         "efs",
	"trail":               "cloudtrail",
	"loggroup":            "cloudwatchlogs",
	"cachecluster":        "elasticache",
	"replicationgroup":    "elasticache",
	"cachesubnetgroup":    "elasticache",
	"eventbus":            "eventbridge",
	"eventrule":           "eventbridge",
	"statemachine":        "sfn",
}

type Infra struct {
	fetcher                      fetch.Fetcher
	region, profile              string
	config                       map[string]any
	log                          *logger.Logger
	Ec2Client                    *ec2.Client
	Elbv2Client                  *elbv2.Client
	ElbClient                    *elb.Client
	RDSClient                    *rds.Client
	AutoscalingClient            *autoscaling.Client
	ECRClient                    *ecr.Client
	ECSClient                    *ecs.Client
	ApplicationautoscalingClient *applicationautoscaling.Client
	ACMClient                    *acm.Client
}

func NewInfra(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	ec2Client := ec2.NewFromConfig(cfg)
	elbv2Client := elbv2.NewFromConfig(cfg)
	elbClient := elb.NewFromConfig(cfg)
	rdsClient := rds.NewFromConfig(cfg)
	autoscalingClient := autoscaling.NewFromConfig(cfg)
	ecrClient := ecr.NewFromConfig(cfg)
	ecsClient := ecs.NewFromConfig(cfg)
	applicationautoscalingClient := applicationautoscaling.NewFromConfig(cfg)
	acmClient := acm.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		ec2Client,
		elbv2Client,
		elbClient,
		rdsClient,
		autoscalingClient,
		ecrClient,
		ecsClient,
		applicationautoscalingClient,
		acmClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Infra{
		Ec2Client:                    ec2Client,
		Elbv2Client:                  elbv2Client,
		ElbClient:                    elbClient,
		RDSClient:                    rdsClient,
		AutoscalingClient:            autoscalingClient,
		ECRClient:                    ecrClient,
		ECSClient:                    ecsClient,
		ApplicationautoscalingClient: applicationautoscalingClient,
		ACMClient:                    acmClient,
		fetcher:                      fetch.NewFetcher(awsfetch.BuildInfraFetchFuncs(fetchConfig)),
		config:                       extraConf,
		region:                       region,
		profile:                      profile,
		log:                          log,
	}
}

func (s *Infra) Name() string {
	return "infra"
}

func (s *Infra) Region() string {
	return s.region
}

func (s *Infra) Profile() string {
	return s.profile
}

func (s *Infra) ResourceTypes() []string {
	return []string{
		"instance",
		"subnet",
		"vpc",
		"keypair",
		"securitygroup",
		"volume",
		"internetgateway",
		"natgateway",
		"routetable",
		"availabilityzone",
		"image",
		"importimagetask",
		"elasticip",
		"snapshot",
		"networkinterface",
		"classicloadbalancer",
		"loadbalancer",
		"targetgroup",
		"listener",
		"database",
		"dbsubnetgroup",
		"launchconfiguration",
		"scalinggroup",
		"scalingpolicy",
		"repository",
		"containercluster",
		"containertask",
		"container",
		"containerinstance",
		"certificate",
	}
}

func (s *Infra) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.infra.instance.sync", true) {
		list, err := s.fetcher.Get("instance_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.Instance); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.Instance' type from fetch context")
		}
		for _, r := range list.([]ec2types.Instance) {
			for _, fn := range addParentsFns["instance"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.Instance) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.subnet.sync", true) {
		list, err := s.fetcher.Get("subnet_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.Subnet); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.Subnet' type from fetch context")
		}
		for _, r := range list.([]ec2types.Subnet) {
			for _, fn := range addParentsFns["subnet"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.Subnet) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.vpc.sync", true) {
		list, err := s.fetcher.Get("vpc_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.Vpc); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.Vpc' type from fetch context")
		}
		for _, r := range list.([]ec2types.Vpc) {
			for _, fn := range addParentsFns["vpc"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.Vpc) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.keypair.sync", true) {
		list, err := s.fetcher.Get("keypair_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.KeyPairInfo); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.KeyPairInfo' type from fetch context")
		}
		for _, r := range list.([]ec2types.KeyPairInfo) {
			for _, fn := range addParentsFns["keypair"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.KeyPairInfo) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.securitygroup.sync", true) {
		list, err := s.fetcher.Get("securitygroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.SecurityGroup); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.SecurityGroup' type from fetch context")
		}
		for _, r := range list.([]ec2types.SecurityGroup) {
			for _, fn := range addParentsFns["securitygroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.SecurityGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.volume.sync", true) {
		list, err := s.fetcher.Get("volume_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.Volume); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.Volume' type from fetch context")
		}
		for _, r := range list.([]ec2types.Volume) {
			for _, fn := range addParentsFns["volume"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.Volume) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.internetgateway.sync", true) {
		list, err := s.fetcher.Get("internetgateway_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.InternetGateway); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.InternetGateway' type from fetch context")
		}
		for _, r := range list.([]ec2types.InternetGateway) {
			for _, fn := range addParentsFns["internetgateway"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.InternetGateway) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.natgateway.sync", true) {
		list, err := s.fetcher.Get("natgateway_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.NatGateway); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.NatGateway' type from fetch context")
		}
		for _, r := range list.([]ec2types.NatGateway) {
			for _, fn := range addParentsFns["natgateway"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.NatGateway) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.routetable.sync", true) {
		list, err := s.fetcher.Get("routetable_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.RouteTable); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.RouteTable' type from fetch context")
		}
		for _, r := range list.([]ec2types.RouteTable) {
			for _, fn := range addParentsFns["routetable"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.RouteTable) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.availabilityzone.sync", true) {
		list, err := s.fetcher.Get("availabilityzone_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.AvailabilityZone); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.AvailabilityZone' type from fetch context")
		}
		for _, r := range list.([]ec2types.AvailabilityZone) {
			for _, fn := range addParentsFns["availabilityzone"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.AvailabilityZone) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.image.sync", true) {
		list, err := s.fetcher.Get("image_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.Image); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.Image' type from fetch context")
		}
		for _, r := range list.([]ec2types.Image) {
			for _, fn := range addParentsFns["image"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.Image) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.importimagetask.sync", true) {
		list, err := s.fetcher.Get("importimagetask_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.ImportImageTask); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.ImportImageTask' type from fetch context")
		}
		for _, r := range list.([]ec2types.ImportImageTask) {
			for _, fn := range addParentsFns["importimagetask"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.ImportImageTask) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.elasticip.sync", true) {
		list, err := s.fetcher.Get("elasticip_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.Address); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.Address' type from fetch context")
		}
		for _, r := range list.([]ec2types.Address) {
			for _, fn := range addParentsFns["elasticip"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.Address) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.snapshot.sync", true) {
		list, err := s.fetcher.Get("snapshot_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.Snapshot); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.Snapshot' type from fetch context")
		}
		for _, r := range list.([]ec2types.Snapshot) {
			for _, fn := range addParentsFns["snapshot"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.Snapshot) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.networkinterface.sync", true) {
		list, err := s.fetcher.Get("networkinterface_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ec2types.NetworkInterface); !ok {
			return gph, errors.New("cannot cast to '[]ec2types.NetworkInterface' type from fetch context")
		}
		for _, r := range list.([]ec2types.NetworkInterface) {
			for _, fn := range addParentsFns["networkinterface"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ec2types.NetworkInterface) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.classicloadbalancer.sync", true) {
		list, err := s.fetcher.Get("classicloadbalancer_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]elbtypes.LoadBalancerDescription); !ok {
			return gph, errors.New("cannot cast to '[]elbtypes.LoadBalancerDescription' type from fetch context")
		}
		for _, r := range list.([]elbtypes.LoadBalancerDescription) {
			for _, fn := range addParentsFns["classicloadbalancer"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elbtypes.LoadBalancerDescription) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.loadbalancer.sync", true) {
		list, err := s.fetcher.Get("loadbalancer_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]elbv2types.LoadBalancer); !ok {
			return gph, errors.New("cannot cast to '[]elbv2types.LoadBalancer' type from fetch context")
		}
		for _, r := range list.([]elbv2types.LoadBalancer) {
			for _, fn := range addParentsFns["loadbalancer"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elbv2types.LoadBalancer) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.targetgroup.sync", true) {
		list, err := s.fetcher.Get("targetgroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]elbv2types.TargetGroup); !ok {
			return gph, errors.New("cannot cast to '[]elbv2types.TargetGroup' type from fetch context")
		}
		for _, r := range list.([]elbv2types.TargetGroup) {
			for _, fn := range addParentsFns["targetgroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elbv2types.TargetGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.listener.sync", true) {
		list, err := s.fetcher.Get("listener_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]elbv2types.Listener); !ok {
			return gph, errors.New("cannot cast to '[]elbv2types.Listener' type from fetch context")
		}
		for _, r := range list.([]elbv2types.Listener) {
			for _, fn := range addParentsFns["listener"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elbv2types.Listener) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.database.sync", true) {
		list, err := s.fetcher.Get("database_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]rdstypes.DBInstance); !ok {
			return gph, errors.New("cannot cast to '[]rdstypes.DBInstance' type from fetch context")
		}
		for _, r := range list.([]rdstypes.DBInstance) {
			for _, fn := range addParentsFns["database"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *rdstypes.DBInstance) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.dbsubnetgroup.sync", true) {
		list, err := s.fetcher.Get("dbsubnetgroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]rdstypes.DBSubnetGroup); !ok {
			return gph, errors.New("cannot cast to '[]rdstypes.DBSubnetGroup' type from fetch context")
		}
		for _, r := range list.([]rdstypes.DBSubnetGroup) {
			for _, fn := range addParentsFns["dbsubnetgroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *rdstypes.DBSubnetGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.launchconfiguration.sync", true) {
		list, err := s.fetcher.Get("launchconfiguration_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]autoscalingtypes.LaunchConfiguration); !ok {
			return gph, errors.New("cannot cast to '[]autoscalingtypes.LaunchConfiguration' type from fetch context")
		}
		for _, r := range list.([]autoscalingtypes.LaunchConfiguration) {
			for _, fn := range addParentsFns["launchconfiguration"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *autoscalingtypes.LaunchConfiguration) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.scalinggroup.sync", true) {
		list, err := s.fetcher.Get("scalinggroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]autoscalingtypes.AutoScalingGroup); !ok {
			return gph, errors.New("cannot cast to '[]autoscalingtypes.AutoScalingGroup' type from fetch context")
		}
		for _, r := range list.([]autoscalingtypes.AutoScalingGroup) {
			for _, fn := range addParentsFns["scalinggroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *autoscalingtypes.AutoScalingGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.scalingpolicy.sync", true) {
		list, err := s.fetcher.Get("scalingpolicy_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]autoscalingtypes.ScalingPolicy); !ok {
			return gph, errors.New("cannot cast to '[]autoscalingtypes.ScalingPolicy' type from fetch context")
		}
		for _, r := range list.([]autoscalingtypes.ScalingPolicy) {
			for _, fn := range addParentsFns["scalingpolicy"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *autoscalingtypes.ScalingPolicy) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.repository.sync", true) {
		list, err := s.fetcher.Get("repository_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ecrtypes.Repository); !ok {
			return gph, errors.New("cannot cast to '[]ecrtypes.Repository' type from fetch context")
		}
		for _, r := range list.([]ecrtypes.Repository) {
			for _, fn := range addParentsFns["repository"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecrtypes.Repository) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.containercluster.sync", true) {
		list, err := s.fetcher.Get("containercluster_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ecstypes.Cluster); !ok {
			return gph, errors.New("cannot cast to '[]ecstypes.Cluster' type from fetch context")
		}
		for _, r := range list.([]ecstypes.Cluster) {
			for _, fn := range addParentsFns["containercluster"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecstypes.Cluster) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.containertask.sync", true) {
		list, err := s.fetcher.Get("containertask_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ecstypes.TaskDefinition); !ok {
			return gph, errors.New("cannot cast to '[]ecstypes.TaskDefinition' type from fetch context")
		}
		for _, r := range list.([]ecstypes.TaskDefinition) {
			for _, fn := range addParentsFns["containertask"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecstypes.TaskDefinition) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.container.sync", true) {
		list, err := s.fetcher.Get("container_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ecstypes.Container); !ok {
			return gph, errors.New("cannot cast to '[]ecstypes.Container' type from fetch context")
		}
		for _, r := range list.([]ecstypes.Container) {
			for _, fn := range addParentsFns["container"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecstypes.Container) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.containerinstance.sync", true) {
		list, err := s.fetcher.Get("containerinstance_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ecstypes.ContainerInstance); !ok {
			return gph, errors.New("cannot cast to '[]ecstypes.ContainerInstance' type from fetch context")
		}
		for _, r := range list.([]ecstypes.ContainerInstance) {
			for _, fn := range addParentsFns["containerinstance"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ecstypes.ContainerInstance) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.infra.certificate.sync", true) {
		list, err := s.fetcher.Get("certificate_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]acmtypes.CertificateSummary); !ok {
			return gph, errors.New("cannot cast to '[]acmtypes.CertificateSummary' type from fetch context")
		}
		for _, r := range list.([]acmtypes.CertificateSummary) {
			for _, fn := range addParentsFns["certificate"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *acmtypes.CertificateSummary) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Infra) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Infra) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.infra.sync", true)
}

type Access struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	IAMClient       *iam.Client
	STSClient       *sts.Client
}

func NewAccess(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := "global"
	iamClient := iam.NewFromConfig(cfg)
	stsClient := sts.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		iamClient,
		stsClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Access{
		IAMClient: iamClient,
		STSClient: stsClient,
		fetcher:   fetch.NewFetcher(awsfetch.BuildAccessFetchFuncs(fetchConfig)),
		config:    extraConf,
		region:    region,
		profile:   profile,
		log:       log,
	}
}

func (s *Access) Name() string {
	return "access"
}

func (s *Access) Region() string {
	return s.region
}

func (s *Access) Profile() string {
	return s.profile
}

func (s *Access) ResourceTypes() []string {
	return []string{
		"user",
		"group",
		"role",
		"policy",
		"accesskey",
		"instanceprofile",
		"mfadevice",
	}
}

func (s *Access) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.access.user.sync", true) {
		list, err := s.fetcher.Get("user_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]iamtypes.UserDetail); !ok {
			return gph, errors.New("cannot cast to '[]iamtypes.UserDetail' type from fetch context")
		}
		for _, r := range list.([]iamtypes.UserDetail) {
			for _, fn := range addParentsFns["user"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iamtypes.UserDetail) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.access.group.sync", true) {
		list, err := s.fetcher.Get("group_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]iamtypes.GroupDetail); !ok {
			return gph, errors.New("cannot cast to '[]iamtypes.GroupDetail' type from fetch context")
		}
		for _, r := range list.([]iamtypes.GroupDetail) {
			for _, fn := range addParentsFns["group"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iamtypes.GroupDetail) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.access.role.sync", true) {
		list, err := s.fetcher.Get("role_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]iamtypes.RoleDetail); !ok {
			return gph, errors.New("cannot cast to '[]iamtypes.RoleDetail' type from fetch context")
		}
		for _, r := range list.([]iamtypes.RoleDetail) {
			for _, fn := range addParentsFns["role"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iamtypes.RoleDetail) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.access.policy.sync", true) {
		list, err := s.fetcher.Get("policy_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]iamtypes.Policy); !ok {
			return gph, errors.New("cannot cast to '[]iamtypes.Policy' type from fetch context")
		}
		for _, r := range list.([]iamtypes.Policy) {
			for _, fn := range addParentsFns["policy"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iamtypes.Policy) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.access.accesskey.sync", true) {
		list, err := s.fetcher.Get("accesskey_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]iamtypes.AccessKeyMetadata); !ok {
			return gph, errors.New("cannot cast to '[]iamtypes.AccessKeyMetadata' type from fetch context")
		}
		for _, r := range list.([]iamtypes.AccessKeyMetadata) {
			for _, fn := range addParentsFns["accesskey"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iamtypes.AccessKeyMetadata) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.access.instanceprofile.sync", true) {
		list, err := s.fetcher.Get("instanceprofile_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]iamtypes.InstanceProfile); !ok {
			return gph, errors.New("cannot cast to '[]iamtypes.InstanceProfile' type from fetch context")
		}
		for _, r := range list.([]iamtypes.InstanceProfile) {
			for _, fn := range addParentsFns["instanceprofile"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iamtypes.InstanceProfile) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.access.mfadevice.sync", true) {
		list, err := s.fetcher.Get("mfadevice_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]iamtypes.VirtualMFADevice); !ok {
			return gph, errors.New("cannot cast to '[]iamtypes.VirtualMFADevice' type from fetch context")
		}
		for _, r := range list.([]iamtypes.VirtualMFADevice) {
			for _, fn := range addParentsFns["mfadevice"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *iamtypes.VirtualMFADevice) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Access) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Access) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.access.sync", true)
}

type Storage struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	S3Client        *s3.Client
}

func NewStorage(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	s3Client := s3.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		s3Client,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Storage{
		S3Client: s3Client,
		fetcher:  fetch.NewFetcher(awsfetch.BuildStorageFetchFuncs(fetchConfig)),
		config:   extraConf,
		region:   region,
		profile:  profile,
		log:      log,
	}
}

func (s *Storage) Name() string {
	return "storage"
}

func (s *Storage) Region() string {
	return s.region
}

func (s *Storage) Profile() string {
	return s.profile
}

func (s *Storage) ResourceTypes() []string {
	return []string{
		"bucket",
		"s3object",
	}
}

func (s *Storage) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.storage.bucket.sync", true) {
		list, err := s.fetcher.Get("bucket_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]s3types.Bucket); !ok {
			return gph, errors.New("cannot cast to '[]s3types.Bucket' type from fetch context")
		}
		for _, r := range list.([]s3types.Bucket) {
			for _, fn := range addParentsFns["bucket"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *s3types.Bucket) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.storage.s3object.sync", true) {
		list, err := s.fetcher.Get("s3object_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]s3types.Object); !ok {
			return gph, errors.New("cannot cast to '[]s3types.Object' type from fetch context")
		}
		for _, r := range list.([]s3types.Object) {
			for _, fn := range addParentsFns["s3object"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *s3types.Object) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Storage) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Storage) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.storage.sync", true)
}

type Messaging struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	SNSClient       *sns.Client
	SQSClient       *sqs.Client
}

func NewMessaging(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	snsClient := sns.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		snsClient,
		sqsClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Messaging{
		SNSClient: snsClient,
		SQSClient: sqsClient,
		fetcher:   fetch.NewFetcher(awsfetch.BuildMessagingFetchFuncs(fetchConfig)),
		config:    extraConf,
		region:    region,
		profile:   profile,
		log:       log,
	}
}

func (s *Messaging) Name() string {
	return "messaging"
}

func (s *Messaging) Region() string {
	return s.region
}

func (s *Messaging) Profile() string {
	return s.profile
}

func (s *Messaging) ResourceTypes() []string {
	return []string{
		"subscription",
		"topic",
		"queue",
	}
}

func (s *Messaging) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.messaging.subscription.sync", true) {
		list, err := s.fetcher.Get("subscription_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]snstypes.Subscription); !ok {
			return gph, errors.New("cannot cast to '[]snstypes.Subscription' type from fetch context")
		}
		for _, r := range list.([]snstypes.Subscription) {
			for _, fn := range addParentsFns["subscription"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *snstypes.Subscription) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.messaging.topic.sync", true) {
		list, err := s.fetcher.Get("topic_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]snstypes.Topic); !ok {
			return gph, errors.New("cannot cast to '[]snstypes.Topic' type from fetch context")
		}
		for _, r := range list.([]snstypes.Topic) {
			for _, fn := range addParentsFns["topic"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *snstypes.Topic) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.messaging.queue.sync", true) {
		list, err := s.fetcher.Get("queue_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]string); !ok {
			return gph, errors.New("cannot cast to '[]string' type from fetch context")
		}
		for _, r := range list.([]string) {
			for _, fn := range addParentsFns["queue"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *string) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Messaging) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Messaging) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.messaging.sync", true)
}

type DNS struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	Route53Client   *route53.Client
}

func NewDNS(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := "global"
	route53Client := route53.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		route53Client,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &DNS{
		Route53Client: route53Client,
		fetcher:       fetch.NewFetcher(awsfetch.BuildDNSFetchFuncs(fetchConfig)),
		config:        extraConf,
		region:        region,
		profile:       profile,
		log:           log,
	}
}

func (s *DNS) Name() string {
	return "dns"
}

func (s *DNS) Region() string {
	return s.region
}

func (s *DNS) Profile() string {
	return s.profile
}

func (s *DNS) ResourceTypes() []string {
	return []string{
		"zone",
		"record",
	}
}

func (s *DNS) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.dns.zone.sync", true) {
		list, err := s.fetcher.Get("zone_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]route53types.HostedZone); !ok {
			return gph, errors.New("cannot cast to '[]route53types.HostedZone' type from fetch context")
		}
		for _, r := range list.([]route53types.HostedZone) {
			for _, fn := range addParentsFns["zone"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *route53types.HostedZone) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.dns.record.sync", true) {
		list, err := s.fetcher.Get("record_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]route53types.ResourceRecordSet); !ok {
			return gph, errors.New("cannot cast to '[]route53types.ResourceRecordSet' type from fetch context")
		}
		for _, r := range list.([]route53types.ResourceRecordSet) {
			for _, fn := range addParentsFns["record"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *route53types.ResourceRecordSet) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *DNS) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *DNS) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.dns.sync", true)
}

type Lambda struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	LambdaClient    *lambda.Client
}

func NewLambda(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	lambdaClient := lambda.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		lambdaClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Lambda{
		LambdaClient: lambdaClient,
		fetcher:      fetch.NewFetcher(awsfetch.BuildLambdaFetchFuncs(fetchConfig)),
		config:       extraConf,
		region:       region,
		profile:      profile,
		log:          log,
	}
}

func (s *Lambda) Name() string {
	return "lambda"
}

func (s *Lambda) Region() string {
	return s.region
}

func (s *Lambda) Profile() string {
	return s.profile
}

func (s *Lambda) ResourceTypes() []string {
	return []string{
		"function",
	}
}

func (s *Lambda) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.lambda.function.sync", true) {
		list, err := s.fetcher.Get("function_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]lambdatypes.FunctionConfiguration); !ok {
			return gph, errors.New("cannot cast to '[]lambdatypes.FunctionConfiguration' type from fetch context")
		}
		for _, r := range list.([]lambdatypes.FunctionConfiguration) {
			for _, fn := range addParentsFns["function"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *lambdatypes.FunctionConfiguration) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Lambda) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Lambda) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.lambda.sync", true)
}

type Monitoring struct {
	fetcher          fetch.Fetcher
	region, profile  string
	config           map[string]any
	log              *logger.Logger
	CloudwatchClient *cloudwatch.Client
}

func NewMonitoring(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	cloudwatchClient := cloudwatch.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		cloudwatchClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Monitoring{
		CloudwatchClient: cloudwatchClient,
		fetcher:          fetch.NewFetcher(awsfetch.BuildMonitoringFetchFuncs(fetchConfig)),
		config:           extraConf,
		region:           region,
		profile:          profile,
		log:              log,
	}
}

func (s *Monitoring) Name() string {
	return "monitoring"
}

func (s *Monitoring) Region() string {
	return s.region
}

func (s *Monitoring) Profile() string {
	return s.profile
}

func (s *Monitoring) ResourceTypes() []string {
	return []string{
		"metric",
		"alarm",
	}
}

func (s *Monitoring) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.monitoring.metric.sync", true) {
		list, err := s.fetcher.Get("metric_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]cloudwatchtypes.Metric); !ok {
			return gph, errors.New("cannot cast to '[]cloudwatchtypes.Metric' type from fetch context")
		}
		for _, r := range list.([]cloudwatchtypes.Metric) {
			for _, fn := range addParentsFns["metric"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudwatchtypes.Metric) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.monitoring.alarm.sync", true) {
		list, err := s.fetcher.Get("alarm_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]cloudwatchtypes.MetricAlarm); !ok {
			return gph, errors.New("cannot cast to '[]cloudwatchtypes.MetricAlarm' type from fetch context")
		}
		for _, r := range list.([]cloudwatchtypes.MetricAlarm) {
			for _, fn := range addParentsFns["alarm"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudwatchtypes.MetricAlarm) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Monitoring) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Monitoring) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.monitoring.sync", true)
}

type CDN struct {
	fetcher          fetch.Fetcher
	region, profile  string
	config           map[string]any
	log              *logger.Logger
	CloudfrontClient *cloudfront.Client
}

func NewCDN(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := "global"
	cloudfrontClient := cloudfront.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		cloudfrontClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &CDN{
		CloudfrontClient: cloudfrontClient,
		fetcher:          fetch.NewFetcher(awsfetch.BuildCDNFetchFuncs(fetchConfig)),
		config:           extraConf,
		region:           region,
		profile:          profile,
		log:              log,
	}
}

func (s *CDN) Name() string {
	return "cdn"
}

func (s *CDN) Region() string {
	return s.region
}

func (s *CDN) Profile() string {
	return s.profile
}

func (s *CDN) ResourceTypes() []string {
	return []string{
		"distribution",
	}
}

func (s *CDN) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.cdn.distribution.sync", true) {
		list, err := s.fetcher.Get("distribution_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]cloudfronttypes.DistributionSummary); !ok {
			return gph, errors.New("cannot cast to '[]cloudfronttypes.DistributionSummary' type from fetch context")
		}
		for _, r := range list.([]cloudfronttypes.DistributionSummary) {
			for _, fn := range addParentsFns["distribution"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudfronttypes.DistributionSummary) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *CDN) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *CDN) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.cdn.sync", true)
}

type Cloudformation struct {
	fetcher              fetch.Fetcher
	region, profile      string
	config               map[string]any
	log                  *logger.Logger
	CloudformationClient *cloudformation.Client
}

func NewCloudformation(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	cloudformationClient := cloudformation.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		cloudformationClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Cloudformation{
		CloudformationClient: cloudformationClient,
		fetcher:              fetch.NewFetcher(awsfetch.BuildCloudformationFetchFuncs(fetchConfig)),
		config:               extraConf,
		region:               region,
		profile:              profile,
		log:                  log,
	}
}

func (s *Cloudformation) Name() string {
	return "cloudformation"
}

func (s *Cloudformation) Region() string {
	return s.region
}

func (s *Cloudformation) Profile() string {
	return s.profile
}

func (s *Cloudformation) ResourceTypes() []string {
	return []string{
		"stack",
	}
}

func (s *Cloudformation) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.cloudformation.stack.sync", true) {
		list, err := s.fetcher.Get("stack_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]cloudformationtypes.Stack); !ok {
			return gph, errors.New("cannot cast to '[]cloudformationtypes.Stack' type from fetch context")
		}
		for _, r := range list.([]cloudformationtypes.Stack) {
			for _, fn := range addParentsFns["stack"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudformationtypes.Stack) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Cloudformation) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Cloudformation) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.cloudformation.sync", true)
}

type EKS struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	EKSClient       *eks.Client
}

func NewEKS(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	eksClient := eks.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		eksClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &EKS{
		EKSClient: eksClient,
		fetcher:   fetch.NewFetcher(awsfetch.BuildEKSFetchFuncs(fetchConfig)),
		config:    extraConf,
		region:    region,
		profile:   profile,
		log:       log,
	}
}

func (s *EKS) Name() string {
	return "eks"
}

func (s *EKS) Region() string {
	return s.region
}

func (s *EKS) Profile() string {
	return s.profile
}

func (s *EKS) ResourceTypes() []string {
	return []string{
		"ekscluster",
		"eksnodegroup",
	}
}

func (s *EKS) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.eks.ekscluster.sync", true) {
		list, err := s.fetcher.Get("ekscluster_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ekstypes.Cluster); !ok {
			return gph, errors.New("cannot cast to '[]ekstypes.Cluster' type from fetch context")
		}
		for _, r := range list.([]ekstypes.Cluster) {
			for _, fn := range addParentsFns["ekscluster"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ekstypes.Cluster) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.eks.eksnodegroup.sync", true) {
		list, err := s.fetcher.Get("eksnodegroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ekstypes.Nodegroup); !ok {
			return gph, errors.New("cannot cast to '[]ekstypes.Nodegroup' type from fetch context")
		}
		for _, r := range list.([]ekstypes.Nodegroup) {
			for _, fn := range addParentsFns["eksnodegroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ekstypes.Nodegroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *EKS) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *EKS) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.eks.sync", true)
}

type Dynamodb struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	DynamodbClient  *dynamodb.Client
}

func NewDynamodb(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		dynamodbClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Dynamodb{
		DynamodbClient: dynamodbClient,
		fetcher:        fetch.NewFetcher(awsfetch.BuildDynamodbFetchFuncs(fetchConfig)),
		config:         extraConf,
		region:         region,
		profile:        profile,
		log:            log,
	}
}

func (s *Dynamodb) Name() string {
	return "dynamodb"
}

func (s *Dynamodb) Region() string {
	return s.region
}

func (s *Dynamodb) Profile() string {
	return s.profile
}

func (s *Dynamodb) ResourceTypes() []string {
	return []string{
		"dynamodbtable",
	}
}

func (s *Dynamodb) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.dynamodb.dynamodbtable.sync", true) {
		list, err := s.fetcher.Get("dynamodbtable_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]dynamodbtypes.TableDescription); !ok {
			return gph, errors.New("cannot cast to '[]dynamodbtypes.TableDescription' type from fetch context")
		}
		for _, r := range list.([]dynamodbtypes.TableDescription) {
			for _, fn := range addParentsFns["dynamodbtable"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *dynamodbtypes.TableDescription) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Dynamodb) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Dynamodb) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.dynamodb.sync", true)
}

type Secretsmanager struct {
	fetcher              fetch.Fetcher
	region, profile      string
	config               map[string]any
	log                  *logger.Logger
	SecretsmanagerClient *secretsmanager.Client
	KMSClient            *kms.Client
}

func NewSecretsmanager(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	secretsmanagerClient := secretsmanager.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		secretsmanagerClient,
		kmsClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Secretsmanager{
		SecretsmanagerClient: secretsmanagerClient,
		KMSClient:            kmsClient,
		fetcher:              fetch.NewFetcher(awsfetch.BuildSecretsmanagerFetchFuncs(fetchConfig)),
		config:               extraConf,
		region:               region,
		profile:              profile,
		log:                  log,
	}
}

func (s *Secretsmanager) Name() string {
	return "secretsmanager"
}

func (s *Secretsmanager) Region() string {
	return s.region
}

func (s *Secretsmanager) Profile() string {
	return s.profile
}

func (s *Secretsmanager) ResourceTypes() []string {
	return []string{
		"secret",
		"key",
	}
}

func (s *Secretsmanager) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.secretsmanager.secret.sync", true) {
		list, err := s.fetcher.Get("secret_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]secretsmanagertypes.SecretListEntry); !ok {
			return gph, errors.New("cannot cast to '[]secretsmanagertypes.SecretListEntry' type from fetch context")
		}
		for _, r := range list.([]secretsmanagertypes.SecretListEntry) {
			for _, fn := range addParentsFns["secret"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *secretsmanagertypes.SecretListEntry) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.secretsmanager.key.sync", true) {
		list, err := s.fetcher.Get("key_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]kmstypes.KeyMetadata); !ok {
			return gph, errors.New("cannot cast to '[]kmstypes.KeyMetadata' type from fetch context")
		}
		for _, r := range list.([]kmstypes.KeyMetadata) {
			for _, fn := range addParentsFns["key"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *kmstypes.KeyMetadata) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Secretsmanager) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Secretsmanager) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.secretsmanager.sync", true)
}

type Apigateway struct {
	fetcher            fetch.Fetcher
	region, profile    string
	config             map[string]any
	log                *logger.Logger
	Apigatewayv2Client *apigatewayv2.Client
}

func NewApigateway(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	apigatewayv2Client := apigatewayv2.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		apigatewayv2Client,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Apigateway{
		Apigatewayv2Client: apigatewayv2Client,
		fetcher:            fetch.NewFetcher(awsfetch.BuildApigatewayFetchFuncs(fetchConfig)),
		config:             extraConf,
		region:             region,
		profile:            profile,
		log:                log,
	}
}

func (s *Apigateway) Name() string {
	return "apigateway"
}

func (s *Apigateway) Region() string {
	return s.region
}

func (s *Apigateway) Profile() string {
	return s.profile
}

func (s *Apigateway) ResourceTypes() []string {
	return []string{
		"apigateway",
		"apigatewayroute",
		"apigatewaystage",
	}
}

func (s *Apigateway) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.apigateway.apigateway.sync", true) {
		list, err := s.fetcher.Get("apigateway_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]apigatewayv2types.Api); !ok {
			return gph, errors.New("cannot cast to '[]apigatewayv2types.Api' type from fetch context")
		}
		for _, r := range list.([]apigatewayv2types.Api) {
			for _, fn := range addParentsFns["apigateway"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *apigatewayv2types.Api) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.apigateway.apigatewayroute.sync", true) {
		list, err := s.fetcher.Get("apigatewayroute_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]apigatewayv2types.Route); !ok {
			return gph, errors.New("cannot cast to '[]apigatewayv2types.Route' type from fetch context")
		}
		for _, r := range list.([]apigatewayv2types.Route) {
			for _, fn := range addParentsFns["apigatewayroute"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *apigatewayv2types.Route) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.apigateway.apigatewaystage.sync", true) {
		list, err := s.fetcher.Get("apigatewaystage_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]apigatewayv2types.Stage); !ok {
			return gph, errors.New("cannot cast to '[]apigatewayv2types.Stage' type from fetch context")
		}
		for _, r := range list.([]apigatewayv2types.Stage) {
			for _, fn := range addParentsFns["apigatewaystage"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *apigatewayv2types.Stage) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Apigateway) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Apigateway) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.apigateway.sync", true)
}

type SSM struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	SSMClient       *ssm.Client
}

func NewSSM(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	ssmClient := ssm.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		ssmClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &SSM{
		SSMClient: ssmClient,
		fetcher:   fetch.NewFetcher(awsfetch.BuildSSMFetchFuncs(fetchConfig)),
		config:    extraConf,
		region:    region,
		profile:   profile,
		log:       log,
	}
}

func (s *SSM) Name() string {
	return "ssm"
}

func (s *SSM) Region() string {
	return s.region
}

func (s *SSM) Profile() string {
	return s.profile
}

func (s *SSM) ResourceTypes() []string {
	return []string{
		"ssmparameter",
	}
}

func (s *SSM) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.ssm.ssmparameter.sync", true) {
		list, err := s.fetcher.Get("ssmparameter_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]ssmtypes.ParameterMetadata); !ok {
			return gph, errors.New("cannot cast to '[]ssmtypes.ParameterMetadata' type from fetch context")
		}
		for _, r := range list.([]ssmtypes.ParameterMetadata) {
			for _, fn := range addParentsFns["ssmparameter"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *ssmtypes.ParameterMetadata) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *SSM) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *SSM) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.ssm.sync", true)
}

type EFS struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	EFSClient       *efs.Client
}

func NewEFS(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	efsClient := efs.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		efsClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &EFS{
		EFSClient: efsClient,
		fetcher:   fetch.NewFetcher(awsfetch.BuildEFSFetchFuncs(fetchConfig)),
		config:    extraConf,
		region:    region,
		profile:   profile,
		log:       log,
	}
}

func (s *EFS) Name() string {
	return "efs"
}

func (s *EFS) Region() string {
	return s.region
}

func (s *EFS) Profile() string {
	return s.profile
}

func (s *EFS) ResourceTypes() []string {
	return []string{
		"filesystem",
		"mounttarget",
	}
}

func (s *EFS) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.efs.filesystem.sync", true) {
		list, err := s.fetcher.Get("filesystem_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]efstypes.FileSystemDescription); !ok {
			return gph, errors.New("cannot cast to '[]efstypes.FileSystemDescription' type from fetch context")
		}
		for _, r := range list.([]efstypes.FileSystemDescription) {
			for _, fn := range addParentsFns["filesystem"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *efstypes.FileSystemDescription) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.efs.mounttarget.sync", true) {
		list, err := s.fetcher.Get("mounttarget_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]efstypes.MountTargetDescription); !ok {
			return gph, errors.New("cannot cast to '[]efstypes.MountTargetDescription' type from fetch context")
		}
		for _, r := range list.([]efstypes.MountTargetDescription) {
			for _, fn := range addParentsFns["mounttarget"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *efstypes.MountTargetDescription) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *EFS) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *EFS) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.efs.sync", true)
}

type Cloudtrail struct {
	fetcher          fetch.Fetcher
	region, profile  string
	config           map[string]any
	log              *logger.Logger
	CloudtrailClient *cloudtrail.Client
}

func NewCloudtrail(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	cloudtrailClient := cloudtrail.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		cloudtrailClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Cloudtrail{
		CloudtrailClient: cloudtrailClient,
		fetcher:          fetch.NewFetcher(awsfetch.BuildCloudtrailFetchFuncs(fetchConfig)),
		config:           extraConf,
		region:           region,
		profile:          profile,
		log:              log,
	}
}

func (s *Cloudtrail) Name() string {
	return "cloudtrail"
}

func (s *Cloudtrail) Region() string {
	return s.region
}

func (s *Cloudtrail) Profile() string {
	return s.profile
}

func (s *Cloudtrail) ResourceTypes() []string {
	return []string{
		"trail",
	}
}

func (s *Cloudtrail) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.cloudtrail.trail.sync", true) {
		list, err := s.fetcher.Get("trail_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]cloudtrailtypes.Trail); !ok {
			return gph, errors.New("cannot cast to '[]cloudtrailtypes.Trail' type from fetch context")
		}
		for _, r := range list.([]cloudtrailtypes.Trail) {
			for _, fn := range addParentsFns["trail"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudtrailtypes.Trail) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Cloudtrail) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Cloudtrail) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.cloudtrail.sync", true)
}

type Cloudwatchlogs struct {
	fetcher              fetch.Fetcher
	region, profile      string
	config               map[string]any
	log                  *logger.Logger
	CloudwatchlogsClient *cloudwatchlogs.Client
}

func NewCloudwatchlogs(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	cloudwatchlogsClient := cloudwatchlogs.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		cloudwatchlogsClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Cloudwatchlogs{
		CloudwatchlogsClient: cloudwatchlogsClient,
		fetcher:              fetch.NewFetcher(awsfetch.BuildCloudwatchlogsFetchFuncs(fetchConfig)),
		config:               extraConf,
		region:               region,
		profile:              profile,
		log:                  log,
	}
}

func (s *Cloudwatchlogs) Name() string {
	return "cloudwatchlogs"
}

func (s *Cloudwatchlogs) Region() string {
	return s.region
}

func (s *Cloudwatchlogs) Profile() string {
	return s.profile
}

func (s *Cloudwatchlogs) ResourceTypes() []string {
	return []string{
		"loggroup",
	}
}

func (s *Cloudwatchlogs) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.cloudwatchlogs.loggroup.sync", true) {
		list, err := s.fetcher.Get("loggroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]cloudwatchlogstypes.LogGroup); !ok {
			return gph, errors.New("cannot cast to '[]cloudwatchlogstypes.LogGroup' type from fetch context")
		}
		for _, r := range list.([]cloudwatchlogstypes.LogGroup) {
			for _, fn := range addParentsFns["loggroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *cloudwatchlogstypes.LogGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Cloudwatchlogs) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Cloudwatchlogs) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.cloudwatchlogs.sync", true)
}

type Elasticache struct {
	fetcher           fetch.Fetcher
	region, profile   string
	config            map[string]any
	log               *logger.Logger
	ElasticacheClient *elasticache.Client
}

func NewElasticache(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	elasticacheClient := elasticache.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		elasticacheClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Elasticache{
		ElasticacheClient: elasticacheClient,
		fetcher:           fetch.NewFetcher(awsfetch.BuildElasticacheFetchFuncs(fetchConfig)),
		config:            extraConf,
		region:            region,
		profile:           profile,
		log:               log,
	}
}

func (s *Elasticache) Name() string {
	return "elasticache"
}

func (s *Elasticache) Region() string {
	return s.region
}

func (s *Elasticache) Profile() string {
	return s.profile
}

func (s *Elasticache) ResourceTypes() []string {
	return []string{
		"cachecluster",
		"replicationgroup",
		"cachesubnetgroup",
	}
}

func (s *Elasticache) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.elasticache.cachecluster.sync", true) {
		list, err := s.fetcher.Get("cachecluster_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]elasticachetypes.CacheCluster); !ok {
			return gph, errors.New("cannot cast to '[]elasticachetypes.CacheCluster' type from fetch context")
		}
		for _, r := range list.([]elasticachetypes.CacheCluster) {
			for _, fn := range addParentsFns["cachecluster"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elasticachetypes.CacheCluster) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.elasticache.replicationgroup.sync", true) {
		list, err := s.fetcher.Get("replicationgroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]elasticachetypes.ReplicationGroup); !ok {
			return gph, errors.New("cannot cast to '[]elasticachetypes.ReplicationGroup' type from fetch context")
		}
		for _, r := range list.([]elasticachetypes.ReplicationGroup) {
			for _, fn := range addParentsFns["replicationgroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elasticachetypes.ReplicationGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.elasticache.cachesubnetgroup.sync", true) {
		list, err := s.fetcher.Get("cachesubnetgroup_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]elasticachetypes.CacheSubnetGroup); !ok {
			return gph, errors.New("cannot cast to '[]elasticachetypes.CacheSubnetGroup' type from fetch context")
		}
		for _, r := range list.([]elasticachetypes.CacheSubnetGroup) {
			for _, fn := range addParentsFns["cachesubnetgroup"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *elasticachetypes.CacheSubnetGroup) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Elasticache) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Elasticache) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.elasticache.sync", true)
}

type Eventbridge struct {
	fetcher           fetch.Fetcher
	region, profile   string
	config            map[string]any
	log               *logger.Logger
	EventbridgeClient *eventbridge.Client
}

func NewEventbridge(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	eventbridgeClient := eventbridge.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		eventbridgeClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Eventbridge{
		EventbridgeClient: eventbridgeClient,
		fetcher:           fetch.NewFetcher(awsfetch.BuildEventbridgeFetchFuncs(fetchConfig)),
		config:            extraConf,
		region:            region,
		profile:           profile,
		log:               log,
	}
}

func (s *Eventbridge) Name() string {
	return "eventbridge"
}

func (s *Eventbridge) Region() string {
	return s.region
}

func (s *Eventbridge) Profile() string {
	return s.profile
}

func (s *Eventbridge) ResourceTypes() []string {
	return []string{
		"eventbus",
		"eventrule",
	}
}

func (s *Eventbridge) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.eventbridge.eventbus.sync", true) {
		list, err := s.fetcher.Get("eventbus_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]eventbridgetypes.EventBus); !ok {
			return gph, errors.New("cannot cast to '[]eventbridgetypes.EventBus' type from fetch context")
		}
		for _, r := range list.([]eventbridgetypes.EventBus) {
			for _, fn := range addParentsFns["eventbus"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *eventbridgetypes.EventBus) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}
	if getBool(s.config, "aws.eventbridge.eventrule.sync", true) {
		list, err := s.fetcher.Get("eventrule_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]eventbridgetypes.Rule); !ok {
			return gph, errors.New("cannot cast to '[]eventbridgetypes.Rule' type from fetch context")
		}
		for _, r := range list.([]eventbridgetypes.Rule) {
			for _, fn := range addParentsFns["eventrule"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *eventbridgetypes.Rule) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Eventbridge) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Eventbridge) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.eventbridge.sync", true)
}

type Stepfunctions struct {
	fetcher         fetch.Fetcher
	region, profile string
	config          map[string]any
	log             *logger.Logger
	SfnClient       *sfn.Client
}

func NewStepfunctions(cfg aws.Config, profile string, extraConf map[string]any, log *logger.Logger) cloud.Service {
	region := cfg.Region
	sfnClient := sfn.NewFromConfig(cfg)

	fetchConfig := awsfetch.NewConfig(
		sfnClient,
	)
	fetchConfig.Extra = extraConf
	fetchConfig.Log = log

	return &Stepfunctions{
		SfnClient: sfnClient,
		fetcher:   fetch.NewFetcher(awsfetch.BuildStepfunctionsFetchFuncs(fetchConfig)),
		config:    extraConf,
		region:    region,
		profile:   profile,
		log:       log,
	}
}

func (s *Stepfunctions) Name() string {
	return "stepfunctions"
}

func (s *Stepfunctions) Region() string {
	return s.region
}

func (s *Stepfunctions) Profile() string {
	return s.profile
}

func (s *Stepfunctions) ResourceTypes() []string {
	return []string{
		"statemachine",
	}
}

func (s *Stepfunctions) Fetch(ctx context.Context) (cloud.GraphAPI, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	gph, err := s.fetcher.Fetch(context.WithValue(ctx, "region", s.region))
	defer s.fetcher.Reset()

	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case nil:
			continue
		default:
			var ae smithy.APIError
			if errors.As(ee, &ae) && ae.ErrorMessage() == accessDenied {
				allErrors.Add(cloud.ErrFetchAccessDenied)
			} else {
				allErrors.Add(ee)
			}
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	snap := gph.AsRDFGraphSnaphot()

	errc := make(chan error)
	var wg sync.WaitGroup
	if getBool(s.config, "aws.stepfunctions.statemachine.sync", true) {
		list, err := s.fetcher.Get("statemachine_objects")
		if err != nil {
			return gph, err
		}
		if _, ok := list.([]sfntypes.StateMachineListItem); !ok {
			return gph, errors.New("cannot cast to '[]sfntypes.StateMachineListItem' type from fetch context")
		}
		for _, r := range list.([]sfntypes.StateMachineListItem) {
			for _, fn := range addParentsFns["statemachine"] {
				wg.Add(1)
				go func(f addParentFn, snap tstore.RDFGraph, region string, res *sfntypes.StateMachineListItem) {
					defer wg.Done()
					err := f(gph, snap, region, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, snap, s.region, &r)
			}
		}
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *Stepfunctions) FetchByType(ctx context.Context, t string) (cloud.GraphAPI, error) {
	defer s.fetcher.Reset()
	return s.fetcher.FetchByType(context.WithValue(ctx, "region", s.region), t)
}

func (s *Stepfunctions) IsSyncDisabled() bool {
	return !getBool(s.config, "aws.stepfunctions.sync", true)
}
