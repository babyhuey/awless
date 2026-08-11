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
package awsspec

import (
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
)

type Factory interface {
	Build(key string) func() any
}

var CommandFactory Factory

var MockAWSSessionFactory = &AWSFactory{
	Log: logger.DiscardLogger,
	Cfg: aws.Config{},
}

type AWSFactory struct {
	Log   *logger.Logger
	Cfg   aws.Config
	Graph cloud.GraphAPI
}

func (f *AWSFactory) Build(key string) func() any {
	switch key {
	case "attachalarm":
		return func() any { return NewAttachAlarm(f.Cfg, f.Graph, f.Log) }
	case "attachclassicloadbalancer":
		return func() any { return NewAttachClassicLoadbalancer(f.Cfg, f.Graph, f.Log) }
	case "attachcontainertask":
		return func() any { return NewAttachContainertask(f.Cfg, f.Graph, f.Log) }
	case "attachelasticip":
		return func() any { return NewAttachElasticip(f.Cfg, f.Graph, f.Log) }
	case "attacheventtarget":
		return func() any { return NewAttachEventtarget(f.Cfg, f.Graph, f.Log) }
	case "attachinstance":
		return func() any { return NewAttachInstance(f.Cfg, f.Graph, f.Log) }
	case "attachinstanceprofile":
		return func() any { return NewAttachInstanceprofile(f.Cfg, f.Graph, f.Log) }
	case "attachinternetgateway":
		return func() any { return NewAttachInternetgateway(f.Cfg, f.Graph, f.Log) }
	case "attachlistener":
		return func() any { return NewAttachListener(f.Cfg, f.Graph, f.Log) }
	case "attachmfadevice":
		return func() any { return NewAttachMfadevice(f.Cfg, f.Graph, f.Log) }
	case "attachnetworkinterface":
		return func() any { return NewAttachNetworkinterface(f.Cfg, f.Graph, f.Log) }
	case "attachpolicy":
		return func() any { return NewAttachPolicy(f.Cfg, f.Graph, f.Log) }
	case "attachrole":
		return func() any { return NewAttachRole(f.Cfg, f.Graph, f.Log) }
	case "attachroutetable":
		return func() any { return NewAttachRoutetable(f.Cfg, f.Graph, f.Log) }
	case "attachsecuritygroup":
		return func() any { return NewAttachSecuritygroup(f.Cfg, f.Graph, f.Log) }
	case "attachuser":
		return func() any { return NewAttachUser(f.Cfg, f.Graph, f.Log) }
	case "attachvolume":
		return func() any { return NewAttachVolume(f.Cfg, f.Graph, f.Log) }
	case "authenticateregistry":
		return func() any { return NewAuthenticateRegistry(f.Cfg, f.Graph, f.Log) }
	case "checkcertificate":
		return func() any { return NewCheckCertificate(f.Cfg, f.Graph, f.Log) }
	case "checkdatabase":
		return func() any { return NewCheckDatabase(f.Cfg, f.Graph, f.Log) }
	case "checkdistribution":
		return func() any { return NewCheckDistribution(f.Cfg, f.Graph, f.Log) }
	case "checkinstance":
		return func() any { return NewCheckInstance(f.Cfg, f.Graph, f.Log) }
	case "checkloadbalancer":
		return func() any { return NewCheckLoadbalancer(f.Cfg, f.Graph, f.Log) }
	case "checknatgateway":
		return func() any { return NewCheckNatgateway(f.Cfg, f.Graph, f.Log) }
	case "checknetworkinterface":
		return func() any { return NewCheckNetworkinterface(f.Cfg, f.Graph, f.Log) }
	case "checkscalinggroup":
		return func() any { return NewCheckScalinggroup(f.Cfg, f.Graph, f.Log) }
	case "checksecuritygroup":
		return func() any { return NewCheckSecuritygroup(f.Cfg, f.Graph, f.Log) }
	case "checkvolume":
		return func() any { return NewCheckVolume(f.Cfg, f.Graph, f.Log) }
	case "copyimage":
		return func() any { return NewCopyImage(f.Cfg, f.Graph, f.Log) }
	case "copysnapshot":
		return func() any { return NewCopySnapshot(f.Cfg, f.Graph, f.Log) }
	case "createaccesskey":
		return func() any { return NewCreateAccesskey(f.Cfg, f.Graph, f.Log) }
	case "createalarm":
		return func() any { return NewCreateAlarm(f.Cfg, f.Graph, f.Log) }
	case "createapigateway":
		return func() any { return NewCreateApigateway(f.Cfg, f.Graph, f.Log) }
	case "createapigatewayroute":
		return func() any { return NewCreateApigatewayroute(f.Cfg, f.Graph, f.Log) }
	case "createapigatewaystage":
		return func() any { return NewCreateApigatewaystage(f.Cfg, f.Graph, f.Log) }
	case "createappscalingpolicy":
		return func() any { return NewCreateAppscalingpolicy(f.Cfg, f.Graph, f.Log) }
	case "createappscalingtarget":
		return func() any { return NewCreateAppscalingtarget(f.Cfg, f.Graph, f.Log) }
	case "createbucket":
		return func() any { return NewCreateBucket(f.Cfg, f.Graph, f.Log) }
	case "createcachecluster":
		return func() any { return NewCreateCachecluster(f.Cfg, f.Graph, f.Log) }
	case "createcachesubnetgroup":
		return func() any { return NewCreateCachesubnetgroup(f.Cfg, f.Graph, f.Log) }
	case "createcertificate":
		return func() any { return NewCreateCertificate(f.Cfg, f.Graph, f.Log) }
	case "createclassicloadbalancer":
		return func() any { return NewCreateClassicLoadbalancer(f.Cfg, f.Graph, f.Log) }
	case "createcontainercluster":
		return func() any { return NewCreateContainercluster(f.Cfg, f.Graph, f.Log) }
	case "createdatabase":
		return func() any { return NewCreateDatabase(f.Cfg, f.Graph, f.Log) }
	case "createdbsubnetgroup":
		return func() any { return NewCreateDbsubnetgroup(f.Cfg, f.Graph, f.Log) }
	case "createdistribution":
		return func() any { return NewCreateDistribution(f.Cfg, f.Graph, f.Log) }
	case "createdynamodbtable":
		return func() any { return NewCreateDynamodbtable(f.Cfg, f.Graph, f.Log) }
	case "createekscluster":
		return func() any { return NewCreateEkscluster(f.Cfg, f.Graph, f.Log) }
	case "createeksnodegroup":
		return func() any { return NewCreateEksnodegroup(f.Cfg, f.Graph, f.Log) }
	case "createelasticip":
		return func() any { return NewCreateElasticip(f.Cfg, f.Graph, f.Log) }
	case "createeventbus":
		return func() any { return NewCreateEventbus(f.Cfg, f.Graph, f.Log) }
	case "createeventrule":
		return func() any { return NewCreateEventrule(f.Cfg, f.Graph, f.Log) }
	case "createfilesystem":
		return func() any { return NewCreateFilesystem(f.Cfg, f.Graph, f.Log) }
	case "createfunction":
		return func() any { return NewCreateFunction(f.Cfg, f.Graph, f.Log) }
	case "creategroup":
		return func() any { return NewCreateGroup(f.Cfg, f.Graph, f.Log) }
	case "createimage":
		return func() any { return NewCreateImage(f.Cfg, f.Graph, f.Log) }
	case "createinstance":
		return func() any { return NewCreateInstance(f.Cfg, f.Graph, f.Log) }
	case "createinstanceprofile":
		return func() any { return NewCreateInstanceprofile(f.Cfg, f.Graph, f.Log) }
	case "createinternetgateway":
		return func() any { return NewCreateInternetgateway(f.Cfg, f.Graph, f.Log) }
	case "createipset":
		return func() any { return NewCreateIpset(f.Cfg, f.Graph, f.Log) }
	case "createkeypair":
		return func() any { return NewCreateKeypair(f.Cfg, f.Graph, f.Log) }
	case "createlaunchconfiguration":
		return func() any { return NewCreateLaunchconfiguration(f.Cfg, f.Graph, f.Log) }
	case "createlistener":
		return func() any { return NewCreateListener(f.Cfg, f.Graph, f.Log) }
	case "createloadbalancer":
		return func() any { return NewCreateLoadbalancer(f.Cfg, f.Graph, f.Log) }
	case "createloggroup":
		return func() any { return NewCreateLoggroup(f.Cfg, f.Graph, f.Log) }
	case "createloginprofile":
		return func() any { return NewCreateLoginprofile(f.Cfg, f.Graph, f.Log) }
	case "createmfadevice":
		return func() any { return NewCreateMfadevice(f.Cfg, f.Graph, f.Log) }
	case "createnatgateway":
		return func() any { return NewCreateNatgateway(f.Cfg, f.Graph, f.Log) }
	case "createnetworkinterface":
		return func() any { return NewCreateNetworkinterface(f.Cfg, f.Graph, f.Log) }
	case "createpolicy":
		return func() any { return NewCreatePolicy(f.Cfg, f.Graph, f.Log) }
	case "createqueue":
		return func() any { return NewCreateQueue(f.Cfg, f.Graph, f.Log) }
	case "createrecord":
		return func() any { return NewCreateRecord(f.Cfg, f.Graph, f.Log) }
	case "createreplicationgroup":
		return func() any { return NewCreateReplicationgroup(f.Cfg, f.Graph, f.Log) }
	case "createrepository":
		return func() any { return NewCreateRepository(f.Cfg, f.Graph, f.Log) }
	case "createrole":
		return func() any { return NewCreateRole(f.Cfg, f.Graph, f.Log) }
	case "createroute":
		return func() any { return NewCreateRoute(f.Cfg, f.Graph, f.Log) }
	case "createroutetable":
		return func() any { return NewCreateRoutetable(f.Cfg, f.Graph, f.Log) }
	case "creates3object":
		return func() any { return NewCreateS3object(f.Cfg, f.Graph, f.Log) }
	case "createscalinggroup":
		return func() any { return NewCreateScalinggroup(f.Cfg, f.Graph, f.Log) }
	case "createscalingpolicy":
		return func() any { return NewCreateScalingpolicy(f.Cfg, f.Graph, f.Log) }
	case "createsecret":
		return func() any { return NewCreateSecret(f.Cfg, f.Graph, f.Log) }
	case "createsecuritygroup":
		return func() any { return NewCreateSecuritygroup(f.Cfg, f.Graph, f.Log) }
	case "createsnapshot":
		return func() any { return NewCreateSnapshot(f.Cfg, f.Graph, f.Log) }
	case "createssmparameter":
		return func() any { return NewCreateSsmparameter(f.Cfg, f.Graph, f.Log) }
	case "createstack":
		return func() any { return NewCreateStack(f.Cfg, f.Graph, f.Log) }
	case "createstatemachine":
		return func() any { return NewCreateStatemachine(f.Cfg, f.Graph, f.Log) }
	case "createsubnet":
		return func() any { return NewCreateSubnet(f.Cfg, f.Graph, f.Log) }
	case "createsubscription":
		return func() any { return NewCreateSubscription(f.Cfg, f.Graph, f.Log) }
	case "createtag":
		return func() any { return NewCreateTag(f.Cfg, f.Graph, f.Log) }
	case "createtargetgroup":
		return func() any { return NewCreateTargetgroup(f.Cfg, f.Graph, f.Log) }
	case "createtopic":
		return func() any { return NewCreateTopic(f.Cfg, f.Graph, f.Log) }
	case "createtrail":
		return func() any { return NewCreateTrail(f.Cfg, f.Graph, f.Log) }
	case "createuser":
		return func() any { return NewCreateUser(f.Cfg, f.Graph, f.Log) }
	case "createvolume":
		return func() any { return NewCreateVolume(f.Cfg, f.Graph, f.Log) }
	case "createvpc":
		return func() any { return NewCreateVpc(f.Cfg, f.Graph, f.Log) }
	case "createzone":
		return func() any { return NewCreateZone(f.Cfg, f.Graph, f.Log) }
	case "deleteaccesskey":
		return func() any { return NewDeleteAccesskey(f.Cfg, f.Graph, f.Log) }
	case "deletealarm":
		return func() any { return NewDeleteAlarm(f.Cfg, f.Graph, f.Log) }
	case "deleteapigateway":
		return func() any { return NewDeleteApigateway(f.Cfg, f.Graph, f.Log) }
	case "deleteapigatewayroute":
		return func() any { return NewDeleteApigatewayroute(f.Cfg, f.Graph, f.Log) }
	case "deleteapigatewaystage":
		return func() any { return NewDeleteApigatewaystage(f.Cfg, f.Graph, f.Log) }
	case "deleteappscalingpolicy":
		return func() any { return NewDeleteAppscalingpolicy(f.Cfg, f.Graph, f.Log) }
	case "deleteappscalingtarget":
		return func() any { return NewDeleteAppscalingtarget(f.Cfg, f.Graph, f.Log) }
	case "deletebucket":
		return func() any { return NewDeleteBucket(f.Cfg, f.Graph, f.Log) }
	case "deletecachecluster":
		return func() any { return NewDeleteCachecluster(f.Cfg, f.Graph, f.Log) }
	case "deletecachesubnetgroup":
		return func() any { return NewDeleteCachesubnetgroup(f.Cfg, f.Graph, f.Log) }
	case "deletecertificate":
		return func() any { return NewDeleteCertificate(f.Cfg, f.Graph, f.Log) }
	case "deleteclassicloadbalancer":
		return func() any { return NewDeleteClassicLoadbalancer(f.Cfg, f.Graph, f.Log) }
	case "deletecontainercluster":
		return func() any { return NewDeleteContainercluster(f.Cfg, f.Graph, f.Log) }
	case "deletecontainertask":
		return func() any { return NewDeleteContainertask(f.Cfg, f.Graph, f.Log) }
	case "deletedatabase":
		return func() any { return NewDeleteDatabase(f.Cfg, f.Graph, f.Log) }
	case "deletedbsubnetgroup":
		return func() any { return NewDeleteDbsubnetgroup(f.Cfg, f.Graph, f.Log) }
	case "deletedistribution":
		return func() any { return NewDeleteDistribution(f.Cfg, f.Graph, f.Log) }
	case "deletedynamodbtable":
		return func() any { return NewDeleteDynamodbtable(f.Cfg, f.Graph, f.Log) }
	case "deleteekscluster":
		return func() any { return NewDeleteEkscluster(f.Cfg, f.Graph, f.Log) }
	case "deleteeksnodegroup":
		return func() any { return NewDeleteEksnodegroup(f.Cfg, f.Graph, f.Log) }
	case "deleteelasticip":
		return func() any { return NewDeleteElasticip(f.Cfg, f.Graph, f.Log) }
	case "deleteeventbus":
		return func() any { return NewDeleteEventbus(f.Cfg, f.Graph, f.Log) }
	case "deleteeventrule":
		return func() any { return NewDeleteEventrule(f.Cfg, f.Graph, f.Log) }
	case "deletefilesystem":
		return func() any { return NewDeleteFilesystem(f.Cfg, f.Graph, f.Log) }
	case "deletefunction":
		return func() any { return NewDeleteFunction(f.Cfg, f.Graph, f.Log) }
	case "deletegroup":
		return func() any { return NewDeleteGroup(f.Cfg, f.Graph, f.Log) }
	case "deleteimage":
		return func() any { return NewDeleteImage(f.Cfg, f.Graph, f.Log) }
	case "deleteinstance":
		return func() any { return NewDeleteInstance(f.Cfg, f.Graph, f.Log) }
	case "deleteinstanceprofile":
		return func() any { return NewDeleteInstanceprofile(f.Cfg, f.Graph, f.Log) }
	case "deleteinternetgateway":
		return func() any { return NewDeleteInternetgateway(f.Cfg, f.Graph, f.Log) }
	case "deleteipset":
		return func() any { return NewDeleteIpset(f.Cfg, f.Graph, f.Log) }
	case "deletekeypair":
		return func() any { return NewDeleteKeypair(f.Cfg, f.Graph, f.Log) }
	case "deletelaunchconfiguration":
		return func() any { return NewDeleteLaunchconfiguration(f.Cfg, f.Graph, f.Log) }
	case "deletelistener":
		return func() any { return NewDeleteListener(f.Cfg, f.Graph, f.Log) }
	case "deleteloadbalancer":
		return func() any { return NewDeleteLoadbalancer(f.Cfg, f.Graph, f.Log) }
	case "deleteloggroup":
		return func() any { return NewDeleteLoggroup(f.Cfg, f.Graph, f.Log) }
	case "deleteloginprofile":
		return func() any { return NewDeleteLoginprofile(f.Cfg, f.Graph, f.Log) }
	case "deletemfadevice":
		return func() any { return NewDeleteMfadevice(f.Cfg, f.Graph, f.Log) }
	case "deletenatgateway":
		return func() any { return NewDeleteNatgateway(f.Cfg, f.Graph, f.Log) }
	case "deletenetworkinterface":
		return func() any { return NewDeleteNetworkinterface(f.Cfg, f.Graph, f.Log) }
	case "deletepolicy":
		return func() any { return NewDeletePolicy(f.Cfg, f.Graph, f.Log) }
	case "deletequeue":
		return func() any { return NewDeleteQueue(f.Cfg, f.Graph, f.Log) }
	case "deleterecord":
		return func() any { return NewDeleteRecord(f.Cfg, f.Graph, f.Log) }
	case "deletereplicationgroup":
		return func() any { return NewDeleteReplicationgroup(f.Cfg, f.Graph, f.Log) }
	case "deleterepository":
		return func() any { return NewDeleteRepository(f.Cfg, f.Graph, f.Log) }
	case "deleterole":
		return func() any { return NewDeleteRole(f.Cfg, f.Graph, f.Log) }
	case "deleteroute":
		return func() any { return NewDeleteRoute(f.Cfg, f.Graph, f.Log) }
	case "deleteroutetable":
		return func() any { return NewDeleteRoutetable(f.Cfg, f.Graph, f.Log) }
	case "deletes3object":
		return func() any { return NewDeleteS3object(f.Cfg, f.Graph, f.Log) }
	case "deletescalinggroup":
		return func() any { return NewDeleteScalinggroup(f.Cfg, f.Graph, f.Log) }
	case "deletescalingpolicy":
		return func() any { return NewDeleteScalingpolicy(f.Cfg, f.Graph, f.Log) }
	case "deletesecret":
		return func() any { return NewDeleteSecret(f.Cfg, f.Graph, f.Log) }
	case "deletesecuritygroup":
		return func() any { return NewDeleteSecuritygroup(f.Cfg, f.Graph, f.Log) }
	case "deletesnapshot":
		return func() any { return NewDeleteSnapshot(f.Cfg, f.Graph, f.Log) }
	case "deletessmparameter":
		return func() any { return NewDeleteSsmparameter(f.Cfg, f.Graph, f.Log) }
	case "deletestack":
		return func() any { return NewDeleteStack(f.Cfg, f.Graph, f.Log) }
	case "deletestatemachine":
		return func() any { return NewDeleteStatemachine(f.Cfg, f.Graph, f.Log) }
	case "deletesubnet":
		return func() any { return NewDeleteSubnet(f.Cfg, f.Graph, f.Log) }
	case "deletesubscription":
		return func() any { return NewDeleteSubscription(f.Cfg, f.Graph, f.Log) }
	case "deletetag":
		return func() any { return NewDeleteTag(f.Cfg, f.Graph, f.Log) }
	case "deletetargetgroup":
		return func() any { return NewDeleteTargetgroup(f.Cfg, f.Graph, f.Log) }
	case "deletetopic":
		return func() any { return NewDeleteTopic(f.Cfg, f.Graph, f.Log) }
	case "deletetrail":
		return func() any { return NewDeleteTrail(f.Cfg, f.Graph, f.Log) }
	case "deleteuser":
		return func() any { return NewDeleteUser(f.Cfg, f.Graph, f.Log) }
	case "deletevolume":
		return func() any { return NewDeleteVolume(f.Cfg, f.Graph, f.Log) }
	case "deletevpc":
		return func() any { return NewDeleteVpc(f.Cfg, f.Graph, f.Log) }
	case "deletezone":
		return func() any { return NewDeleteZone(f.Cfg, f.Graph, f.Log) }
	case "detachalarm":
		return func() any { return NewDetachAlarm(f.Cfg, f.Graph, f.Log) }
	case "detachclassicloadbalancer":
		return func() any { return NewDetachClassicLoadbalancer(f.Cfg, f.Graph, f.Log) }
	case "detachcontainertask":
		return func() any { return NewDetachContainertask(f.Cfg, f.Graph, f.Log) }
	case "detachelasticip":
		return func() any { return NewDetachElasticip(f.Cfg, f.Graph, f.Log) }
	case "detacheventtarget":
		return func() any { return NewDetachEventtarget(f.Cfg, f.Graph, f.Log) }
	case "detachinstance":
		return func() any { return NewDetachInstance(f.Cfg, f.Graph, f.Log) }
	case "detachinstanceprofile":
		return func() any { return NewDetachInstanceprofile(f.Cfg, f.Graph, f.Log) }
	case "detachinternetgateway":
		return func() any { return NewDetachInternetgateway(f.Cfg, f.Graph, f.Log) }
	case "detachmfadevice":
		return func() any { return NewDetachMfadevice(f.Cfg, f.Graph, f.Log) }
	case "detachnetworkinterface":
		return func() any { return NewDetachNetworkinterface(f.Cfg, f.Graph, f.Log) }
	case "detachpolicy":
		return func() any { return NewDetachPolicy(f.Cfg, f.Graph, f.Log) }
	case "detachrole":
		return func() any { return NewDetachRole(f.Cfg, f.Graph, f.Log) }
	case "detachroutetable":
		return func() any { return NewDetachRoutetable(f.Cfg, f.Graph, f.Log) }
	case "detachsecuritygroup":
		return func() any { return NewDetachSecuritygroup(f.Cfg, f.Graph, f.Log) }
	case "detachuser":
		return func() any { return NewDetachUser(f.Cfg, f.Graph, f.Log) }
	case "detachvolume":
		return func() any { return NewDetachVolume(f.Cfg, f.Graph, f.Log) }
	case "importimage":
		return func() any { return NewImportImage(f.Cfg, f.Graph, f.Log) }
	case "restartdatabase":
		return func() any { return NewRestartDatabase(f.Cfg, f.Graph, f.Log) }
	case "restartinstance":
		return func() any { return NewRestartInstance(f.Cfg, f.Graph, f.Log) }
	case "startalarm":
		return func() any { return NewStartAlarm(f.Cfg, f.Graph, f.Log) }
	case "startcontainertask":
		return func() any { return NewStartContainertask(f.Cfg, f.Graph, f.Log) }
	case "startdatabase":
		return func() any { return NewStartDatabase(f.Cfg, f.Graph, f.Log) }
	case "starteventrule":
		return func() any { return NewStartEventrule(f.Cfg, f.Graph, f.Log) }
	case "startexecution":
		return func() any { return NewStartExecution(f.Cfg, f.Graph, f.Log) }
	case "startinstance":
		return func() any { return NewStartInstance(f.Cfg, f.Graph, f.Log) }
	case "starttrail":
		return func() any { return NewStartTrail(f.Cfg, f.Graph, f.Log) }
	case "stopalarm":
		return func() any { return NewStopAlarm(f.Cfg, f.Graph, f.Log) }
	case "stopcontainertask":
		return func() any { return NewStopContainertask(f.Cfg, f.Graph, f.Log) }
	case "stopdatabase":
		return func() any { return NewStopDatabase(f.Cfg, f.Graph, f.Log) }
	case "stopeventrule":
		return func() any { return NewStopEventrule(f.Cfg, f.Graph, f.Log) }
	case "stopexecution":
		return func() any { return NewStopExecution(f.Cfg, f.Graph, f.Log) }
	case "stopinstance":
		return func() any { return NewStopInstance(f.Cfg, f.Graph, f.Log) }
	case "stoptrail":
		return func() any { return NewStopTrail(f.Cfg, f.Graph, f.Log) }
	case "updatebucket":
		return func() any { return NewUpdateBucket(f.Cfg, f.Graph, f.Log) }
	case "updatecachecluster":
		return func() any { return NewUpdateCachecluster(f.Cfg, f.Graph, f.Log) }
	case "updatecachesubnetgroup":
		return func() any { return NewUpdateCachesubnetgroup(f.Cfg, f.Graph, f.Log) }
	case "updateclassicloadbalancer":
		return func() any { return NewUpdateClassicLoadbalancer(f.Cfg, f.Graph, f.Log) }
	case "updatecontainertask":
		return func() any { return NewUpdateContainertask(f.Cfg, f.Graph, f.Log) }
	case "updatedistribution":
		return func() any { return NewUpdateDistribution(f.Cfg, f.Graph, f.Log) }
	case "updateeventrule":
		return func() any { return NewUpdateEventrule(f.Cfg, f.Graph, f.Log) }
	case "updateimage":
		return func() any { return NewUpdateImage(f.Cfg, f.Graph, f.Log) }
	case "updateinstance":
		return func() any { return NewUpdateInstance(f.Cfg, f.Graph, f.Log) }
	case "updateipset":
		return func() any { return NewUpdateIpset(f.Cfg, f.Graph, f.Log) }
	case "updateloggroup":
		return func() any { return NewUpdateLoggroup(f.Cfg, f.Graph, f.Log) }
	case "updateloginprofile":
		return func() any { return NewUpdateLoginprofile(f.Cfg, f.Graph, f.Log) }
	case "updatepolicy":
		return func() any { return NewUpdatePolicy(f.Cfg, f.Graph, f.Log) }
	case "updaterecord":
		return func() any { return NewUpdateRecord(f.Cfg, f.Graph, f.Log) }
	case "updates3object":
		return func() any { return NewUpdateS3object(f.Cfg, f.Graph, f.Log) }
	case "updatescalinggroup":
		return func() any { return NewUpdateScalinggroup(f.Cfg, f.Graph, f.Log) }
	case "updatesecret":
		return func() any { return NewUpdateSecret(f.Cfg, f.Graph, f.Log) }
	case "updatesecuritygroup":
		return func() any { return NewUpdateSecuritygroup(f.Cfg, f.Graph, f.Log) }
	case "updatessmparameter":
		return func() any { return NewUpdateSsmparameter(f.Cfg, f.Graph, f.Log) }
	case "updatestack":
		return func() any { return NewUpdateStack(f.Cfg, f.Graph, f.Log) }
	case "updatestatemachine":
		return func() any { return NewUpdateStatemachine(f.Cfg, f.Graph, f.Log) }
	case "updatesubnet":
		return func() any { return NewUpdateSubnet(f.Cfg, f.Graph, f.Log) }
	case "updatetargetgroup":
		return func() any { return NewUpdateTargetgroup(f.Cfg, f.Graph, f.Log) }
	}
	return nil
}

var (
	_ command = &AttachAlarm{}
	_ command = &AttachClassicLoadbalancer{}
	_ command = &AttachContainertask{}
	_ command = &AttachElasticip{}
	_ command = &AttachEventtarget{}
	_ command = &AttachInstance{}
	_ command = &AttachInstanceprofile{}
	_ command = &AttachInternetgateway{}
	_ command = &AttachListener{}
	_ command = &AttachMfadevice{}
	_ command = &AttachNetworkinterface{}
	_ command = &AttachPolicy{}
	_ command = &AttachRole{}
	_ command = &AttachRoutetable{}
	_ command = &AttachSecuritygroup{}
	_ command = &AttachUser{}
	_ command = &AttachVolume{}
	_ command = &AuthenticateRegistry{}
	_ command = &CheckCertificate{}
	_ command = &CheckDatabase{}
	_ command = &CheckDistribution{}
	_ command = &CheckInstance{}
	_ command = &CheckLoadbalancer{}
	_ command = &CheckNatgateway{}
	_ command = &CheckNetworkinterface{}
	_ command = &CheckScalinggroup{}
	_ command = &CheckSecuritygroup{}
	_ command = &CheckVolume{}
	_ command = &CopyImage{}
	_ command = &CopySnapshot{}
	_ command = &CreateAccesskey{}
	_ command = &CreateAlarm{}
	_ command = &CreateApigateway{}
	_ command = &CreateApigatewayroute{}
	_ command = &CreateApigatewaystage{}
	_ command = &CreateAppscalingpolicy{}
	_ command = &CreateAppscalingtarget{}
	_ command = &CreateBucket{}
	_ command = &CreateCachecluster{}
	_ command = &CreateCachesubnetgroup{}
	_ command = &CreateCertificate{}
	_ command = &CreateClassicLoadbalancer{}
	_ command = &CreateContainercluster{}
	_ command = &CreateDatabase{}
	_ command = &CreateDbsubnetgroup{}
	_ command = &CreateDistribution{}
	_ command = &CreateDynamodbtable{}
	_ command = &CreateEkscluster{}
	_ command = &CreateEksnodegroup{}
	_ command = &CreateElasticip{}
	_ command = &CreateEventbus{}
	_ command = &CreateEventrule{}
	_ command = &CreateFilesystem{}
	_ command = &CreateFunction{}
	_ command = &CreateGroup{}
	_ command = &CreateImage{}
	_ command = &CreateInstance{}
	_ command = &CreateInstanceprofile{}
	_ command = &CreateInternetgateway{}
	_ command = &CreateIpset{}
	_ command = &CreateKeypair{}
	_ command = &CreateLaunchconfiguration{}
	_ command = &CreateListener{}
	_ command = &CreateLoadbalancer{}
	_ command = &CreateLoggroup{}
	_ command = &CreateLoginprofile{}
	_ command = &CreateMfadevice{}
	_ command = &CreateNatgateway{}
	_ command = &CreateNetworkinterface{}
	_ command = &CreatePolicy{}
	_ command = &CreateQueue{}
	_ command = &CreateRecord{}
	_ command = &CreateReplicationgroup{}
	_ command = &CreateRepository{}
	_ command = &CreateRole{}
	_ command = &CreateRoute{}
	_ command = &CreateRoutetable{}
	_ command = &CreateS3object{}
	_ command = &CreateScalinggroup{}
	_ command = &CreateScalingpolicy{}
	_ command = &CreateSecret{}
	_ command = &CreateSecuritygroup{}
	_ command = &CreateSnapshot{}
	_ command = &CreateSsmparameter{}
	_ command = &CreateStack{}
	_ command = &CreateStatemachine{}
	_ command = &CreateSubnet{}
	_ command = &CreateSubscription{}
	_ command = &CreateTag{}
	_ command = &CreateTargetgroup{}
	_ command = &CreateTopic{}
	_ command = &CreateTrail{}
	_ command = &CreateUser{}
	_ command = &CreateVolume{}
	_ command = &CreateVpc{}
	_ command = &CreateZone{}
	_ command = &DeleteAccesskey{}
	_ command = &DeleteAlarm{}
	_ command = &DeleteApigateway{}
	_ command = &DeleteApigatewayroute{}
	_ command = &DeleteApigatewaystage{}
	_ command = &DeleteAppscalingpolicy{}
	_ command = &DeleteAppscalingtarget{}
	_ command = &DeleteBucket{}
	_ command = &DeleteCachecluster{}
	_ command = &DeleteCachesubnetgroup{}
	_ command = &DeleteCertificate{}
	_ command = &DeleteClassicLoadbalancer{}
	_ command = &DeleteContainercluster{}
	_ command = &DeleteContainertask{}
	_ command = &DeleteDatabase{}
	_ command = &DeleteDbsubnetgroup{}
	_ command = &DeleteDistribution{}
	_ command = &DeleteDynamodbtable{}
	_ command = &DeleteEkscluster{}
	_ command = &DeleteEksnodegroup{}
	_ command = &DeleteElasticip{}
	_ command = &DeleteEventbus{}
	_ command = &DeleteEventrule{}
	_ command = &DeleteFilesystem{}
	_ command = &DeleteFunction{}
	_ command = &DeleteGroup{}
	_ command = &DeleteImage{}
	_ command = &DeleteInstance{}
	_ command = &DeleteInstanceprofile{}
	_ command = &DeleteInternetgateway{}
	_ command = &DeleteIpset{}
	_ command = &DeleteKeypair{}
	_ command = &DeleteLaunchconfiguration{}
	_ command = &DeleteListener{}
	_ command = &DeleteLoadbalancer{}
	_ command = &DeleteLoggroup{}
	_ command = &DeleteLoginprofile{}
	_ command = &DeleteMfadevice{}
	_ command = &DeleteNatgateway{}
	_ command = &DeleteNetworkinterface{}
	_ command = &DeletePolicy{}
	_ command = &DeleteQueue{}
	_ command = &DeleteRecord{}
	_ command = &DeleteReplicationgroup{}
	_ command = &DeleteRepository{}
	_ command = &DeleteRole{}
	_ command = &DeleteRoute{}
	_ command = &DeleteRoutetable{}
	_ command = &DeleteS3object{}
	_ command = &DeleteScalinggroup{}
	_ command = &DeleteScalingpolicy{}
	_ command = &DeleteSecret{}
	_ command = &DeleteSecuritygroup{}
	_ command = &DeleteSnapshot{}
	_ command = &DeleteSsmparameter{}
	_ command = &DeleteStack{}
	_ command = &DeleteStatemachine{}
	_ command = &DeleteSubnet{}
	_ command = &DeleteSubscription{}
	_ command = &DeleteTag{}
	_ command = &DeleteTargetgroup{}
	_ command = &DeleteTopic{}
	_ command = &DeleteTrail{}
	_ command = &DeleteUser{}
	_ command = &DeleteVolume{}
	_ command = &DeleteVpc{}
	_ command = &DeleteZone{}
	_ command = &DetachAlarm{}
	_ command = &DetachClassicLoadbalancer{}
	_ command = &DetachContainertask{}
	_ command = &DetachElasticip{}
	_ command = &DetachEventtarget{}
	_ command = &DetachInstance{}
	_ command = &DetachInstanceprofile{}
	_ command = &DetachInternetgateway{}
	_ command = &DetachMfadevice{}
	_ command = &DetachNetworkinterface{}
	_ command = &DetachPolicy{}
	_ command = &DetachRole{}
	_ command = &DetachRoutetable{}
	_ command = &DetachSecuritygroup{}
	_ command = &DetachUser{}
	_ command = &DetachVolume{}
	_ command = &ImportImage{}
	_ command = &RestartDatabase{}
	_ command = &RestartInstance{}
	_ command = &StartAlarm{}
	_ command = &StartContainertask{}
	_ command = &StartDatabase{}
	_ command = &StartEventrule{}
	_ command = &StartExecution{}
	_ command = &StartInstance{}
	_ command = &StartTrail{}
	_ command = &StopAlarm{}
	_ command = &StopContainertask{}
	_ command = &StopDatabase{}
	_ command = &StopEventrule{}
	_ command = &StopExecution{}
	_ command = &StopInstance{}
	_ command = &StopTrail{}
	_ command = &UpdateBucket{}
	_ command = &UpdateCachecluster{}
	_ command = &UpdateCachesubnetgroup{}
	_ command = &UpdateClassicLoadbalancer{}
	_ command = &UpdateContainertask{}
	_ command = &UpdateDistribution{}
	_ command = &UpdateEventrule{}
	_ command = &UpdateImage{}
	_ command = &UpdateInstance{}
	_ command = &UpdateIpset{}
	_ command = &UpdateLoggroup{}
	_ command = &UpdateLoginprofile{}
	_ command = &UpdatePolicy{}
	_ command = &UpdateRecord{}
	_ command = &UpdateS3object{}
	_ command = &UpdateScalinggroup{}
	_ command = &UpdateSecret{}
	_ command = &UpdateSecuritygroup{}
	_ command = &UpdateSsmparameter{}
	_ command = &UpdateStack{}
	_ command = &UpdateStatemachine{}
	_ command = &UpdateSubnet{}
	_ command = &UpdateTargetgroup{}
)
