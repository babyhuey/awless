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

package awsspec

import (
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/networkmanager"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// Global Network

type CreateGlobalnetwork struct {
	_           string `action:"create" entity:"globalnetwork" awsAPI:"networkmanager" awsCall:"CreateGlobalNetwork" awsInput:"networkmanager.CreateGlobalNetworkInput" awsOutput:"networkmanager.CreateGlobalNetworkOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *networkmanager.Client
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateGlobalnetwork) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Opt("description")))
}

func (cmd *CreateGlobalnetwork) ExtractResult(i any) string {
	out, ok := i.(*networkmanager.CreateGlobalNetworkOutput)
	if !ok || out.GlobalNetwork == nil {
		return ""
	}
	return awssdk.ToString(out.GlobalNetwork.GlobalNetworkId)
}

type DeleteGlobalnetwork struct {
	_      string `action:"delete" entity:"globalnetwork" awsAPI:"networkmanager" awsCall:"DeleteGlobalNetwork" awsInput:"networkmanager.DeleteGlobalNetworkInput" awsOutput:"networkmanager.DeleteGlobalNetworkOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *networkmanager.Client
	ID     *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteGlobalnetwork) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

// Core Network

type CreateCorenetwork struct {
	_             string `action:"create" entity:"corenetwork" awsAPI:"networkmanager" awsCall:"CreateCoreNetwork" awsInput:"networkmanager.CreateCoreNetworkInput" awsOutput:"networkmanager.CreateCoreNetworkOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *networkmanager.Client
	GlobalNetwork *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
	Description   *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateCorenetwork) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("global-network"), params.Opt("description")))
}

func (cmd *CreateCorenetwork) ExtractResult(i any) string {
	out, ok := i.(*networkmanager.CreateCoreNetworkOutput)
	if !ok || out.CoreNetwork == nil {
		return ""
	}
	return awssdk.ToString(out.CoreNetwork.CoreNetworkId)
}

type DeleteCorenetwork struct {
	_      string `action:"delete" entity:"corenetwork" awsAPI:"networkmanager" awsCall:"DeleteCoreNetwork" awsInput:"networkmanager.DeleteCoreNetworkInput" awsOutput:"networkmanager.DeleteCoreNetworkOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *networkmanager.Client
	ID     *string `awsName:"CoreNetworkId" awsType:"awsstr" templateName:"id"`
}

func (cmd *DeleteCorenetwork) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id")))
}

// Site

type CreateNetworkmanagersite struct {
	_             string `action:"create" entity:"networkmanagersite" awsAPI:"networkmanager" awsCall:"CreateSite" awsInput:"networkmanager.CreateSiteInput" awsOutput:"networkmanager.CreateSiteOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *networkmanager.Client
	GlobalNetwork *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
	Description   *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateNetworkmanagersite) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("global-network"), params.Opt("description")))
}

func (cmd *CreateNetworkmanagersite) ExtractResult(i any) string {
	out, ok := i.(*networkmanager.CreateSiteOutput)
	if !ok || out.Site == nil {
		return ""
	}
	return awssdk.ToString(out.Site.SiteId)
}

type DeleteNetworkmanagersite struct {
	_             string `action:"delete" entity:"networkmanagersite" awsAPI:"networkmanager" awsCall:"DeleteSite" awsInput:"networkmanager.DeleteSiteInput" awsOutput:"networkmanager.DeleteSiteOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *networkmanager.Client
	ID            *string `awsName:"SiteId" awsType:"awsstr" templateName:"id"`
	GlobalNetwork *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
}

func (cmd *DeleteNetworkmanagersite) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("global-network")))
}

// Device

type CreateNetworkmanagerdevice struct {
	_             string `action:"create" entity:"networkmanagerdevice" awsAPI:"networkmanager" awsCall:"CreateDevice" awsInput:"networkmanager.CreateDeviceInput" awsOutput:"networkmanager.CreateDeviceOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *networkmanager.Client
	GlobalNetwork *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
	Description   *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Model         *string `awsName:"Model" awsType:"awsstr" templateName:"model"`
	SerialNumber  *string `awsName:"SerialNumber" awsType:"awsstr" templateName:"serial-number"`
	Type          *string `awsName:"Type" awsType:"awsstr" templateName:"type"`
	Vendor        *string `awsName:"Vendor" awsType:"awsstr" templateName:"vendor"`
	SiteID        *string `awsName:"SiteId" awsType:"awsstr" templateName:"site"`
}

func (cmd *CreateNetworkmanagerdevice) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("global-network"),
		params.Opt("description", "model", "serial-number", "type", "vendor", "site"),
	))
}

func (cmd *CreateNetworkmanagerdevice) ExtractResult(i any) string {
	out, ok := i.(*networkmanager.CreateDeviceOutput)
	if !ok || out.Device == nil {
		return ""
	}
	return awssdk.ToString(out.Device.DeviceId)
}

type DeleteNetworkmanagerdevice struct {
	_             string `action:"delete" entity:"networkmanagerdevice" awsAPI:"networkmanager" awsCall:"DeleteDevice" awsInput:"networkmanager.DeleteDeviceInput" awsOutput:"networkmanager.DeleteDeviceOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *networkmanager.Client
	ID            *string `awsName:"DeviceId" awsType:"awsstr" templateName:"id"`
	GlobalNetwork *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
}

func (cmd *DeleteNetworkmanagerdevice) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("global-network")))
}

// Link

type CreateNetworkmanagerlink struct {
	_             string `action:"create" entity:"networkmanagerlink" awsAPI:"networkmanager" awsCall:"CreateLink" awsInput:"networkmanager.CreateLinkInput" awsOutput:"networkmanager.CreateLinkOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *networkmanager.Client
	GlobalNetwork *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
	SiteID        *string `awsName:"SiteId" awsType:"awsstr" templateName:"site"`
	Description   *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	Type          *string `awsName:"Type" awsType:"awsstr" templateName:"type"`
	Provider      *string `awsName:"Provider" awsType:"awsstr" templateName:"provider"`
}

func (cmd *CreateNetworkmanagerlink) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("global-network"), params.Key("site"),
		params.Opt("description", "type", "provider"),
	))
}

func (cmd *CreateNetworkmanagerlink) ExtractResult(i any) string {
	out, ok := i.(*networkmanager.CreateLinkOutput)
	if !ok || out.Link == nil {
		return ""
	}
	return awssdk.ToString(out.Link.LinkId)
}

type DeleteNetworkmanagerlink struct {
	_             string `action:"delete" entity:"networkmanagerlink" awsAPI:"networkmanager" awsCall:"DeleteLink" awsInput:"networkmanager.DeleteLinkInput" awsOutput:"networkmanager.DeleteLinkOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *networkmanager.Client
	ID            *string `awsName:"LinkId" awsType:"awsstr" templateName:"id"`
	GlobalNetwork *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
}

func (cmd *DeleteNetworkmanagerlink) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("global-network")))
}

// Connection

type CreateNetworkmanagerconnection struct {
	_               string `action:"create" entity:"networkmanagerconnection" awsAPI:"networkmanager" awsCall:"CreateConnection" awsInput:"networkmanager.CreateConnectionInput" awsOutput:"networkmanager.CreateConnectionOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             *networkmanager.Client
	GlobalNetwork   *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
	Device          *string `awsName:"DeviceId" awsType:"awsstr" templateName:"device"`
	ConnectedDevice *string `awsName:"ConnectedDeviceId" awsType:"awsstr" templateName:"connected-device"`
	Link            *string `awsName:"LinkId" awsType:"awsstr" templateName:"link"`
	ConnectedLink   *string `awsName:"ConnectedLinkId" awsType:"awsstr" templateName:"connected-link"`
	Description     *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateNetworkmanagerconnection) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("global-network"), params.Key("device"), params.Key("connected-device"),
		params.Opt("link", "connected-link", "description"),
	))
}

func (cmd *CreateNetworkmanagerconnection) ExtractResult(i any) string {
	out, ok := i.(*networkmanager.CreateConnectionOutput)
	if !ok || out.Connection == nil {
		return ""
	}
	return awssdk.ToString(out.Connection.ConnectionId)
}

type DeleteNetworkmanagerconnection struct {
	_             string `action:"delete" entity:"networkmanagerconnection" awsAPI:"networkmanager" awsCall:"DeleteConnection" awsInput:"networkmanager.DeleteConnectionInput" awsOutput:"networkmanager.DeleteConnectionOutput"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *networkmanager.Client
	ID            *string `awsName:"ConnectionId" awsType:"awsstr" templateName:"id"`
	GlobalNetwork *string `awsName:"GlobalNetworkId" awsType:"awsstr" templateName:"global-network"`
}

func (cmd *DeleteNetworkmanagerconnection) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("id"), params.Key("global-network")))
}
