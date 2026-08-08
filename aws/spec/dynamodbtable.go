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
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/bootswithdefer/awless/cloud"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/template/env"
	"github.com/bootswithdefer/awless/template/params"
)

// CreateDynamodbtable creates a table with a partition key and an optional sort
// key.
//
// The API models the key schema and attribute definitions as parallel lists that
// have to agree, and billing mode determines whether throughput is required at
// all, so BeforeRun assembles the request rather than relying on field tags.
type CreateDynamodbtable struct {
	_             string `action:"create" entity:"dynamodbtable" awsAPI:"dynamodb"`
	logger        *logger.Logger
	graph         cloud.GraphAPI
	api           *dynamodb.Client
	Name          *string `templateName:"name"`
	PartitionKey  *string `templateName:"partition-key"`
	PartitionType *string `templateName:"partition-type"`
	SortKey       *string `templateName:"sort-key"`
	SortType      *string `templateName:"sort-type"`
	BillingMode   *string `templateName:"billing-mode"`
	ReadCapacity  *int64  `templateName:"read-capacity"`
	WriteCapacity *int64  `templateName:"write-capacity"`
}

func (cmd *CreateDynamodbtable) ParamsSpec() params.Spec {
	return params.NewSpec(
		params.AllOf(
			params.Key("name"), params.Key("partition-key"),
			params.Opt(
				params.Suggested("partition-type"),
				"sort-key", "sort-type", "billing-mode",
				"read-capacity", "write-capacity",
			),
		),
		params.Validators{
			"partition-type": isDynamoAttributeType,
			"sort-type":      isDynamoAttributeType,
			"billing-mode":   isDynamoBillingMode,
		})
}

func (cmd *CreateDynamodbtable) ManualRun(renv env.Running) (any, error) {
	partitionType := StringValue(cmd.PartitionType)
	if partitionType == "" {
		partitionType = "S"
	}

	in := &dynamodb.CreateTableInput{
		TableName: cmd.Name,
		KeySchema: []dynamodbtypes.KeySchemaElement{{
			AttributeName: cmd.PartitionKey,
			KeyType:       dynamodbtypes.KeyTypeHash,
		}},
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{{
			AttributeName: cmd.PartitionKey,
			AttributeType: dynamodbtypes.ScalarAttributeType(partitionType),
		}},
	}

	if sortKey := StringValue(cmd.SortKey); sortKey != "" {
		sortType := StringValue(cmd.SortType)
		if sortType == "" {
			sortType = "S"
		}
		in.KeySchema = append(in.KeySchema, dynamodbtypes.KeySchemaElement{
			AttributeName: cmd.SortKey,
			KeyType:       dynamodbtypes.KeyTypeRange,
		})
		in.AttributeDefinitions = append(in.AttributeDefinitions, dynamodbtypes.AttributeDefinition{
			AttributeName: cmd.SortKey,
			AttributeType: dynamodbtypes.ScalarAttributeType(sortType),
		})
	}

	// PAY_PER_REQUEST rejects a throughput block, and PROVISIONED requires one, so
	// default to on-demand unless capacity was asked for.
	mode := StringValue(cmd.BillingMode)
	if mode == "" {
		if cmd.ReadCapacity != nil || cmd.WriteCapacity != nil {
			mode = "PROVISIONED"
		} else {
			mode = "PAY_PER_REQUEST"
		}
	}
	in.BillingMode = dynamodbtypes.BillingMode(mode)

	if mode == "PROVISIONED" {
		read, write := int64(5), int64(5)
		if cmd.ReadCapacity != nil {
			read = *cmd.ReadCapacity
		}
		if cmd.WriteCapacity != nil {
			write = *cmd.WriteCapacity
		}
		in.ProvisionedThroughput = &dynamodbtypes.ProvisionedThroughput{
			ReadCapacityUnits:  awssdk.Int64(read),
			WriteCapacityUnits: awssdk.Int64(write),
		}
	} else if cmd.ReadCapacity != nil || cmd.WriteCapacity != nil {
		return nil, fmt.Errorf("read-capacity and write-capacity require billing-mode=PROVISIONED")
	}

	start := time.Now()
	output, err := cmd.api.CreateTable(renv.RequestContext(), in)
	cmd.logger.ExtraVerbosef("dynamodb.CreateTable call took %s", time.Since(start))
	if err != nil {
		return nil, fmt.Errorf("create dynamodbtable: %w", err)
	}
	return output, nil
}

func (cmd *CreateDynamodbtable) ExtractResult(i any) string {
	return awssdk.ToString(i.(*dynamodb.CreateTableOutput).TableDescription.TableName)
}

type DeleteDynamodbtable struct {
	_      string `action:"delete" entity:"dynamodbtable" awsAPI:"dynamodb" awsCall:"DeleteTable" awsInput:"dynamodb.DeleteTableInput" awsOutput:"dynamodb.DeleteTableOutput"`
	logger *logger.Logger
	graph  cloud.GraphAPI
	api    *dynamodb.Client
	Name   *string `awsName:"TableName" awsType:"awsstr" templateName:"name"`
}

func (cmd *DeleteDynamodbtable) ParamsSpec() params.Spec {
	return params.NewSpec(params.AllOf(params.Key("name")))
}

func isDynamoAttributeType(i any, _ map[string]any) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected a string, got %T", i)
	}
	switch s {
	case "S", "N", "B":
		return nil
	}
	return fmt.Errorf("invalid attribute type '%s', expected S (string), N (number) or B (binary)", s)
}

func isDynamoBillingMode(i any, _ map[string]any) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("expected a string, got %T", i)
	}
	switch s {
	case "PROVISIONED", "PAY_PER_REQUEST":
		return nil
	}
	return fmt.Errorf("invalid billing mode '%s', expected PROVISIONED or PAY_PER_REQUEST", s)
}
