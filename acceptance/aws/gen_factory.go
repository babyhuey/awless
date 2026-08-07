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

type AcceptanceFactory struct {
	Mock   any
	Logger *logger.Logger
	Graph  cloud.GraphAPI
}

func NewAcceptanceFactory(mock any, g cloud.GraphAPI, l ...*logger.Logger) *AcceptanceFactory {
	lg := logger.DiscardLogger
	if len(l) > 0 {
		lg = l[0]
	}
	return &AcceptanceFactory{Mock: mock, Graph: g, Logger: lg}
}

func (f *AcceptanceFactory) Build(key string) func() any {
	switch key {
	case "attachalarm":
		return func() any {
			cmd := awsspec.NewAttachAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachclassicloadbalancer":
		return func() any {
			cmd := awsspec.NewAttachClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachcontainertask":
		return func() any {
			cmd := awsspec.NewAttachContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachelasticip":
		return func() any {
			cmd := awsspec.NewAttachElasticip(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachinstance":
		return func() any {
			cmd := awsspec.NewAttachInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachinstanceprofile":
		return func() any {
			cmd := awsspec.NewAttachInstanceprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachinternetgateway":
		return func() any {
			cmd := awsspec.NewAttachInternetgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachlistener":
		return func() any {
			cmd := awsspec.NewAttachListener(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachmfadevice":
		return func() any {
			cmd := awsspec.NewAttachMfadevice(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachnetworkinterface":
		return func() any {
			cmd := awsspec.NewAttachNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachpolicy":
		return func() any {
			cmd := awsspec.NewAttachPolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachrole":
		return func() any {
			cmd := awsspec.NewAttachRole(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachroutetable":
		return func() any {
			cmd := awsspec.NewAttachRoutetable(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachsecuritygroup":
		return func() any {
			cmd := awsspec.NewAttachSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachuser":
		return func() any {
			cmd := awsspec.NewAttachUser(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachvolume":
		return func() any {
			cmd := awsspec.NewAttachVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "authenticateregistry":
		return func() any {
			cmd := awsspec.NewAuthenticateRegistry(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkcertificate":
		return func() any {
			cmd := awsspec.NewCheckCertificate(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkdatabase":
		return func() any {
			cmd := awsspec.NewCheckDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkdistribution":
		return func() any {
			cmd := awsspec.NewCheckDistribution(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkinstance":
		return func() any {
			cmd := awsspec.NewCheckInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkloadbalancer":
		return func() any {
			cmd := awsspec.NewCheckLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checknatgateway":
		return func() any {
			cmd := awsspec.NewCheckNatgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checknetworkinterface":
		return func() any {
			cmd := awsspec.NewCheckNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkscalinggroup":
		return func() any {
			cmd := awsspec.NewCheckScalinggroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checksecuritygroup":
		return func() any {
			cmd := awsspec.NewCheckSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkvolume":
		return func() any {
			cmd := awsspec.NewCheckVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "copyimage":
		return func() any {
			cmd := awsspec.NewCopyImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "copysnapshot":
		return func() any {
			cmd := awsspec.NewCopySnapshot(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createaccesskey":
		return func() any {
			cmd := awsspec.NewCreateAccesskey(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createalarm":
		return func() any {
			cmd := awsspec.NewCreateAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createappscalingpolicy":
		return func() any {
			cmd := awsspec.NewCreateAppscalingpolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createappscalingtarget":
		return func() any {
			cmd := awsspec.NewCreateAppscalingtarget(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createbucket":
		return func() any {
			cmd := awsspec.NewCreateBucket(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createcertificate":
		return func() any {
			cmd := awsspec.NewCreateCertificate(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createclassicloadbalancer":
		return func() any {
			cmd := awsspec.NewCreateClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createcontainercluster":
		return func() any {
			cmd := awsspec.NewCreateContainercluster(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createdatabase":
		return func() any {
			cmd := awsspec.NewCreateDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createdbsubnetgroup":
		return func() any {
			cmd := awsspec.NewCreateDbsubnetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createdistribution":
		return func() any {
			cmd := awsspec.NewCreateDistribution(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createelasticip":
		return func() any {
			cmd := awsspec.NewCreateElasticip(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createfunction":
		return func() any {
			cmd := awsspec.NewCreateFunction(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "creategroup":
		return func() any {
			cmd := awsspec.NewCreateGroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createimage":
		return func() any {
			cmd := awsspec.NewCreateImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createinstance":
		return func() any {
			cmd := awsspec.NewCreateInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createinstanceprofile":
		return func() any {
			cmd := awsspec.NewCreateInstanceprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createinternetgateway":
		return func() any {
			cmd := awsspec.NewCreateInternetgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createkeypair":
		return func() any {
			cmd := awsspec.NewCreateKeypair(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createlaunchconfiguration":
		return func() any {
			cmd := awsspec.NewCreateLaunchconfiguration(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createlistener":
		return func() any {
			cmd := awsspec.NewCreateListener(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createloadbalancer":
		return func() any {
			cmd := awsspec.NewCreateLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createloginprofile":
		return func() any {
			cmd := awsspec.NewCreateLoginprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createmfadevice":
		return func() any {
			cmd := awsspec.NewCreateMfadevice(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createnatgateway":
		return func() any {
			cmd := awsspec.NewCreateNatgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createnetworkinterface":
		return func() any {
			cmd := awsspec.NewCreateNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createpolicy":
		return func() any {
			cmd := awsspec.NewCreatePolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createqueue":
		return func() any {
			cmd := awsspec.NewCreateQueue(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createrecord":
		return func() any {
			cmd := awsspec.NewCreateRecord(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createrepository":
		return func() any {
			cmd := awsspec.NewCreateRepository(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createrole":
		return func() any {
			cmd := awsspec.NewCreateRole(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createroute":
		return func() any {
			cmd := awsspec.NewCreateRoute(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createroutetable":
		return func() any {
			cmd := awsspec.NewCreateRoutetable(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "creates3object":
		return func() any {
			cmd := awsspec.NewCreateS3object(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createscalinggroup":
		return func() any {
			cmd := awsspec.NewCreateScalinggroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createscalingpolicy":
		return func() any {
			cmd := awsspec.NewCreateScalingpolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createsecuritygroup":
		return func() any {
			cmd := awsspec.NewCreateSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createsnapshot":
		return func() any {
			cmd := awsspec.NewCreateSnapshot(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createstack":
		return func() any {
			cmd := awsspec.NewCreateStack(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createsubnet":
		return func() any {
			cmd := awsspec.NewCreateSubnet(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createsubscription":
		return func() any {
			cmd := awsspec.NewCreateSubscription(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createtag":
		return func() any {
			cmd := awsspec.NewCreateTag(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createtargetgroup":
		return func() any {
			cmd := awsspec.NewCreateTargetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createtopic":
		return func() any {
			cmd := awsspec.NewCreateTopic(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createuser":
		return func() any {
			cmd := awsspec.NewCreateUser(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createvolume":
		return func() any {
			cmd := awsspec.NewCreateVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createvpc":
		return func() any {
			cmd := awsspec.NewCreateVpc(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createzone":
		return func() any {
			cmd := awsspec.NewCreateZone(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteaccesskey":
		return func() any {
			cmd := awsspec.NewDeleteAccesskey(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletealarm":
		return func() any {
			cmd := awsspec.NewDeleteAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteappscalingpolicy":
		return func() any {
			cmd := awsspec.NewDeleteAppscalingpolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteappscalingtarget":
		return func() any {
			cmd := awsspec.NewDeleteAppscalingtarget(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletebucket":
		return func() any {
			cmd := awsspec.NewDeleteBucket(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletecertificate":
		return func() any {
			cmd := awsspec.NewDeleteCertificate(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteclassicloadbalancer":
		return func() any {
			cmd := awsspec.NewDeleteClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletecontainercluster":
		return func() any {
			cmd := awsspec.NewDeleteContainercluster(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletecontainertask":
		return func() any {
			cmd := awsspec.NewDeleteContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletedatabase":
		return func() any {
			cmd := awsspec.NewDeleteDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletedbsubnetgroup":
		return func() any {
			cmd := awsspec.NewDeleteDbsubnetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletedistribution":
		return func() any {
			cmd := awsspec.NewDeleteDistribution(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteelasticip":
		return func() any {
			cmd := awsspec.NewDeleteElasticip(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletefunction":
		return func() any {
			cmd := awsspec.NewDeleteFunction(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletegroup":
		return func() any {
			cmd := awsspec.NewDeleteGroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteimage":
		return func() any {
			cmd := awsspec.NewDeleteImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteinstance":
		return func() any {
			cmd := awsspec.NewDeleteInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteinstanceprofile":
		return func() any {
			cmd := awsspec.NewDeleteInstanceprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteinternetgateway":
		return func() any {
			cmd := awsspec.NewDeleteInternetgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletekeypair":
		return func() any {
			cmd := awsspec.NewDeleteKeypair(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletelaunchconfiguration":
		return func() any {
			cmd := awsspec.NewDeleteLaunchconfiguration(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletelistener":
		return func() any {
			cmd := awsspec.NewDeleteListener(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteloadbalancer":
		return func() any {
			cmd := awsspec.NewDeleteLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteloginprofile":
		return func() any {
			cmd := awsspec.NewDeleteLoginprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletemfadevice":
		return func() any {
			cmd := awsspec.NewDeleteMfadevice(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletenatgateway":
		return func() any {
			cmd := awsspec.NewDeleteNatgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletenetworkinterface":
		return func() any {
			cmd := awsspec.NewDeleteNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletepolicy":
		return func() any {
			cmd := awsspec.NewDeletePolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletequeue":
		return func() any {
			cmd := awsspec.NewDeleteQueue(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleterecord":
		return func() any {
			cmd := awsspec.NewDeleteRecord(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleterepository":
		return func() any {
			cmd := awsspec.NewDeleteRepository(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleterole":
		return func() any {
			cmd := awsspec.NewDeleteRole(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteroute":
		return func() any {
			cmd := awsspec.NewDeleteRoute(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteroutetable":
		return func() any {
			cmd := awsspec.NewDeleteRoutetable(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletes3object":
		return func() any {
			cmd := awsspec.NewDeleteS3object(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletescalinggroup":
		return func() any {
			cmd := awsspec.NewDeleteScalinggroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletescalingpolicy":
		return func() any {
			cmd := awsspec.NewDeleteScalingpolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletesecuritygroup":
		return func() any {
			cmd := awsspec.NewDeleteSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletesnapshot":
		return func() any {
			cmd := awsspec.NewDeleteSnapshot(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletestack":
		return func() any {
			cmd := awsspec.NewDeleteStack(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletesubnet":
		return func() any {
			cmd := awsspec.NewDeleteSubnet(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletesubscription":
		return func() any {
			cmd := awsspec.NewDeleteSubscription(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletetag":
		return func() any {
			cmd := awsspec.NewDeleteTag(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletetargetgroup":
		return func() any {
			cmd := awsspec.NewDeleteTargetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletetopic":
		return func() any {
			cmd := awsspec.NewDeleteTopic(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteuser":
		return func() any {
			cmd := awsspec.NewDeleteUser(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletevolume":
		return func() any {
			cmd := awsspec.NewDeleteVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletevpc":
		return func() any {
			cmd := awsspec.NewDeleteVpc(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletezone":
		return func() any {
			cmd := awsspec.NewDeleteZone(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachalarm":
		return func() any {
			cmd := awsspec.NewDetachAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachclassicloadbalancer":
		return func() any {
			cmd := awsspec.NewDetachClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachcontainertask":
		return func() any {
			cmd := awsspec.NewDetachContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachelasticip":
		return func() any {
			cmd := awsspec.NewDetachElasticip(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachinstance":
		return func() any {
			cmd := awsspec.NewDetachInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachinstanceprofile":
		return func() any {
			cmd := awsspec.NewDetachInstanceprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachinternetgateway":
		return func() any {
			cmd := awsspec.NewDetachInternetgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachmfadevice":
		return func() any {
			cmd := awsspec.NewDetachMfadevice(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachnetworkinterface":
		return func() any {
			cmd := awsspec.NewDetachNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachpolicy":
		return func() any {
			cmd := awsspec.NewDetachPolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachrole":
		return func() any {
			cmd := awsspec.NewDetachRole(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachroutetable":
		return func() any {
			cmd := awsspec.NewDetachRoutetable(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachsecuritygroup":
		return func() any {
			cmd := awsspec.NewDetachSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachuser":
		return func() any {
			cmd := awsspec.NewDetachUser(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachvolume":
		return func() any {
			cmd := awsspec.NewDetachVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "importimage":
		return func() any {
			cmd := awsspec.NewImportImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "restartdatabase":
		return func() any {
			cmd := awsspec.NewRestartDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "restartinstance":
		return func() any {
			cmd := awsspec.NewRestartInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "startalarm":
		return func() any {
			cmd := awsspec.NewStartAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "startcontainertask":
		return func() any {
			cmd := awsspec.NewStartContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "startdatabase":
		return func() any {
			cmd := awsspec.NewStartDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "startinstance":
		return func() any {
			cmd := awsspec.NewStartInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "stopalarm":
		return func() any {
			cmd := awsspec.NewStopAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "stopcontainertask":
		return func() any {
			cmd := awsspec.NewStopContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "stopdatabase":
		return func() any {
			cmd := awsspec.NewStopDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "stopinstance":
		return func() any {
			cmd := awsspec.NewStopInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatebucket":
		return func() any {
			cmd := awsspec.NewUpdateBucket(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updateclassicloadbalancer":
		return func() any {
			cmd := awsspec.NewUpdateClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatecontainertask":
		return func() any {
			cmd := awsspec.NewUpdateContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatedistribution":
		return func() any {
			cmd := awsspec.NewUpdateDistribution(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updateimage":
		return func() any {
			cmd := awsspec.NewUpdateImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updateinstance":
		return func() any {
			cmd := awsspec.NewUpdateInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updateloginprofile":
		return func() any {
			cmd := awsspec.NewUpdateLoginprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatepolicy":
		return func() any {
			cmd := awsspec.NewUpdatePolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updaterecord":
		return func() any {
			cmd := awsspec.NewUpdateRecord(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updates3object":
		return func() any {
			cmd := awsspec.NewUpdateS3object(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatescalinggroup":
		return func() any {
			cmd := awsspec.NewUpdateScalinggroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatesecuritygroup":
		return func() any {
			cmd := awsspec.NewUpdateSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatestack":
		return func() any {
			cmd := awsspec.NewUpdateStack(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatesubnet":
		return func() any {
			cmd := awsspec.NewUpdateSubnet(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatetargetgroup":
		return func() any {
			cmd := awsspec.NewUpdateTargetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	}
	return nil
}
