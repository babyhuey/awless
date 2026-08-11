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
	"github.com/aws/aws-sdk-go-v2/service/glue"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/params"
)

// The entities are gluedatabase and gluetable: "database" belongs to RDS, and a bare
// "table" would not say whether it means the data catalog or DynamoDB.

type CreateGluedatabase struct {
	_           string `action:"create" entity:"gluedatabase" awsAPI:"glue" awsCall:"CreateDatabase" awsInput:"glue.CreateDatabaseInput" awsOutput:"glue.CreateDatabaseOutput"`
	logger      *logger.Logger
	graph       cloud.GraphAPI
	api         *glue.Client
	Name        *string `awsName:"DatabaseInput.Name" awsType:"awsstr" templateName:"name"`
	Description *string `awsName:"DatabaseInput.Description" awsType:"awsstr" templateName:"description"`
	Location    *string `awsName:"DatabaseInput.LocationUri" awsType:"awsstr" templateName:"location"`
}

func (cmd *CreateGluedatabase) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt(params.Suggested("description"), "location"),
	))
}

func (cmd *CreateGluedatabase) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteGluedatabase struct {
	_      string `action:"delete" entity:"gluedatabase" awsAPI:"glue" awsCall:"DeleteDatabase" awsInput:"glue.DeleteDatabaseInput" awsOutput:"glue.DeleteDatabaseOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *glue.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

// Deleting a database deletes every table in it, which is worth knowing but is AWS's
// behavior rather than something awless can soften.
func (cmd *DeleteGluedatabase) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type CreateCrawler struct {
	_        string `action:"create" entity:"crawler" awsAPI:"glue" awsCall:"CreateCrawler" awsInput:"glue.CreateCrawlerInput" awsOutput:"glue.CreateCrawlerOutput"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *glue.Client
	Name     *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Role     *string `awsName:"Role" awsType:"awsstr" templateName:"role"`
	Database *string `awsName:"DatabaseName" awsType:"awsstr" templateName:"database"`
	// What to crawl: S3 paths, JDBC connections, DynamoDB tables and more, each with a
	// different shape, so it is a document.
	TargetsFile *string `awsName:"Targets" awsType:"awsfiletostruct" templateName:"targets-file"`
	Schedule    *string `awsName:"Schedule" awsType:"awsstr" templateName:"schedule"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
}

func (cmd *CreateCrawler) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("role"), params.Key("targets-file"),
		params.Opt(params.Suggested("database"), "schedule", "description"),
	))
}

func (cmd *CreateCrawler) ExtractResult(i any) string {
	return StringValue(cmd.Name)
}

type DeleteCrawler struct {
	_      string `action:"delete" entity:"crawler" awsAPI:"glue" awsCall:"DeleteCrawler" awsInput:"glue.DeleteCrawlerInput" awsOutput:"glue.DeleteCrawlerOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *glue.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteCrawler) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type StartCrawler struct {
	_      string `action:"start" entity:"crawler" awsAPI:"glue" awsCall:"StartCrawler" awsInput:"glue.StartCrawlerInput" awsOutput:"glue.StartCrawlerOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *glue.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *StartCrawler) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type StopCrawler struct {
	_      string `action:"stop" entity:"crawler" awsAPI:"glue" awsCall:"StopCrawler" awsInput:"glue.StopCrawlerInput" awsOutput:"glue.StopCrawlerOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *glue.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
}

func (cmd *StopCrawler) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type CreateJob struct {
	_      string `action:"create" entity:"job" awsAPI:"glue" awsCall:"CreateJob" awsInput:"glue.CreateJobInput" awsOutput:"glue.CreateJobOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *glue.Client
	Name   *string `awsName:"Name" awsType:"awsstr" templateName:"name"`
	Role   *string `awsName:"Role" awsType:"awsstr" templateName:"role"`
	// The command says which runtime and script to use. Its three fields are flat
	// enough to spell out rather than needing a document.
	CommandName *string `awsName:"Command.Name" awsType:"awsstr" templateName:"command"`
	ScriptPath  *string `awsName:"Command.ScriptLocation" awsType:"awsstr" templateName:"script"`
	PythonVer   *string `awsName:"Command.PythonVersion" awsType:"awsstr" templateName:"python-version"`
	GlueVersion *string `awsName:"GlueVersion" awsType:"awsstr" templateName:"glue-version"`
	WorkerType  *string `awsName:"WorkerType" awsType:"awsstr" templateName:"worker-type"`
	Workers     *int64  `awsName:"NumberOfWorkers" awsType:"awsint64" templateName:"workers"`
	Description *string `awsName:"Description" awsType:"awsstr" templateName:"description"`
	MaxRetries  *int64  `awsName:"MaxRetries" awsType:"awsint64" templateName:"max-retries"`
	TimeoutMins *int64  `awsName:"Timeout" awsType:"awsint64" templateName:"timeout"`
}

func (cmd *CreateJob) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"), params.Key("role"), params.Key("command"), params.Key("script"),
		params.Opt(
			params.Suggested("glue-version", "worker-type", "workers"),
			"python-version", "description", "max-retries", "timeout",
		),
	))
}

func (cmd *CreateJob) ExtractResult(i any) string {
	out, ok := i.(*glue.CreateJobOutput)
	if !ok || out.Name == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.Name)
}

type DeleteJob struct {
	_      string `action:"delete" entity:"job" awsAPI:"glue" awsCall:"DeleteJob" awsInput:"glue.DeleteJobInput" awsOutput:"glue.DeleteJobOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *glue.Client
	Name   *string `awsName:"JobName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteJob) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

type StartJob struct {
	_      string `action:"start" entity:"job" awsAPI:"glue" awsCall:"StartJobRun" awsInput:"glue.StartJobRunInput" awsOutput:"glue.StartJobRunOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *glue.Client
	Name   *string `awsName:"JobName" awsType:"awsstr" templateName:"name"`
	// Overrides for one run, so a job can be given more capacity without editing it.
	Workers    *int64  `awsName:"NumberOfWorkers" awsType:"awsint64" templateName:"workers"`
	WorkerType *string `awsName:"WorkerType" awsType:"awsstr" templateName:"worker-type"`
}

func (cmd *StartJob) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(
		params.Key("name"),
		params.Opt("workers", "worker-type"),
	))
}

// The run id is what identifies this execution afterwards.
func (cmd *StartJob) ExtractResult(i any) string {
	out, ok := i.(*glue.StartJobRunOutput)
	if !ok || out.JobRunId == nil {
		return StringValue(cmd.Name)
	}
	return awssdk.ToString(out.JobRunId)
}
