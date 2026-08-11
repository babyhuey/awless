/* Copyright 2017 WALLIX

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

// DO NOT EDIT
// This file was automatically generated with go generate
package awsat

import (
	"github.com/aws/aws-sdk-go-v2/aws"

	awsspec "github.com/bootswithdefer/awless/aws/spec"
	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
)

// AcceptanceFactory builds commands whose AWS clients are routed through a Mock.
//
// SDK v2 clients are concrete structs, so they cannot be replaced with an
// interface. Instead every command is constructed from Mock.Config(), whose
// APIOptions carry a middleware that intercepts the call before it is signed or
// sent. The generated constructors already do service.NewFromConfig(cfg), so no
// SetApi call is needed.
type AcceptanceFactory struct {
	Mock   *Mock
	Logger *logger.Logger
	Graph  cloud.GraphAPI
}

func NewAcceptanceFactory(mock *Mock, g cloud.GraphAPI, l ...*logger.Logger) *AcceptanceFactory {
	lg := logger.DiscardLogger
	if len(l) > 0 {
		lg = l[0]
	}
	return &AcceptanceFactory{Mock: mock, Graph: g, Logger: lg}
}

func (f *AcceptanceFactory) config() aws.Config {
	if f.Mock == nil {
		return aws.Config{}
	}
	return f.Mock.Config()
}

func (f *AcceptanceFactory) Build(key string) func() any {
	switch key {
	case "attachalarm":
		return func() any {
			return awsspec.NewAttachAlarm(f.config(), f.Graph, f.Logger)
		}
	case "attachclassicloadbalancer":
		return func() any {
			return awsspec.NewAttachClassicLoadbalancer(f.config(), f.Graph, f.Logger)
		}
	case "attachcontainertask":
		return func() any {
			return awsspec.NewAttachContainertask(f.config(), f.Graph, f.Logger)
		}
	case "attachelasticip":
		return func() any {
			return awsspec.NewAttachElasticip(f.config(), f.Graph, f.Logger)
		}
	case "attacheventtarget":
		return func() any {
			return awsspec.NewAttachEventtarget(f.config(), f.Graph, f.Logger)
		}
	case "attachinstance":
		return func() any {
			return awsspec.NewAttachInstance(f.config(), f.Graph, f.Logger)
		}
	case "attachinstanceprofile":
		return func() any {
			return awsspec.NewAttachInstanceprofile(f.config(), f.Graph, f.Logger)
		}
	case "attachinternetgateway":
		return func() any {
			return awsspec.NewAttachInternetgateway(f.config(), f.Graph, f.Logger)
		}
	case "attachlistener":
		return func() any {
			return awsspec.NewAttachListener(f.config(), f.Graph, f.Logger)
		}
	case "attachmfadevice":
		return func() any {
			return awsspec.NewAttachMfadevice(f.config(), f.Graph, f.Logger)
		}
	case "attachnetworkinterface":
		return func() any {
			return awsspec.NewAttachNetworkinterface(f.config(), f.Graph, f.Logger)
		}
	case "attachpolicy":
		return func() any {
			return awsspec.NewAttachPolicy(f.config(), f.Graph, f.Logger)
		}
	case "attachrole":
		return func() any {
			return awsspec.NewAttachRole(f.config(), f.Graph, f.Logger)
		}
	case "attachroutetable":
		return func() any {
			return awsspec.NewAttachRoutetable(f.config(), f.Graph, f.Logger)
		}
	case "attachsecuritygroup":
		return func() any {
			return awsspec.NewAttachSecuritygroup(f.config(), f.Graph, f.Logger)
		}
	case "attachuser":
		return func() any {
			return awsspec.NewAttachUser(f.config(), f.Graph, f.Logger)
		}
	case "attachvolume":
		return func() any {
			return awsspec.NewAttachVolume(f.config(), f.Graph, f.Logger)
		}
	case "authenticateregistry":
		return func() any {
			return awsspec.NewAuthenticateRegistry(f.config(), f.Graph, f.Logger)
		}
	case "checkcertificate":
		return func() any {
			return awsspec.NewCheckCertificate(f.config(), f.Graph, f.Logger)
		}
	case "checkdatabase":
		return func() any {
			return awsspec.NewCheckDatabase(f.config(), f.Graph, f.Logger)
		}
	case "checkdistribution":
		return func() any {
			return awsspec.NewCheckDistribution(f.config(), f.Graph, f.Logger)
		}
	case "checkinstance":
		return func() any {
			return awsspec.NewCheckInstance(f.config(), f.Graph, f.Logger)
		}
	case "checkloadbalancer":
		return func() any {
			return awsspec.NewCheckLoadbalancer(f.config(), f.Graph, f.Logger)
		}
	case "checknatgateway":
		return func() any {
			return awsspec.NewCheckNatgateway(f.config(), f.Graph, f.Logger)
		}
	case "checknetworkinterface":
		return func() any {
			return awsspec.NewCheckNetworkinterface(f.config(), f.Graph, f.Logger)
		}
	case "checkscalinggroup":
		return func() any {
			return awsspec.NewCheckScalinggroup(f.config(), f.Graph, f.Logger)
		}
	case "checksecuritygroup":
		return func() any {
			return awsspec.NewCheckSecuritygroup(f.config(), f.Graph, f.Logger)
		}
	case "checkvolume":
		return func() any {
			return awsspec.NewCheckVolume(f.config(), f.Graph, f.Logger)
		}
	case "copyimage":
		return func() any {
			return awsspec.NewCopyImage(f.config(), f.Graph, f.Logger)
		}
	case "copysnapshot":
		return func() any {
			return awsspec.NewCopySnapshot(f.config(), f.Graph, f.Logger)
		}
	case "createaccesskey":
		return func() any {
			return awsspec.NewCreateAccesskey(f.config(), f.Graph, f.Logger)
		}
	case "createalarm":
		return func() any {
			return awsspec.NewCreateAlarm(f.config(), f.Graph, f.Logger)
		}
	case "createapigateway":
		return func() any {
			return awsspec.NewCreateApigateway(f.config(), f.Graph, f.Logger)
		}
	case "createapigatewayroute":
		return func() any {
			return awsspec.NewCreateApigatewayroute(f.config(), f.Graph, f.Logger)
		}
	case "createapigatewaystage":
		return func() any {
			return awsspec.NewCreateApigatewaystage(f.config(), f.Graph, f.Logger)
		}
	case "createapplication":
		return func() any {
			return awsspec.NewCreateApplication(f.config(), f.Graph, f.Logger)
		}
	case "createappscalingpolicy":
		return func() any {
			return awsspec.NewCreateAppscalingpolicy(f.config(), f.Graph, f.Logger)
		}
	case "createappscalingtarget":
		return func() any {
			return awsspec.NewCreateAppscalingtarget(f.config(), f.Graph, f.Logger)
		}
	case "createbroker":
		return func() any {
			return awsspec.NewCreateBroker(f.config(), f.Graph, f.Logger)
		}
	case "createbucket":
		return func() any {
			return awsspec.NewCreateBucket(f.config(), f.Graph, f.Logger)
		}
	case "createbuildproject":
		return func() any {
			return awsspec.NewCreateBuildproject(f.config(), f.Graph, f.Logger)
		}
	case "createcachecluster":
		return func() any {
			return awsspec.NewCreateCachecluster(f.config(), f.Graph, f.Logger)
		}
	case "createcachesubnetgroup":
		return func() any {
			return awsspec.NewCreateCachesubnetgroup(f.config(), f.Graph, f.Logger)
		}
	case "createcertificate":
		return func() any {
			return awsspec.NewCreateCertificate(f.config(), f.Graph, f.Logger)
		}
	case "createclassicloadbalancer":
		return func() any {
			return awsspec.NewCreateClassicLoadbalancer(f.config(), f.Graph, f.Logger)
		}
	case "createconfigrule":
		return func() any {
			return awsspec.NewCreateConfigrule(f.config(), f.Graph, f.Logger)
		}
	case "createconfigurationset":
		return func() any {
			return awsspec.NewCreateConfigurationset(f.config(), f.Graph, f.Logger)
		}
	case "createcontainercluster":
		return func() any {
			return awsspec.NewCreateContainercluster(f.config(), f.Graph, f.Logger)
		}
	case "createcrawler":
		return func() any {
			return awsspec.NewCreateCrawler(f.config(), f.Graph, f.Logger)
		}
	case "createdatabase":
		return func() any {
			return awsspec.NewCreateDatabase(f.config(), f.Graph, f.Logger)
		}
	case "createdbsubnetgroup":
		return func() any {
			return awsspec.NewCreateDbsubnetgroup(f.config(), f.Graph, f.Logger)
		}
	case "createdeployapplication":
		return func() any {
			return awsspec.NewCreateDeployapplication(f.config(), f.Graph, f.Logger)
		}
	case "createdeployment":
		return func() any {
			return awsspec.NewCreateDeployment(f.config(), f.Graph, f.Logger)
		}
	case "createdeploymentgroup":
		return func() any {
			return awsspec.NewCreateDeploymentgroup(f.config(), f.Graph, f.Logger)
		}
	case "createdistribution":
		return func() any {
			return awsspec.NewCreateDistribution(f.config(), f.Graph, f.Logger)
		}
	case "createdynamodbtable":
		return func() any {
			return awsspec.NewCreateDynamodbtable(f.config(), f.Graph, f.Logger)
		}
	case "createekscluster":
		return func() any {
			return awsspec.NewCreateEkscluster(f.config(), f.Graph, f.Logger)
		}
	case "createeksnodegroup":
		return func() any {
			return awsspec.NewCreateEksnodegroup(f.config(), f.Graph, f.Logger)
		}
	case "createelasticip":
		return func() any {
			return awsspec.NewCreateElasticip(f.config(), f.Graph, f.Logger)
		}
	case "createemailidentity":
		return func() any {
			return awsspec.NewCreateEmailidentity(f.config(), f.Graph, f.Logger)
		}
	case "createenvironment":
		return func() any {
			return awsspec.NewCreateEnvironment(f.config(), f.Graph, f.Logger)
		}
	case "createeventbus":
		return func() any {
			return awsspec.NewCreateEventbus(f.config(), f.Graph, f.Logger)
		}
	case "createeventrule":
		return func() any {
			return awsspec.NewCreateEventrule(f.config(), f.Graph, f.Logger)
		}
	case "createfilesystem":
		return func() any {
			return awsspec.NewCreateFilesystem(f.config(), f.Graph, f.Logger)
		}
	case "createfunction":
		return func() any {
			return awsspec.NewCreateFunction(f.config(), f.Graph, f.Logger)
		}
	case "creategluedatabase":
		return func() any {
			return awsspec.NewCreateGluedatabase(f.config(), f.Graph, f.Logger)
		}
	case "creategroup":
		return func() any {
			return awsspec.NewCreateGroup(f.config(), f.Graph, f.Logger)
		}
	case "createidentitypool":
		return func() any {
			return awsspec.NewCreateIdentitypool(f.config(), f.Graph, f.Logger)
		}
	case "createimage":
		return func() any {
			return awsspec.NewCreateImage(f.config(), f.Graph, f.Logger)
		}
	case "createinstance":
		return func() any {
			return awsspec.NewCreateInstance(f.config(), f.Graph, f.Logger)
		}
	case "createinstanceprofile":
		return func() any {
			return awsspec.NewCreateInstanceprofile(f.config(), f.Graph, f.Logger)
		}
	case "createinternetgateway":
		return func() any {
			return awsspec.NewCreateInternetgateway(f.config(), f.Graph, f.Logger)
		}
	case "createipset":
		return func() any {
			return awsspec.NewCreateIpset(f.config(), f.Graph, f.Logger)
		}
	case "createjob":
		return func() any {
			return awsspec.NewCreateJob(f.config(), f.Graph, f.Logger)
		}
	case "createkafkacluster":
		return func() any {
			return awsspec.NewCreateKafkacluster(f.config(), f.Graph, f.Logger)
		}
	case "createkeypair":
		return func() any {
			return awsspec.NewCreateKeypair(f.config(), f.Graph, f.Logger)
		}
	case "createlaunchconfiguration":
		return func() any {
			return awsspec.NewCreateLaunchconfiguration(f.config(), f.Graph, f.Logger)
		}
	case "createlistener":
		return func() any {
			return awsspec.NewCreateListener(f.config(), f.Graph, f.Logger)
		}
	case "createloadbalancer":
		return func() any {
			return awsspec.NewCreateLoadbalancer(f.config(), f.Graph, f.Logger)
		}
	case "createloggroup":
		return func() any {
			return awsspec.NewCreateLoggroup(f.config(), f.Graph, f.Logger)
		}
	case "createloginprofile":
		return func() any {
			return awsspec.NewCreateLoginprofile(f.config(), f.Graph, f.Logger)
		}
	case "createmfadevice":
		return func() any {
			return awsspec.NewCreateMfadevice(f.config(), f.Graph, f.Logger)
		}
	case "createnatgateway":
		return func() any {
			return awsspec.NewCreateNatgateway(f.config(), f.Graph, f.Logger)
		}
	case "createnetworkinterface":
		return func() any {
			return awsspec.NewCreateNetworkinterface(f.config(), f.Graph, f.Logger)
		}
	case "createpipeline":
		return func() any {
			return awsspec.NewCreatePipeline(f.config(), f.Graph, f.Logger)
		}
	case "createpolicy":
		return func() any {
			return awsspec.NewCreatePolicy(f.config(), f.Graph, f.Logger)
		}
	case "createqueue":
		return func() any {
			return awsspec.NewCreateQueue(f.config(), f.Graph, f.Logger)
		}
	case "createrecord":
		return func() any {
			return awsspec.NewCreateRecord(f.config(), f.Graph, f.Logger)
		}
	case "createredshiftcluster":
		return func() any {
			return awsspec.NewCreateRedshiftcluster(f.config(), f.Graph, f.Logger)
		}
	case "createredshiftsubnetgroup":
		return func() any {
			return awsspec.NewCreateRedshiftsubnetgroup(f.config(), f.Graph, f.Logger)
		}
	case "createreplicationgroup":
		return func() any {
			return awsspec.NewCreateReplicationgroup(f.config(), f.Graph, f.Logger)
		}
	case "createrepository":
		return func() any {
			return awsspec.NewCreateRepository(f.config(), f.Graph, f.Logger)
		}
	case "createrole":
		return func() any {
			return awsspec.NewCreateRole(f.config(), f.Graph, f.Logger)
		}
	case "createroute":
		return func() any {
			return awsspec.NewCreateRoute(f.config(), f.Graph, f.Logger)
		}
	case "createroutetable":
		return func() any {
			return awsspec.NewCreateRoutetable(f.config(), f.Graph, f.Logger)
		}
	case "createrulegroup":
		return func() any {
			return awsspec.NewCreateRulegroup(f.config(), f.Graph, f.Logger)
		}
	case "creates3object":
		return func() any {
			return awsspec.NewCreateS3object(f.config(), f.Graph, f.Logger)
		}
	case "createscalinggroup":
		return func() any {
			return awsspec.NewCreateScalinggroup(f.config(), f.Graph, f.Logger)
		}
	case "createscalingpolicy":
		return func() any {
			return awsspec.NewCreateScalingpolicy(f.config(), f.Graph, f.Logger)
		}
	case "createsecret":
		return func() any {
			return awsspec.NewCreateSecret(f.config(), f.Graph, f.Logger)
		}
	case "createsecuritygroup":
		return func() any {
			return awsspec.NewCreateSecuritygroup(f.config(), f.Graph, f.Logger)
		}
	case "createsnapshot":
		return func() any {
			return awsspec.NewCreateSnapshot(f.config(), f.Graph, f.Logger)
		}
	case "createssmparameter":
		return func() any {
			return awsspec.NewCreateSsmparameter(f.config(), f.Graph, f.Logger)
		}
	case "createstack":
		return func() any {
			return awsspec.NewCreateStack(f.config(), f.Graph, f.Logger)
		}
	case "createstatemachine":
		return func() any {
			return awsspec.NewCreateStatemachine(f.config(), f.Graph, f.Logger)
		}
	case "createstream":
		return func() any {
			return awsspec.NewCreateStream(f.config(), f.Graph, f.Logger)
		}
	case "createsubnet":
		return func() any {
			return awsspec.NewCreateSubnet(f.config(), f.Graph, f.Logger)
		}
	case "createsubscription":
		return func() any {
			return awsspec.NewCreateSubscription(f.config(), f.Graph, f.Logger)
		}
	case "createtag":
		return func() any {
			return awsspec.NewCreateTag(f.config(), f.Graph, f.Logger)
		}
	case "createtargetgroup":
		return func() any {
			return awsspec.NewCreateTargetgroup(f.config(), f.Graph, f.Logger)
		}
	case "createtopic":
		return func() any {
			return awsspec.NewCreateTopic(f.config(), f.Graph, f.Logger)
		}
	case "createtrail":
		return func() any {
			return awsspec.NewCreateTrail(f.config(), f.Graph, f.Logger)
		}
	case "createtransitgateway":
		return func() any {
			return awsspec.NewCreateTransitgateway(f.config(), f.Graph, f.Logger)
		}
	case "createtransitgatewayattachment":
		return func() any {
			return awsspec.NewCreateTransitgatewayattachment(f.config(), f.Graph, f.Logger)
		}
	case "createtransitgatewayroutetable":
		return func() any {
			return awsspec.NewCreateTransitgatewayroutetable(f.config(), f.Graph, f.Logger)
		}
	case "createuser":
		return func() any {
			return awsspec.NewCreateUser(f.config(), f.Graph, f.Logger)
		}
	case "createuserpool":
		return func() any {
			return awsspec.NewCreateUserpool(f.config(), f.Graph, f.Logger)
		}
	case "createvolume":
		return func() any {
			return awsspec.NewCreateVolume(f.config(), f.Graph, f.Logger)
		}
	case "createvpc":
		return func() any {
			return awsspec.NewCreateVpc(f.config(), f.Graph, f.Logger)
		}
	case "createvpcendpoint":
		return func() any {
			return awsspec.NewCreateVpcendpoint(f.config(), f.Graph, f.Logger)
		}
	case "createwebacl":
		return func() any {
			return awsspec.NewCreateWebacl(f.config(), f.Graph, f.Logger)
		}
	case "createzone":
		return func() any {
			return awsspec.NewCreateZone(f.config(), f.Graph, f.Logger)
		}
	case "deleteaccesskey":
		return func() any {
			return awsspec.NewDeleteAccesskey(f.config(), f.Graph, f.Logger)
		}
	case "deletealarm":
		return func() any {
			return awsspec.NewDeleteAlarm(f.config(), f.Graph, f.Logger)
		}
	case "deleteapigateway":
		return func() any {
			return awsspec.NewDeleteApigateway(f.config(), f.Graph, f.Logger)
		}
	case "deleteapigatewayroute":
		return func() any {
			return awsspec.NewDeleteApigatewayroute(f.config(), f.Graph, f.Logger)
		}
	case "deleteapigatewaystage":
		return func() any {
			return awsspec.NewDeleteApigatewaystage(f.config(), f.Graph, f.Logger)
		}
	case "deleteapplication":
		return func() any {
			return awsspec.NewDeleteApplication(f.config(), f.Graph, f.Logger)
		}
	case "deleteappscalingpolicy":
		return func() any {
			return awsspec.NewDeleteAppscalingpolicy(f.config(), f.Graph, f.Logger)
		}
	case "deleteappscalingtarget":
		return func() any {
			return awsspec.NewDeleteAppscalingtarget(f.config(), f.Graph, f.Logger)
		}
	case "deletebroker":
		return func() any {
			return awsspec.NewDeleteBroker(f.config(), f.Graph, f.Logger)
		}
	case "deletebucket":
		return func() any {
			return awsspec.NewDeleteBucket(f.config(), f.Graph, f.Logger)
		}
	case "deletebuildproject":
		return func() any {
			return awsspec.NewDeleteBuildproject(f.config(), f.Graph, f.Logger)
		}
	case "deletecachecluster":
		return func() any {
			return awsspec.NewDeleteCachecluster(f.config(), f.Graph, f.Logger)
		}
	case "deletecachesubnetgroup":
		return func() any {
			return awsspec.NewDeleteCachesubnetgroup(f.config(), f.Graph, f.Logger)
		}
	case "deletecertificate":
		return func() any {
			return awsspec.NewDeleteCertificate(f.config(), f.Graph, f.Logger)
		}
	case "deleteclassicloadbalancer":
		return func() any {
			return awsspec.NewDeleteClassicLoadbalancer(f.config(), f.Graph, f.Logger)
		}
	case "deleteconfigrule":
		return func() any {
			return awsspec.NewDeleteConfigrule(f.config(), f.Graph, f.Logger)
		}
	case "deleteconfigurationset":
		return func() any {
			return awsspec.NewDeleteConfigurationset(f.config(), f.Graph, f.Logger)
		}
	case "deletecontainercluster":
		return func() any {
			return awsspec.NewDeleteContainercluster(f.config(), f.Graph, f.Logger)
		}
	case "deletecontainertask":
		return func() any {
			return awsspec.NewDeleteContainertask(f.config(), f.Graph, f.Logger)
		}
	case "deletecrawler":
		return func() any {
			return awsspec.NewDeleteCrawler(f.config(), f.Graph, f.Logger)
		}
	case "deletedatabase":
		return func() any {
			return awsspec.NewDeleteDatabase(f.config(), f.Graph, f.Logger)
		}
	case "deletedbsubnetgroup":
		return func() any {
			return awsspec.NewDeleteDbsubnetgroup(f.config(), f.Graph, f.Logger)
		}
	case "deletedeployapplication":
		return func() any {
			return awsspec.NewDeleteDeployapplication(f.config(), f.Graph, f.Logger)
		}
	case "deletedeploymentgroup":
		return func() any {
			return awsspec.NewDeleteDeploymentgroup(f.config(), f.Graph, f.Logger)
		}
	case "deletedistribution":
		return func() any {
			return awsspec.NewDeleteDistribution(f.config(), f.Graph, f.Logger)
		}
	case "deletedynamodbtable":
		return func() any {
			return awsspec.NewDeleteDynamodbtable(f.config(), f.Graph, f.Logger)
		}
	case "deleteekscluster":
		return func() any {
			return awsspec.NewDeleteEkscluster(f.config(), f.Graph, f.Logger)
		}
	case "deleteeksnodegroup":
		return func() any {
			return awsspec.NewDeleteEksnodegroup(f.config(), f.Graph, f.Logger)
		}
	case "deleteelasticip":
		return func() any {
			return awsspec.NewDeleteElasticip(f.config(), f.Graph, f.Logger)
		}
	case "deleteemailidentity":
		return func() any {
			return awsspec.NewDeleteEmailidentity(f.config(), f.Graph, f.Logger)
		}
	case "deleteenvironment":
		return func() any {
			return awsspec.NewDeleteEnvironment(f.config(), f.Graph, f.Logger)
		}
	case "deleteeventbus":
		return func() any {
			return awsspec.NewDeleteEventbus(f.config(), f.Graph, f.Logger)
		}
	case "deleteeventrule":
		return func() any {
			return awsspec.NewDeleteEventrule(f.config(), f.Graph, f.Logger)
		}
	case "deletefilesystem":
		return func() any {
			return awsspec.NewDeleteFilesystem(f.config(), f.Graph, f.Logger)
		}
	case "deletefunction":
		return func() any {
			return awsspec.NewDeleteFunction(f.config(), f.Graph, f.Logger)
		}
	case "deletegluedatabase":
		return func() any {
			return awsspec.NewDeleteGluedatabase(f.config(), f.Graph, f.Logger)
		}
	case "deletegroup":
		return func() any {
			return awsspec.NewDeleteGroup(f.config(), f.Graph, f.Logger)
		}
	case "deleteidentitypool":
		return func() any {
			return awsspec.NewDeleteIdentitypool(f.config(), f.Graph, f.Logger)
		}
	case "deleteimage":
		return func() any {
			return awsspec.NewDeleteImage(f.config(), f.Graph, f.Logger)
		}
	case "deleteinstance":
		return func() any {
			return awsspec.NewDeleteInstance(f.config(), f.Graph, f.Logger)
		}
	case "deleteinstanceprofile":
		return func() any {
			return awsspec.NewDeleteInstanceprofile(f.config(), f.Graph, f.Logger)
		}
	case "deleteinternetgateway":
		return func() any {
			return awsspec.NewDeleteInternetgateway(f.config(), f.Graph, f.Logger)
		}
	case "deleteipset":
		return func() any {
			return awsspec.NewDeleteIpset(f.config(), f.Graph, f.Logger)
		}
	case "deletejob":
		return func() any {
			return awsspec.NewDeleteJob(f.config(), f.Graph, f.Logger)
		}
	case "deletekafkacluster":
		return func() any {
			return awsspec.NewDeleteKafkacluster(f.config(), f.Graph, f.Logger)
		}
	case "deletekeypair":
		return func() any {
			return awsspec.NewDeleteKeypair(f.config(), f.Graph, f.Logger)
		}
	case "deletelaunchconfiguration":
		return func() any {
			return awsspec.NewDeleteLaunchconfiguration(f.config(), f.Graph, f.Logger)
		}
	case "deletelistener":
		return func() any {
			return awsspec.NewDeleteListener(f.config(), f.Graph, f.Logger)
		}
	case "deleteloadbalancer":
		return func() any {
			return awsspec.NewDeleteLoadbalancer(f.config(), f.Graph, f.Logger)
		}
	case "deleteloggroup":
		return func() any {
			return awsspec.NewDeleteLoggroup(f.config(), f.Graph, f.Logger)
		}
	case "deleteloginprofile":
		return func() any {
			return awsspec.NewDeleteLoginprofile(f.config(), f.Graph, f.Logger)
		}
	case "deletemfadevice":
		return func() any {
			return awsspec.NewDeleteMfadevice(f.config(), f.Graph, f.Logger)
		}
	case "deletenatgateway":
		return func() any {
			return awsspec.NewDeleteNatgateway(f.config(), f.Graph, f.Logger)
		}
	case "deletenetworkinterface":
		return func() any {
			return awsspec.NewDeleteNetworkinterface(f.config(), f.Graph, f.Logger)
		}
	case "deletepipeline":
		return func() any {
			return awsspec.NewDeletePipeline(f.config(), f.Graph, f.Logger)
		}
	case "deletepolicy":
		return func() any {
			return awsspec.NewDeletePolicy(f.config(), f.Graph, f.Logger)
		}
	case "deletequeue":
		return func() any {
			return awsspec.NewDeleteQueue(f.config(), f.Graph, f.Logger)
		}
	case "deleterecord":
		return func() any {
			return awsspec.NewDeleteRecord(f.config(), f.Graph, f.Logger)
		}
	case "deleteredshiftcluster":
		return func() any {
			return awsspec.NewDeleteRedshiftcluster(f.config(), f.Graph, f.Logger)
		}
	case "deleteredshiftsubnetgroup":
		return func() any {
			return awsspec.NewDeleteRedshiftsubnetgroup(f.config(), f.Graph, f.Logger)
		}
	case "deletereplicationgroup":
		return func() any {
			return awsspec.NewDeleteReplicationgroup(f.config(), f.Graph, f.Logger)
		}
	case "deleterepository":
		return func() any {
			return awsspec.NewDeleteRepository(f.config(), f.Graph, f.Logger)
		}
	case "deleterole":
		return func() any {
			return awsspec.NewDeleteRole(f.config(), f.Graph, f.Logger)
		}
	case "deleteroute":
		return func() any {
			return awsspec.NewDeleteRoute(f.config(), f.Graph, f.Logger)
		}
	case "deleteroutetable":
		return func() any {
			return awsspec.NewDeleteRoutetable(f.config(), f.Graph, f.Logger)
		}
	case "deleterulegroup":
		return func() any {
			return awsspec.NewDeleteRulegroup(f.config(), f.Graph, f.Logger)
		}
	case "deletes3object":
		return func() any {
			return awsspec.NewDeleteS3object(f.config(), f.Graph, f.Logger)
		}
	case "deletescalinggroup":
		return func() any {
			return awsspec.NewDeleteScalinggroup(f.config(), f.Graph, f.Logger)
		}
	case "deletescalingpolicy":
		return func() any {
			return awsspec.NewDeleteScalingpolicy(f.config(), f.Graph, f.Logger)
		}
	case "deletesecret":
		return func() any {
			return awsspec.NewDeleteSecret(f.config(), f.Graph, f.Logger)
		}
	case "deletesecuritygroup":
		return func() any {
			return awsspec.NewDeleteSecuritygroup(f.config(), f.Graph, f.Logger)
		}
	case "deletesnapshot":
		return func() any {
			return awsspec.NewDeleteSnapshot(f.config(), f.Graph, f.Logger)
		}
	case "deletessmparameter":
		return func() any {
			return awsspec.NewDeleteSsmparameter(f.config(), f.Graph, f.Logger)
		}
	case "deletestack":
		return func() any {
			return awsspec.NewDeleteStack(f.config(), f.Graph, f.Logger)
		}
	case "deletestatemachine":
		return func() any {
			return awsspec.NewDeleteStatemachine(f.config(), f.Graph, f.Logger)
		}
	case "deletestream":
		return func() any {
			return awsspec.NewDeleteStream(f.config(), f.Graph, f.Logger)
		}
	case "deletesubnet":
		return func() any {
			return awsspec.NewDeleteSubnet(f.config(), f.Graph, f.Logger)
		}
	case "deletesubscription":
		return func() any {
			return awsspec.NewDeleteSubscription(f.config(), f.Graph, f.Logger)
		}
	case "deletetag":
		return func() any {
			return awsspec.NewDeleteTag(f.config(), f.Graph, f.Logger)
		}
	case "deletetargetgroup":
		return func() any {
			return awsspec.NewDeleteTargetgroup(f.config(), f.Graph, f.Logger)
		}
	case "deletetopic":
		return func() any {
			return awsspec.NewDeleteTopic(f.config(), f.Graph, f.Logger)
		}
	case "deletetrail":
		return func() any {
			return awsspec.NewDeleteTrail(f.config(), f.Graph, f.Logger)
		}
	case "deletetransitgateway":
		return func() any {
			return awsspec.NewDeleteTransitgateway(f.config(), f.Graph, f.Logger)
		}
	case "deletetransitgatewayattachment":
		return func() any {
			return awsspec.NewDeleteTransitgatewayattachment(f.config(), f.Graph, f.Logger)
		}
	case "deletetransitgatewayroutetable":
		return func() any {
			return awsspec.NewDeleteTransitgatewayroutetable(f.config(), f.Graph, f.Logger)
		}
	case "deleteuser":
		return func() any {
			return awsspec.NewDeleteUser(f.config(), f.Graph, f.Logger)
		}
	case "deleteuserpool":
		return func() any {
			return awsspec.NewDeleteUserpool(f.config(), f.Graph, f.Logger)
		}
	case "deletevolume":
		return func() any {
			return awsspec.NewDeleteVolume(f.config(), f.Graph, f.Logger)
		}
	case "deletevpc":
		return func() any {
			return awsspec.NewDeleteVpc(f.config(), f.Graph, f.Logger)
		}
	case "deletevpcendpoint":
		return func() any {
			return awsspec.NewDeleteVpcendpoint(f.config(), f.Graph, f.Logger)
		}
	case "deletewebacl":
		return func() any {
			return awsspec.NewDeleteWebacl(f.config(), f.Graph, f.Logger)
		}
	case "deletezone":
		return func() any {
			return awsspec.NewDeleteZone(f.config(), f.Graph, f.Logger)
		}
	case "detachalarm":
		return func() any {
			return awsspec.NewDetachAlarm(f.config(), f.Graph, f.Logger)
		}
	case "detachclassicloadbalancer":
		return func() any {
			return awsspec.NewDetachClassicLoadbalancer(f.config(), f.Graph, f.Logger)
		}
	case "detachcontainertask":
		return func() any {
			return awsspec.NewDetachContainertask(f.config(), f.Graph, f.Logger)
		}
	case "detachelasticip":
		return func() any {
			return awsspec.NewDetachElasticip(f.config(), f.Graph, f.Logger)
		}
	case "detacheventtarget":
		return func() any {
			return awsspec.NewDetachEventtarget(f.config(), f.Graph, f.Logger)
		}
	case "detachinstance":
		return func() any {
			return awsspec.NewDetachInstance(f.config(), f.Graph, f.Logger)
		}
	case "detachinstanceprofile":
		return func() any {
			return awsspec.NewDetachInstanceprofile(f.config(), f.Graph, f.Logger)
		}
	case "detachinternetgateway":
		return func() any {
			return awsspec.NewDetachInternetgateway(f.config(), f.Graph, f.Logger)
		}
	case "detachmfadevice":
		return func() any {
			return awsspec.NewDetachMfadevice(f.config(), f.Graph, f.Logger)
		}
	case "detachnetworkinterface":
		return func() any {
			return awsspec.NewDetachNetworkinterface(f.config(), f.Graph, f.Logger)
		}
	case "detachpolicy":
		return func() any {
			return awsspec.NewDetachPolicy(f.config(), f.Graph, f.Logger)
		}
	case "detachrole":
		return func() any {
			return awsspec.NewDetachRole(f.config(), f.Graph, f.Logger)
		}
	case "detachroutetable":
		return func() any {
			return awsspec.NewDetachRoutetable(f.config(), f.Graph, f.Logger)
		}
	case "detachsecuritygroup":
		return func() any {
			return awsspec.NewDetachSecuritygroup(f.config(), f.Graph, f.Logger)
		}
	case "detachuser":
		return func() any {
			return awsspec.NewDetachUser(f.config(), f.Graph, f.Logger)
		}
	case "detachvolume":
		return func() any {
			return awsspec.NewDetachVolume(f.config(), f.Graph, f.Logger)
		}
	case "importimage":
		return func() any {
			return awsspec.NewImportImage(f.config(), f.Graph, f.Logger)
		}
	case "restartdatabase":
		return func() any {
			return awsspec.NewRestartDatabase(f.config(), f.Graph, f.Logger)
		}
	case "restartinstance":
		return func() any {
			return awsspec.NewRestartInstance(f.config(), f.Graph, f.Logger)
		}
	case "startalarm":
		return func() any {
			return awsspec.NewStartAlarm(f.config(), f.Graph, f.Logger)
		}
	case "startbuildproject":
		return func() any {
			return awsspec.NewStartBuildproject(f.config(), f.Graph, f.Logger)
		}
	case "startcontainertask":
		return func() any {
			return awsspec.NewStartContainertask(f.config(), f.Graph, f.Logger)
		}
	case "startcrawler":
		return func() any {
			return awsspec.NewStartCrawler(f.config(), f.Graph, f.Logger)
		}
	case "startdatabase":
		return func() any {
			return awsspec.NewStartDatabase(f.config(), f.Graph, f.Logger)
		}
	case "starteventrule":
		return func() any {
			return awsspec.NewStartEventrule(f.config(), f.Graph, f.Logger)
		}
	case "startexecution":
		return func() any {
			return awsspec.NewStartExecution(f.config(), f.Graph, f.Logger)
		}
	case "startinstance":
		return func() any {
			return awsspec.NewStartInstance(f.config(), f.Graph, f.Logger)
		}
	case "startjob":
		return func() any {
			return awsspec.NewStartJob(f.config(), f.Graph, f.Logger)
		}
	case "startpipeline":
		return func() any {
			return awsspec.NewStartPipeline(f.config(), f.Graph, f.Logger)
		}
	case "starttrail":
		return func() any {
			return awsspec.NewStartTrail(f.config(), f.Graph, f.Logger)
		}
	case "stopalarm":
		return func() any {
			return awsspec.NewStopAlarm(f.config(), f.Graph, f.Logger)
		}
	case "stopbuildproject":
		return func() any {
			return awsspec.NewStopBuildproject(f.config(), f.Graph, f.Logger)
		}
	case "stopcontainertask":
		return func() any {
			return awsspec.NewStopContainertask(f.config(), f.Graph, f.Logger)
		}
	case "stopcrawler":
		return func() any {
			return awsspec.NewStopCrawler(f.config(), f.Graph, f.Logger)
		}
	case "stopdatabase":
		return func() any {
			return awsspec.NewStopDatabase(f.config(), f.Graph, f.Logger)
		}
	case "stopdeployment":
		return func() any {
			return awsspec.NewStopDeployment(f.config(), f.Graph, f.Logger)
		}
	case "stopeventrule":
		return func() any {
			return awsspec.NewStopEventrule(f.config(), f.Graph, f.Logger)
		}
	case "stopexecution":
		return func() any {
			return awsspec.NewStopExecution(f.config(), f.Graph, f.Logger)
		}
	case "stopinstance":
		return func() any {
			return awsspec.NewStopInstance(f.config(), f.Graph, f.Logger)
		}
	case "stoppipeline":
		return func() any {
			return awsspec.NewStopPipeline(f.config(), f.Graph, f.Logger)
		}
	case "stoptrail":
		return func() any {
			return awsspec.NewStopTrail(f.config(), f.Graph, f.Logger)
		}
	case "updateapplication":
		return func() any {
			return awsspec.NewUpdateApplication(f.config(), f.Graph, f.Logger)
		}
	case "updatebucket":
		return func() any {
			return awsspec.NewUpdateBucket(f.config(), f.Graph, f.Logger)
		}
	case "updatebuildproject":
		return func() any {
			return awsspec.NewUpdateBuildproject(f.config(), f.Graph, f.Logger)
		}
	case "updatecachecluster":
		return func() any {
			return awsspec.NewUpdateCachecluster(f.config(), f.Graph, f.Logger)
		}
	case "updatecachesubnetgroup":
		return func() any {
			return awsspec.NewUpdateCachesubnetgroup(f.config(), f.Graph, f.Logger)
		}
	case "updateclassicloadbalancer":
		return func() any {
			return awsspec.NewUpdateClassicLoadbalancer(f.config(), f.Graph, f.Logger)
		}
	case "updateconfigrule":
		return func() any {
			return awsspec.NewUpdateConfigrule(f.config(), f.Graph, f.Logger)
		}
	case "updatecontainertask":
		return func() any {
			return awsspec.NewUpdateContainertask(f.config(), f.Graph, f.Logger)
		}
	case "updatedistribution":
		return func() any {
			return awsspec.NewUpdateDistribution(f.config(), f.Graph, f.Logger)
		}
	case "updateenvironment":
		return func() any {
			return awsspec.NewUpdateEnvironment(f.config(), f.Graph, f.Logger)
		}
	case "updateeventrule":
		return func() any {
			return awsspec.NewUpdateEventrule(f.config(), f.Graph, f.Logger)
		}
	case "updateimage":
		return func() any {
			return awsspec.NewUpdateImage(f.config(), f.Graph, f.Logger)
		}
	case "updateinstance":
		return func() any {
			return awsspec.NewUpdateInstance(f.config(), f.Graph, f.Logger)
		}
	case "updateipset":
		return func() any {
			return awsspec.NewUpdateIpset(f.config(), f.Graph, f.Logger)
		}
	case "updateloggroup":
		return func() any {
			return awsspec.NewUpdateLoggroup(f.config(), f.Graph, f.Logger)
		}
	case "updateloginprofile":
		return func() any {
			return awsspec.NewUpdateLoginprofile(f.config(), f.Graph, f.Logger)
		}
	case "updatepolicy":
		return func() any {
			return awsspec.NewUpdatePolicy(f.config(), f.Graph, f.Logger)
		}
	case "updaterecord":
		return func() any {
			return awsspec.NewUpdateRecord(f.config(), f.Graph, f.Logger)
		}
	case "updateredshiftcluster":
		return func() any {
			return awsspec.NewUpdateRedshiftcluster(f.config(), f.Graph, f.Logger)
		}
	case "updates3object":
		return func() any {
			return awsspec.NewUpdateS3object(f.config(), f.Graph, f.Logger)
		}
	case "updatescalinggroup":
		return func() any {
			return awsspec.NewUpdateScalinggroup(f.config(), f.Graph, f.Logger)
		}
	case "updatesecret":
		return func() any {
			return awsspec.NewUpdateSecret(f.config(), f.Graph, f.Logger)
		}
	case "updatesecuritygroup":
		return func() any {
			return awsspec.NewUpdateSecuritygroup(f.config(), f.Graph, f.Logger)
		}
	case "updatessmparameter":
		return func() any {
			return awsspec.NewUpdateSsmparameter(f.config(), f.Graph, f.Logger)
		}
	case "updatestack":
		return func() any {
			return awsspec.NewUpdateStack(f.config(), f.Graph, f.Logger)
		}
	case "updatestatemachine":
		return func() any {
			return awsspec.NewUpdateStatemachine(f.config(), f.Graph, f.Logger)
		}
	case "updatestream":
		return func() any {
			return awsspec.NewUpdateStream(f.config(), f.Graph, f.Logger)
		}
	case "updatesubnet":
		return func() any {
			return awsspec.NewUpdateSubnet(f.config(), f.Graph, f.Logger)
		}
	case "updatetargetgroup":
		return func() any {
			return awsspec.NewUpdateTargetgroup(f.config(), f.Graph, f.Logger)
		}
	}
	return nil
}
