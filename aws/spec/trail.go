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

package awsspec

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

type CreateTrail struct {
	_               string `action:"create" entity:"trail" awsAPI:"cloudtrail" awsCall:"CreateTrail" awsInput:"cloudtrail.CreateTrailInput" awsOutput:"cloudtrail.CreateTrailOutput"`
	logger          *logger.Logger
	graph           cloud.GraphAPI
	api             *cloudtrail.Client
	Name            *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Bucket          *string `awsName:"S3BucketName" awsType:"awsstr" templateName:"bucket"`
	Prefix          *string `awsName:"S3KeyPrefix" awsType:"awsstr" templateName:"prefix"`
	Multiregion     *bool   `awsName:"IsMultiRegionTrail" awsType:"awsbool" templateName:"multiregion"`
	Global          *bool   `awsName:"IncludeGlobalServiceEvents" awsType:"awsbool" templateName:"global-events"`
	LogValidation   *bool   `awsName:"EnableLogFileValidation" awsType:"awsbool" templateName:"log-validation"`
	KmsKey          *string `awsName:"KmsKeyId" awsType:"awsstr" templateName:"kms-key"`
	SnsTopic        *string `awsName:"SnsTopicName" awsType:"awsstr" templateName:"sns-topic"`
	OrganizationAll *bool   `awsName:"IsOrganizationTrail" awsType:"awsbool" templateName:"organization"`
}

func (cmd *CreateTrail) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("bucket"),
		params.Opt(
			params.Suggested("multiregion", "log-validation"),
			"prefix", "global-events", "kms-key", "sns-topic", "organization",
		),
	))
}

func (cmd *CreateTrail) ExtractResult(i any) string {
	return awssdk.ToString(i.(*cloudtrail.CreateTrailOutput).Name)
}

type DeleteTrail struct {
	_      string `action:"delete" entity:"trail" awsAPI:"cloudtrail" awsCall:"DeleteTrail" awsInput:"cloudtrail.DeleteTrailInput" awsOutput:"cloudtrail.DeleteTrailOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cloudtrail.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteTrail) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

// StartTrail begins delivering events. A trail created by create trail is not
// logging until this runs, which is easy to miss.
type StartTrail struct {
	_      string `action:"start" entity:"trail" awsAPI:"cloudtrail" awsCall:"StartLogging" awsInput:"cloudtrail.StartLoggingInput" awsOutput:"cloudtrail.StartLoggingOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cloudtrail.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *StartTrail) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type StopTrail struct {
	_      string `action:"stop" entity:"trail" awsAPI:"cloudtrail" awsCall:"StopLogging" awsInput:"cloudtrail.StopLoggingInput" awsOutput:"cloudtrail.StopLoggingOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cloudtrail.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *StopTrail) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type CreateLoggroup struct {
	_      string `action:"create" entity:"loggroup" awsAPI:"cloudwatchlogs" awsCall:"CreateLogGroup" awsInput:"cloudwatchlogs.CreateLogGroupInput" awsOutput:"cloudwatchlogs.CreateLogGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cloudwatchlogs.Client
	Name   *string `awsName:"LogGroupName" awsType:"awsstr" templateName:"name"`
	KmsKey *string `awsName:"KmsKeyId" awsType:"awsstr" templateName:"kms-key"`
}

func (cmd *CreateLoggroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("kms-key"),
	))
}

// ExtractResult returns the name, since CreateLogGroup has an empty output and the
// name is what the graph keys on.
func (cmd *CreateLoggroup) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteLoggroup struct {
	_      string `action:"delete" entity:"loggroup" awsAPI:"cloudwatchlogs" awsCall:"DeleteLogGroup" awsInput:"cloudwatchlogs.DeleteLogGroupInput" awsOutput:"cloudwatchlogs.DeleteLogGroupOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *cloudwatchlogs.Client
	Name   *string `awsName:"LogGroupName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteLoggroup) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

// UpdateLoggroup sets the retention period. Zero means never expire, which the API
// models as DeleteRetentionPolicy rather than a retention of zero, so that case is
// handled separately in ManualRun.
type UpdateLoggroup struct {
	_         string `action:"update" entity:"loggroup" awsAPI:"cloudwatchlogs" awsCall:"PutRetentionPolicy" awsInput:"cloudwatchlogs.PutRetentionPolicyInput" awsOutput:"cloudwatchlogs.PutRetentionPolicyOutput"`
	logger    *logger.Logger
	graph     cloud.GraphAPI
	api       *cloudwatchlogs.Client
	Name      *string `awsName:"LogGroupName" awsType:"awsstr" templateName:"name"`
	Retention *int64  `awsName:"RetentionInDays" awsType:"awsint64" templateName:"retention"`
}

func (cmd *UpdateLoggroup) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(params.Key("name"), params.Key("retention")),
		params.Validators{"retention": isLogRetentionDays})
}

// isLogRetentionDays rejects a value CloudWatch Logs would not accept. The API
// only allows a fixed set, and its error for anything else does not list them.
func isLogRetentionDays(i any, _ map[string]any) error {
	var days int64
	switch v := i.(type) {
	case int:
		days = int64(v)
	case int64:
		days = v
	case string:
		parsed, err := castInt64(v)
		if err != nil {
			return fmt.Errorf("invalid retention '%s', expected a number of days", v)
		}
		days = parsed
	default:
		return fmt.Errorf("invalid retention type %T, expected a number of days", i)
	}

	for _, allowed := range []int64{
		1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096,
		1827, 2192, 2557, 2922, 3288, 3653,
	} {
		if days == allowed {
			return nil
		}
	}
	return fmt.Errorf("invalid retention %d, CloudWatch Logs accepts only 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1096, 1827, 2192, 2557, 2922, 3288 or 3653 days", days)
}
