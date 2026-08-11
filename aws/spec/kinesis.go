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
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	kinesistypes "github.com/aws/aws-sdk-go-v2/service/kinesis/types"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateStream struct {
	_      string `action:"create" entity:"stream" awsAPI:"kinesis" awsCall:"CreateStream" awsInput:"kinesis.CreateStreamInput" awsOutput:"kinesis.CreateStreamOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *kinesis.Client
	Name   *string `awsName:"StreamName" awsType:"awsstr" templateName:"name"`
	// Shards apply to a PROVISIONED stream only; an ON_DEMAND stream scales itself and
	// AWS rejects a shard count alongside it.
	Shards *int64  `awsName:"ShardCount" awsType:"awsint64" templateName:"shards"`
	Mode   *string `awsName:"StreamModeDetails.StreamMode" awsType:"awsstr" templateName:"mode"`
}

func (cmd *CreateStream) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(params.Suggested("mode"), "shards"),
	))
}

// CreateStream returns nothing, and the name is what the graph and delete key on.
func (cmd *CreateStream) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteStream struct {
	_      string `action:"delete" entity:"stream" awsAPI:"kinesis" awsCall:"DeleteStream" awsInput:"kinesis.DeleteStreamInput" awsOutput:"kinesis.DeleteStreamOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *kinesis.Client
	Name   *string `awsName:"StreamName" awsType:"awsstr" templateName:"name"`
	// Without this, deleting a stream that still has registered consumers fails.
	Force *bool `awsName:"EnforceConsumerDeletion" awsType:"awsbool" templateName:"force"`
}

func (cmd *DeleteStream) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("force"),
	))
}

// Resharding is the only in-place change worth exposing; retention and encryption each
// have their own API rather than a general modify call.
type UpdateStream struct {
	_      string `action:"update" entity:"stream" awsAPI:"kinesis" awsCall:"UpdateShardCount" awsInput:"kinesis.UpdateShardCountInput" awsOutput:"kinesis.UpdateShardCountOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *kinesis.Client
	Name   *string `awsName:"StreamName" awsType:"awsstr" templateName:"name"`
	Shards *int64  `awsName:"TargetShardCount" awsType:"awsint64" templateName:"shards"`
	// Required by the API, and UNIFORM_SCALING is the only value it accepts. Defaulted
	// in BeforeRun rather than asked for, since leaving it unset makes every call fail
	// and there is nothing to choose between.
	ScalingType *string `awsName:"ScalingType" awsType:"awsstr" templateName:"scaling-type"`
}

func (cmd *UpdateStream) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("shards"),
		params.Opt("scaling-type"),
	))
}

func (cmd *UpdateStream) BeforeRun(renv env.Running) error {
	if StringValue(cmd.ScalingType) == "" {
		cmd.ScalingType = String(string(kinesistypes.ScalingTypeUniformScaling))
	}
	return nil
}

func (cmd *UpdateStream) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}
