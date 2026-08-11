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
	"github.com/aws/aws-sdk-go-v2/service/kafka"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// The entity is kafkacluster: "cluster" would be ambiguous next to ECS, EKS and Redshift,
// all of which already have their own.
type CreateKafkacluster struct {
	_       string `action:"create" entity:"kafkacluster" awsAPI:"kafka" awsCall:"CreateCluster" awsInput:"kafka.CreateClusterInput" awsOutput:"kafka.CreateClusterOutput"`
	logger  *logger.Logger
	graph   cloud.GraphAPI
	api     *kafka.Client
	Name    *string   `awsName:"ClusterName" awsType:"awsstr" templateName:"name"`
	Version *string   `awsName:"KafkaVersion" awsType:"awsstr" templateName:"version"`
	Brokers *int64    `awsName:"NumberOfBrokerNodes" awsType:"awsint64" templateName:"brokers"`
	Subnets []*string `awsName:"BrokerNodeGroupInfo.ClientSubnets" awsType:"awsstringslice" templateName:"subnets"`
	Type    *string   `awsName:"BrokerNodeGroupInfo.InstanceType" awsType:"awsstr" templateName:"type"`
	// Security groups and storage sit inside the broker node group alongside the subnets.
	Securitygroups []*string `awsName:"BrokerNodeGroupInfo.SecurityGroups" awsType:"awsstringslice" templateName:"securitygroups"`
	StorageGB      *int64    `awsName:"BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.VolumeSize" awsType:"awsint64" templateName:"storage"`
	// Encryption and client authentication are nested structures with several mutually
	// relevant options, so they come from files.
	EncryptionFile *string `awsName:"EncryptionInfo" awsType:"awsfiletostruct" templateName:"encryption-file"`
	AuthFile       *string `awsName:"ClientAuthentication" awsType:"awsfiletostruct" templateName:"auth-file"`
}

// The broker count must be a multiple of the number of subnets, which AWS enforces; there
// is nothing useful awless can default it to.
func (cmd *CreateKafkacluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("version"), params.Key("brokers"),
		params.Key("subnets"), params.Key("type"),
		params.Opt(
			params.Suggested("storage", "securitygroups"),
			"encryption-file", "auth-file",
		),
	))
}

func (cmd *CreateKafkacluster) ExtractResult(i any) string {
	out, ok := i.(*kafka.CreateClusterOutput)
	if !ok || out.ClusterName == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.ClusterName)
}

// DeleteCluster takes the ARN rather than the name, and the ARN is what the graph records.
type DeleteKafkacluster struct {
	_      string `action:"delete" entity:"kafkacluster" awsAPI:"kafka" awsCall:"DeleteCluster" awsInput:"kafka.DeleteClusterInput" awsOutput:"kafka.DeleteClusterOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *kafka.Client
	Arn    *string `awsName:"ClusterArn" awsType:"awsstr" templateName:"arn"`
}

func (cmd *DeleteKafkacluster) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("arn")))
}
