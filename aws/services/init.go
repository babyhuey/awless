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
	"errors"

	awsspec "github.com/bootswithdefer/awless/aws/spec"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/graph"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/sync"
)

var (
	AccessService, InfraService, StorageService, MessagingService, DNSService, LambdaService, MonitoringService, CdnService, CloudformationService       cloud.Service
	EksService, DynamodbService, SecretsmanagerService, ApigatewayService, SsmService, EfsService, CloudtrailService, CloudwatchlogsService              cloud.Service
	ElasticacheService, EventbridgeService, StepfunctionsService, WafService, ConfigserviceService, KinesisService, RedshiftService, CodepipelineService cloud.Service
	CodebuildService, BeanstalkService                                                                                                                   cloud.Service
)

func Init(profile, region string, extraConf map[string]any, log *logger.Logger, profileSetterCallback func(val string) error, enableNetworkMonitor bool) error {
	if region == "" {
		return errors.New("empty AWS region. Set it with `awless config set aws.region`")
	}

	sb := newConfigResolver().withRegion(region).withProfile(profile).withNetworkMonitor(enableNetworkMonitor)
	sb = sb.withProfileSetter(profileSetterCallback).withLogger(log).withCredentialResolvers()

	cfg, err := sb.resolve()
	if err != nil {
		return err
	}

	AccessService = NewAccess(cfg, profile, extraConf, log)
	InfraService = NewInfra(cfg, profile, extraConf, log)
	StorageService = NewStorage(cfg, profile, extraConf, log)
	MessagingService = NewMessaging(cfg, profile, extraConf, log)
	DNSService = NewDNS(cfg, profile, extraConf, log)
	LambdaService = NewLambda(cfg, profile, extraConf, log)
	MonitoringService = NewMonitoring(cfg, profile, extraConf, log)
	CdnService = NewCDN(cfg, profile, extraConf, log)
	CloudformationService = NewCloudformation(cfg, profile, extraConf, log)
	EksService = NewEKS(cfg, profile, extraConf, log)
	DynamodbService = NewDynamodb(cfg, profile, extraConf, log)
	SecretsmanagerService = NewSecretsmanager(cfg, profile, extraConf, log)
	ApigatewayService = NewApigateway(cfg, profile, extraConf, log)
	SsmService = NewSSM(cfg, profile, extraConf, log)
	EfsService = NewEFS(cfg, profile, extraConf, log)
	CloudtrailService = NewCloudtrail(cfg, profile, extraConf, log)
	CloudwatchlogsService = NewCloudwatchlogs(cfg, profile, extraConf, log)
	ElasticacheService = NewElasticache(cfg, profile, extraConf, log)
	EventbridgeService = NewEventbridge(cfg, profile, extraConf, log)
	StepfunctionsService = NewStepfunctions(cfg, profile, extraConf, log)
	WafService = NewWaf(cfg, profile, extraConf, log)
	ConfigserviceService = NewConfigservice(cfg, profile, extraConf, log)
	KinesisService = NewKinesis(cfg, profile, extraConf, log)
	RedshiftService = NewRedshift(cfg, profile, extraConf, log)
	CodepipelineService = NewCodepipeline(cfg, profile, extraConf, log)
	CodebuildService = NewCodebuild(cfg, profile, extraConf, log)
	BeanstalkService = NewBeanstalk(cfg, profile, extraConf, log)

	cloud.ServiceRegistry[InfraService.Name()] = InfraService
	cloud.ServiceRegistry[AccessService.Name()] = AccessService
	cloud.ServiceRegistry[StorageService.Name()] = StorageService
	cloud.ServiceRegistry[MessagingService.Name()] = MessagingService
	cloud.ServiceRegistry[DNSService.Name()] = DNSService
	cloud.ServiceRegistry[LambdaService.Name()] = LambdaService
	cloud.ServiceRegistry[MonitoringService.Name()] = MonitoringService
	cloud.ServiceRegistry[CdnService.Name()] = CdnService
	cloud.ServiceRegistry[CloudformationService.Name()] = CloudformationService
	cloud.ServiceRegistry[EksService.Name()] = EksService
	cloud.ServiceRegistry[DynamodbService.Name()] = DynamodbService
	cloud.ServiceRegistry[SecretsmanagerService.Name()] = SecretsmanagerService
	cloud.ServiceRegistry[ApigatewayService.Name()] = ApigatewayService
	cloud.ServiceRegistry[SsmService.Name()] = SsmService
	cloud.ServiceRegistry[EfsService.Name()] = EfsService
	cloud.ServiceRegistry[CloudtrailService.Name()] = CloudtrailService
	cloud.ServiceRegistry[CloudwatchlogsService.Name()] = CloudwatchlogsService
	cloud.ServiceRegistry[ElasticacheService.Name()] = ElasticacheService
	cloud.ServiceRegistry[EventbridgeService.Name()] = EventbridgeService
	cloud.ServiceRegistry[StepfunctionsService.Name()] = StepfunctionsService
	cloud.ServiceRegistry[WafService.Name()] = WafService
	cloud.ServiceRegistry[ConfigserviceService.Name()] = ConfigserviceService
	cloud.ServiceRegistry[KinesisService.Name()] = KinesisService
	cloud.ServiceRegistry[RedshiftService.Name()] = RedshiftService
	cloud.ServiceRegistry[CodepipelineService.Name()] = CodepipelineService
	cloud.ServiceRegistry[CodebuildService.Name()] = CodebuildService
	cloud.ServiceRegistry[BeanstalkService.Name()] = BeanstalkService

	awsspec.CommandFactory = &awsspec.AWSFactory{
		Log: log,
		Cfg: cfg,
		Graph: &cloud.LazyGraph{LoadingFunc: func() cloud.GraphAPI {
			g, err := sync.LoadLocalGraphs(profile, region)
			if err != nil {
				g = graph.NewGraph()
			}
			return g
		}},
	}

	return nil
}

func getBool(m map[string]any, key string, def bool) bool {
	if b, ok := m[key].(bool); ok {
		return b
	}
	return def
}
