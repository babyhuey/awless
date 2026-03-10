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
	awsspec "github.com/wallix/awless/aws/spec"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/logger"
)

type AcceptanceFactory struct {
	Mock   interface{}
	Logger *logger.Logger
	Graph  cloud.GraphAPI
}

func NewAcceptanceFactory(mock interface{}, g cloud.GraphAPI, l ...*logger.Logger) *AcceptanceFactory {
	lg := logger.DiscardLogger
	if len(l) > 0 {
		lg = l[0]
	}
	return &AcceptanceFactory{Mock: mock, Graph: g, Logger: lg}
}

func (f *AcceptanceFactory) Build(key string) func() interface{} {
	switch key {
	case "attachalarm":
		return func() interface{} {
			cmd := awsspec.NewAttachAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachclassicloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewAttachClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachcontainertask":
		return func() interface{} {
			cmd := awsspec.NewAttachContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachelasticip":
		return func() interface{} {
			cmd := awsspec.NewAttachElasticip(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachinstance":
		return func() interface{} {
			cmd := awsspec.NewAttachInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachinstanceprofile":
		return func() interface{} {
			cmd := awsspec.NewAttachInstanceprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachinternetgateway":
		return func() interface{} {
			cmd := awsspec.NewAttachInternetgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachlistener":
		return func() interface{} {
			cmd := awsspec.NewAttachListener(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachmfadevice":
		return func() interface{} {
			cmd := awsspec.NewAttachMfadevice(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachnetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewAttachNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachpolicy":
		return func() interface{} {
			cmd := awsspec.NewAttachPolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachrole":
		return func() interface{} {
			cmd := awsspec.NewAttachRole(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachroutetable":
		return func() interface{} {
			cmd := awsspec.NewAttachRoutetable(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachsecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewAttachSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachuser":
		return func() interface{} {
			cmd := awsspec.NewAttachUser(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "attachvolume":
		return func() interface{} {
			cmd := awsspec.NewAttachVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "authenticateregistry":
		return func() interface{} {
			cmd := awsspec.NewAuthenticateRegistry(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkcertificate":
		return func() interface{} {
			cmd := awsspec.NewCheckCertificate(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkdatabase":
		return func() interface{} {
			cmd := awsspec.NewCheckDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkdistribution":
		return func() interface{} {
			cmd := awsspec.NewCheckDistribution(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkinstance":
		return func() interface{} {
			cmd := awsspec.NewCheckInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewCheckLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checknatgateway":
		return func() interface{} {
			cmd := awsspec.NewCheckNatgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checknetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewCheckNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkscalinggroup":
		return func() interface{} {
			cmd := awsspec.NewCheckScalinggroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checksecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewCheckSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "checkvolume":
		return func() interface{} {
			cmd := awsspec.NewCheckVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "copyimage":
		return func() interface{} {
			cmd := awsspec.NewCopyImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "copysnapshot":
		return func() interface{} {
			cmd := awsspec.NewCopySnapshot(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createaccesskey":
		return func() interface{} {
			cmd := awsspec.NewCreateAccesskey(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createalarm":
		return func() interface{} {
			cmd := awsspec.NewCreateAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createappscalingpolicy":
		return func() interface{} {
			cmd := awsspec.NewCreateAppscalingpolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createappscalingtarget":
		return func() interface{} {
			cmd := awsspec.NewCreateAppscalingtarget(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createbucket":
		return func() interface{} {
			cmd := awsspec.NewCreateBucket(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createcertificate":
		return func() interface{} {
			cmd := awsspec.NewCreateCertificate(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createclassicloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewCreateClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createcontainercluster":
		return func() interface{} {
			cmd := awsspec.NewCreateContainercluster(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createdatabase":
		return func() interface{} {
			cmd := awsspec.NewCreateDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createdbsubnetgroup":
		return func() interface{} {
			cmd := awsspec.NewCreateDbsubnetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createdistribution":
		return func() interface{} {
			cmd := awsspec.NewCreateDistribution(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createelasticip":
		return func() interface{} {
			cmd := awsspec.NewCreateElasticip(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createfunction":
		return func() interface{} {
			cmd := awsspec.NewCreateFunction(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "creategroup":
		return func() interface{} {
			cmd := awsspec.NewCreateGroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createimage":
		return func() interface{} {
			cmd := awsspec.NewCreateImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createinstance":
		return func() interface{} {
			cmd := awsspec.NewCreateInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createinstanceprofile":
		return func() interface{} {
			cmd := awsspec.NewCreateInstanceprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createinternetgateway":
		return func() interface{} {
			cmd := awsspec.NewCreateInternetgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createkeypair":
		return func() interface{} {
			cmd := awsspec.NewCreateKeypair(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createlaunchconfiguration":
		return func() interface{} {
			cmd := awsspec.NewCreateLaunchconfiguration(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createlistener":
		return func() interface{} {
			cmd := awsspec.NewCreateListener(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewCreateLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createloginprofile":
		return func() interface{} {
			cmd := awsspec.NewCreateLoginprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createmfadevice":
		return func() interface{} {
			cmd := awsspec.NewCreateMfadevice(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createnatgateway":
		return func() interface{} {
			cmd := awsspec.NewCreateNatgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createnetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewCreateNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createpolicy":
		return func() interface{} {
			cmd := awsspec.NewCreatePolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createqueue":
		return func() interface{} {
			cmd := awsspec.NewCreateQueue(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createrecord":
		return func() interface{} {
			cmd := awsspec.NewCreateRecord(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createrepository":
		return func() interface{} {
			cmd := awsspec.NewCreateRepository(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createrole":
		return func() interface{} {
			cmd := awsspec.NewCreateRole(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createroute":
		return func() interface{} {
			cmd := awsspec.NewCreateRoute(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createroutetable":
		return func() interface{} {
			cmd := awsspec.NewCreateRoutetable(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "creates3object":
		return func() interface{} {
			cmd := awsspec.NewCreateS3object(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createscalinggroup":
		return func() interface{} {
			cmd := awsspec.NewCreateScalinggroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createscalingpolicy":
		return func() interface{} {
			cmd := awsspec.NewCreateScalingpolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createsecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewCreateSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createsnapshot":
		return func() interface{} {
			cmd := awsspec.NewCreateSnapshot(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createstack":
		return func() interface{} {
			cmd := awsspec.NewCreateStack(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createsubnet":
		return func() interface{} {
			cmd := awsspec.NewCreateSubnet(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createsubscription":
		return func() interface{} {
			cmd := awsspec.NewCreateSubscription(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createtag":
		return func() interface{} {
			cmd := awsspec.NewCreateTag(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createtargetgroup":
		return func() interface{} {
			cmd := awsspec.NewCreateTargetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createtopic":
		return func() interface{} {
			cmd := awsspec.NewCreateTopic(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createuser":
		return func() interface{} {
			cmd := awsspec.NewCreateUser(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createvolume":
		return func() interface{} {
			cmd := awsspec.NewCreateVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createvpc":
		return func() interface{} {
			cmd := awsspec.NewCreateVpc(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "createzone":
		return func() interface{} {
			cmd := awsspec.NewCreateZone(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteaccesskey":
		return func() interface{} {
			cmd := awsspec.NewDeleteAccesskey(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletealarm":
		return func() interface{} {
			cmd := awsspec.NewDeleteAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteappscalingpolicy":
		return func() interface{} {
			cmd := awsspec.NewDeleteAppscalingpolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteappscalingtarget":
		return func() interface{} {
			cmd := awsspec.NewDeleteAppscalingtarget(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletebucket":
		return func() interface{} {
			cmd := awsspec.NewDeleteBucket(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletecertificate":
		return func() interface{} {
			cmd := awsspec.NewDeleteCertificate(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteclassicloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewDeleteClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletecontainercluster":
		return func() interface{} {
			cmd := awsspec.NewDeleteContainercluster(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletecontainertask":
		return func() interface{} {
			cmd := awsspec.NewDeleteContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletedatabase":
		return func() interface{} {
			cmd := awsspec.NewDeleteDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletedbsubnetgroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteDbsubnetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletedistribution":
		return func() interface{} {
			cmd := awsspec.NewDeleteDistribution(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteelasticip":
		return func() interface{} {
			cmd := awsspec.NewDeleteElasticip(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletefunction":
		return func() interface{} {
			cmd := awsspec.NewDeleteFunction(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletegroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteGroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteimage":
		return func() interface{} {
			cmd := awsspec.NewDeleteImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteinstance":
		return func() interface{} {
			cmd := awsspec.NewDeleteInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteinstanceprofile":
		return func() interface{} {
			cmd := awsspec.NewDeleteInstanceprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteinternetgateway":
		return func() interface{} {
			cmd := awsspec.NewDeleteInternetgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletekeypair":
		return func() interface{} {
			cmd := awsspec.NewDeleteKeypair(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletelaunchconfiguration":
		return func() interface{} {
			cmd := awsspec.NewDeleteLaunchconfiguration(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletelistener":
		return func() interface{} {
			cmd := awsspec.NewDeleteListener(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewDeleteLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteloginprofile":
		return func() interface{} {
			cmd := awsspec.NewDeleteLoginprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletemfadevice":
		return func() interface{} {
			cmd := awsspec.NewDeleteMfadevice(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletenatgateway":
		return func() interface{} {
			cmd := awsspec.NewDeleteNatgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletenetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewDeleteNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletepolicy":
		return func() interface{} {
			cmd := awsspec.NewDeletePolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletequeue":
		return func() interface{} {
			cmd := awsspec.NewDeleteQueue(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleterecord":
		return func() interface{} {
			cmd := awsspec.NewDeleteRecord(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleterepository":
		return func() interface{} {
			cmd := awsspec.NewDeleteRepository(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleterole":
		return func() interface{} {
			cmd := awsspec.NewDeleteRole(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteroute":
		return func() interface{} {
			cmd := awsspec.NewDeleteRoute(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteroutetable":
		return func() interface{} {
			cmd := awsspec.NewDeleteRoutetable(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletes3object":
		return func() interface{} {
			cmd := awsspec.NewDeleteS3object(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletescalinggroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteScalinggroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletescalingpolicy":
		return func() interface{} {
			cmd := awsspec.NewDeleteScalingpolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletesecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletesnapshot":
		return func() interface{} {
			cmd := awsspec.NewDeleteSnapshot(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletestack":
		return func() interface{} {
			cmd := awsspec.NewDeleteStack(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletesubnet":
		return func() interface{} {
			cmd := awsspec.NewDeleteSubnet(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletesubscription":
		return func() interface{} {
			cmd := awsspec.NewDeleteSubscription(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletetag":
		return func() interface{} {
			cmd := awsspec.NewDeleteTag(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletetargetgroup":
		return func() interface{} {
			cmd := awsspec.NewDeleteTargetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletetopic":
		return func() interface{} {
			cmd := awsspec.NewDeleteTopic(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deleteuser":
		return func() interface{} {
			cmd := awsspec.NewDeleteUser(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletevolume":
		return func() interface{} {
			cmd := awsspec.NewDeleteVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletevpc":
		return func() interface{} {
			cmd := awsspec.NewDeleteVpc(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "deletezone":
		return func() interface{} {
			cmd := awsspec.NewDeleteZone(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachalarm":
		return func() interface{} {
			cmd := awsspec.NewDetachAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachclassicloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewDetachClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachcontainertask":
		return func() interface{} {
			cmd := awsspec.NewDetachContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachelasticip":
		return func() interface{} {
			cmd := awsspec.NewDetachElasticip(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachinstance":
		return func() interface{} {
			cmd := awsspec.NewDetachInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachinstanceprofile":
		return func() interface{} {
			cmd := awsspec.NewDetachInstanceprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachinternetgateway":
		return func() interface{} {
			cmd := awsspec.NewDetachInternetgateway(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachmfadevice":
		return func() interface{} {
			cmd := awsspec.NewDetachMfadevice(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachnetworkinterface":
		return func() interface{} {
			cmd := awsspec.NewDetachNetworkinterface(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachpolicy":
		return func() interface{} {
			cmd := awsspec.NewDetachPolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachrole":
		return func() interface{} {
			cmd := awsspec.NewDetachRole(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachroutetable":
		return func() interface{} {
			cmd := awsspec.NewDetachRoutetable(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachsecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewDetachSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachuser":
		return func() interface{} {
			cmd := awsspec.NewDetachUser(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "detachvolume":
		return func() interface{} {
			cmd := awsspec.NewDetachVolume(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "importimage":
		return func() interface{} {
			cmd := awsspec.NewImportImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "restartdatabase":
		return func() interface{} {
			cmd := awsspec.NewRestartDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "restartinstance":
		return func() interface{} {
			cmd := awsspec.NewRestartInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "startalarm":
		return func() interface{} {
			cmd := awsspec.NewStartAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "startcontainertask":
		return func() interface{} {
			cmd := awsspec.NewStartContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "startdatabase":
		return func() interface{} {
			cmd := awsspec.NewStartDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "startinstance":
		return func() interface{} {
			cmd := awsspec.NewStartInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "stopalarm":
		return func() interface{} {
			cmd := awsspec.NewStopAlarm(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "stopcontainertask":
		return func() interface{} {
			cmd := awsspec.NewStopContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "stopdatabase":
		return func() interface{} {
			cmd := awsspec.NewStopDatabase(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "stopinstance":
		return func() interface{} {
			cmd := awsspec.NewStopInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatebucket":
		return func() interface{} {
			cmd := awsspec.NewUpdateBucket(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updateclassicloadbalancer":
		return func() interface{} {
			cmd := awsspec.NewUpdateClassicLoadbalancer(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatecontainertask":
		return func() interface{} {
			cmd := awsspec.NewUpdateContainertask(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatedistribution":
		return func() interface{} {
			cmd := awsspec.NewUpdateDistribution(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updateimage":
		return func() interface{} {
			cmd := awsspec.NewUpdateImage(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updateinstance":
		return func() interface{} {
			cmd := awsspec.NewUpdateInstance(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updateloginprofile":
		return func() interface{} {
			cmd := awsspec.NewUpdateLoginprofile(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatepolicy":
		return func() interface{} {
			cmd := awsspec.NewUpdatePolicy(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updaterecord":
		return func() interface{} {
			cmd := awsspec.NewUpdateRecord(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updates3object":
		return func() interface{} {
			cmd := awsspec.NewUpdateS3object(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatescalinggroup":
		return func() interface{} {
			cmd := awsspec.NewUpdateScalinggroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatesecuritygroup":
		return func() interface{} {
			cmd := awsspec.NewUpdateSecuritygroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatestack":
		return func() interface{} {
			cmd := awsspec.NewUpdateStack(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatesubnet":
		return func() interface{} {
			cmd := awsspec.NewUpdateSubnet(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	case "updatetargetgroup":
		return func() interface{} {
			cmd := awsspec.NewUpdateTargetgroup(aws.Config{}, f.Graph, f.Logger)
			// TODO: SDK v2 mocking needs rework - SetApi expects *service.Client
			_ = cmd
			return cmd
		}
	}
	return nil
}
