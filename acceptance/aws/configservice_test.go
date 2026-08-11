package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
)

// PutConfigRule nests everything under a ConfigRule struct, so these assert that the
// dotted awsName paths built the intermediate structs rather than dropping the values.
// A path that does not resolve is discarded silently.
func TestCreateConfigrule(t *testing.T) {
	mock := NewMock().On("PutConfigRule", &configservice.PutConfigRuleOutput{})

	Template("create configrule name=s3-versioning source=S3_BUCKET_VERSIONING_ENABLED " +
		"description=Versioning frequency=TwentyFour_Hours resource-types=AWS::S3::Bucket").
		Mock(mock).
		ExpectCalls("PutConfigRule").
		ExpectCommandResult("s3-versioning").
		ExpectRevert("delete configrule name=s3-versioning").
		Run(t)

	in := mock.InputFor("PutConfigRule").(*configservice.PutConfigRuleInput)
	if in.ConfigRule == nil {
		t.Fatal("the nested ConfigRule was never built")
	}
	if got := awssdk.ToString(in.ConfigRule.ConfigRuleName); got != "s3-versioning" {
		t.Errorf("ConfigRuleName: got %q, want s3-versioning", got)
	}
	if in.ConfigRule.Source == nil {
		t.Fatal("the nested Source was never built")
	}
	if got := awssdk.ToString(in.ConfigRule.Source.SourceIdentifier); got != "S3_BUCKET_VERSIONING_ENABLED" {
		t.Errorf("Source.SourceIdentifier: got %q", got)
	}
	if got := awssdk.ToString(in.ConfigRule.Description); got != "Versioning" {
		t.Errorf("Description: got %q, want Versioning", got)
	}
	if got := string(in.ConfigRule.MaximumExecutionFrequency); got != "TwentyFour_Hours" {
		t.Errorf("MaximumExecutionFrequency: got %q", got)
	}
	// Two levels of nesting plus a slice.
	if in.ConfigRule.Scope == nil {
		t.Fatal("the nested Scope was never built")
	}
	if len(in.ConfigRule.Scope.ComplianceResourceTypes) != 1 ||
		in.ConfigRule.Scope.ComplianceResourceTypes[0] != "AWS::S3::Bucket" {
		t.Errorf("Scope.ComplianceResourceTypes: got %v", in.ConfigRule.Scope.ComplianceResourceTypes)
	}
}

func TestCreateConfigruleWithExplicitOwner(t *testing.T) {
	mock := NewMock().On("PutConfigRule", &configservice.PutConfigRuleOutput{})

	Template("create configrule name=custom source=arn:aws:lambda:us-west-2:1:function:check owner=CUSTOM_LAMBDA").
		Mock(mock).
		ExpectCalls("PutConfigRule").
		Run(t)

	in := mock.InputFor("PutConfigRule").(*configservice.PutConfigRuleInput)
	if got := string(in.ConfigRule.Source.Owner); got != "CUSTOM_LAMBDA" {
		t.Errorf("Source.Owner: got %q, want CUSTOM_LAMBDA", got)
	}
}

func TestUpdateConfigrule(t *testing.T) {
	mock := NewMock().On("PutConfigRule", &configservice.PutConfigRuleOutput{})

	Template("update configrule name=s3-versioning source=S3_BUCKET_VERSIONING_ENABLED input-parameters='{\"minimumCount\":1}'").
		Mock(mock).
		ExpectCalls("PutConfigRule").
		Run(t)

	in := mock.InputFor("PutConfigRule").(*configservice.PutConfigRuleInput)
	if got := awssdk.ToString(in.ConfigRule.InputParameters); got != `{"minimumCount":1}` {
		t.Errorf("InputParameters: got %q, want the JSON parameters", got)
	}
}

func TestDeleteConfigrule(t *testing.T) {
	mock := NewMock().On("DeleteConfigRule", &configservice.DeleteConfigRuleOutput{})

	Template("delete configrule name=s3-versioning").
		Mock(mock).
		ExpectCalls("DeleteConfigRule").
		Run(t)

	in := mock.InputFor("DeleteConfigRule").(*configservice.DeleteConfigRuleInput)
	if got := awssdk.ToString(in.ConfigRuleName); got != "s3-versioning" {
		t.Errorf("ConfigRuleName: got %q, want s3-versioning", got)
	}
}
