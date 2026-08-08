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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"

	"github.com/bootswithdefer/awless/logger"
)

type CreateTag struct {
	_        string `action:"create" entity:"tag" awsAPI:"ec2" awsDryRun:"manual"` //  awsCall:"CreateTags" awsInput:"ec2.CreateTagsInput" awsOutput:"ec2.CreateTagsOutput"
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *ec2.Client
	Resource *string `awsName:"Resources" awsType:"awsstringslice" templateName:"resource"`
	Key      *string `templateName:"key"`
	Value    *string `templateName:"value"`
}

func (cmd *CreateTag) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("key"), params.Key("resource"), params.Key("value")))
}

func (cmd *CreateTag) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("dry run: cannot set params on command struct: %w", err)
	}

	input := &ec2.CreateTagsInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("dry run: cannot inject in ec2.CreateTagsInput: %w", err)
	}
	input.Tags = []ec2types.Tag{{Key: cmd.Key, Value: cmd.Value}}

	start := time.Now()
	_, err := cmd.api.CreateTags(renv.RequestContext(), input)
	var awsErr smithy.APIError
	if errors.As(err, &awsErr) {
		switch code := awsErr.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			cmd.logger.ExtraVerbosef("dry run: ec2.CreateTags call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create tag ok")
			return fakeDryRunID("tag"), nil
		}
	}

	return nil, fmt.Errorf("dry run: %w", err)
}

func (cmd *CreateTag) ManualRun(renv env.Running) (any, error) {
	input := &ec2.CreateTagsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.CreateTagsInput: %w", err)
	}
	input.Tags = []ec2types.Tag{{Key: cmd.Key, Value: cmd.Value}}

	start := time.Now()
	var err error
	for attempt := 0; attempt < 5; attempt++ {
		_, err = cmd.api.CreateTags(renv.RequestContext(), input)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}
	if err != nil {
		return nil, err
	}
	cmd.logger.ExtraVerbosef("ec2.CreateTags call took %s", time.Since(start))
	return nil, nil
}

type DeleteTag struct {
	_        string `action:"delete" entity:"tag" awsAPI:"ec2" awsDryRun:"manual"`
	logger   *logger.Logger
	graph    cloud.GraphAPI
	api      *ec2.Client
	Resource *string `awsName:"Resources" awsType:"awsstringslice" templateName:"resource"`
	Key      *string `templateName:"key"`
	Value    *string `templateName:"value"`
}

func (cmd *DeleteTag) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("key"), params.Key("resource"),
		params.Opt("value"),
	))
}

func (cmd *DeleteTag) dryRun(renv env.Running, params map[string]any) (any, error) {
	if err := cmd.inject(params); err != nil {
		return nil, fmt.Errorf("cannot set params on command struct: %w", err)
	}

	input := &ec2.DeleteTagsInput{}
	input.DryRun = aws.Bool(true)
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteTagsInput: %w", err)
	}
	input.Tags = []ec2types.Tag{{Key: cmd.Key, Value: cmd.Value}}

	start := time.Now()
	_, err := cmd.api.DeleteTags(renv.RequestContext(), input)
	var awsErr smithy.APIError
	if errors.As(err, &awsErr) {
		switch code := awsErr.ErrorCode(); {
		case code == dryRunOperation, strings.HasSuffix(code, notFound):
			cmd.logger.ExtraVerbosef("dry run: ec2.DeleteTags call took %s", time.Since(start))
			cmd.logger.Verbose("dry run: create tag ok")
			return fakeDryRunID("tag"), nil
		}
	}

	return nil, err
}

func (cmd *DeleteTag) ManualRun(renv env.Running) (any, error) {
	input := &ec2.DeleteTagsInput{}
	if err := structInjector(cmd, input, renv); err != nil {
		return nil, fmt.Errorf("cannot inject in ec2.DeleteTagsInput: %w", err)
	}
	input.Tags = []ec2types.Tag{{Key: cmd.Key, Value: cmd.Value}}

	start := time.Now()
	_, err := cmd.api.DeleteTags(renv.RequestContext(), input)
	cmd.logger.ExtraVerbosef("ec2.DeleteTags call took %s", time.Since(start))
	return nil, err
}

// createNameTag tags a freshly created resource with its Name.
//
// name is nil whenever the user did not pass `name=`, which is legal for every
// command that calls this from AfterRun. Attempting the tag then panicked while
// injecting a nil Value into CreateTag, so return early instead.
func createNameTag(resource, name *string, renv env.Running) error {
	if name == nil || *name == "" {
		return nil
	}
	createTag := CommandFactory.Build("createtag")().(*CreateTag)
	entries := map[string]any{
		"key":      "Name",
		"value":    name,
		"resource": resource,
	}
	if err := params.Validate(createTag.ParamsSpec().Validators(), entries); err != nil {
		return err
	}
	_, err := createTag.Run(renv, entries)
	return err
}
