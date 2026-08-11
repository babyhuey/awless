package awsat

import (
	"os"
	"path/filepath"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	cptypes "github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafv2types "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func docFile(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("writing %s: %s", name, err)
	}
	return path
}

// A pipeline is created from a JSON declaration rather than flags. This asserts the whole
// document reached the request, nested stages and all, since a document that decoded to an
// empty struct would look like success.
func TestCreatePipelineFromADocument(t *testing.T) {
	def := docFile(t, "pipeline.json", `{
	  "name": "build-and-deploy",
	  "roleArn": "arn:aws:iam::123456789012:role/CodePipeline",
	  "artifactStore": {"type": "S3", "location": "my-artifacts"},
	  "stages": [
	    {"name": "Source", "actions": [
	      {"name": "Source",
	       "actionTypeId": {"category":"Source","owner":"AWS","provider":"S3","version":"1"},
	       "configuration": {"S3Bucket":"src","S3ObjectKey":"app.zip"}}
	    ]},
	    {"name": "Build", "actions": [
	      {"name": "Build",
	       "actionTypeId": {"category":"Build","owner":"AWS","provider":"CodeBuild","version":"1"},
	       "configuration": {"ProjectName":"api-build"}}
	    ]}
	  ]
	}`)

	mock := NewMock().On("CreatePipeline", &codepipeline.CreatePipelineOutput{
		Pipeline: &cptypes.PipelineDeclaration{Name: awssdk.String("build-and-deploy")},
	})

	Template("create pipeline definition-file=" + def).
		Mock(mock).
		ExpectCalls("CreatePipeline").
		ExpectCommandResult("build-and-deploy").
		Run(t)

	in := mock.InputFor("CreatePipeline").(*codepipeline.CreatePipelineInput)
	if in.Pipeline == nil {
		t.Fatal("the declaration was not decoded")
	}
	if got := awssdk.ToString(in.Pipeline.Name); got != "build-and-deploy" {
		t.Errorf("Name: got %q, want build-and-deploy", got)
	}
	if got := awssdk.ToString(in.Pipeline.RoleArn); got != "arn:aws:iam::123456789012:role/CodePipeline" {
		t.Errorf("RoleArn: got %q", got)
	}
	if len(in.Pipeline.Stages) != 2 {
		t.Fatalf("Stages: got %d, want 2", len(in.Pipeline.Stages))
	}
	if got := awssdk.ToString(in.Pipeline.Stages[1].Name); got != "Build" {
		t.Errorf("Stages[1].Name: got %q, want Build", got)
	}
	if got := in.Pipeline.Stages[1].Actions[0].Configuration["ProjectName"]; got != "api-build" {
		t.Errorf("the build action's configuration did not survive: got %q", got)
	}
}

// A typo in the document must fail rather than produce a pipeline missing a stage.
func TestCreatePipelineRejectsABadDocument(t *testing.T) {
	def := docFile(t, "pipeline.json", `{"name": "p", "stagez": []}`)

	err := Template("create pipeline definition-file=" + def).
		Mock(NewMock()).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected a misspelled key to be rejected")
	}
}

func TestCreateWebacl(t *testing.T) {
	action := docFile(t, "action.json", `{"allow": {}}`)
	visibility := docFile(t, "visibility.json",
		`{"sampledRequestsEnabled": true, "cloudWatchMetricsEnabled": true, "metricName": "publicApi"}`)
	rules := docFile(t, "rules.json",
		`[{"name":"RateLimit","priority":1,
		   "statement":{"rateBasedStatement":{"limit":2000,"aggregateKeyType":"IP"}},
		   "action":{"block":{}},
		   "visibilityConfig":{"sampledRequestsEnabled":true,"cloudWatchMetricsEnabled":true,"metricName":"rateLimit"}}]`)

	mock := NewMock().On("CreateWebACL", &wafv2.CreateWebACLOutput{
		Summary: &wafv2types.WebACLSummary{Name: awssdk.String("public-api")},
	})

	Template("create webacl name=public-api default-action-file=" + action +
		" visibility-file=" + visibility + " rules-file=" + rules).
		Mock(mock).
		ExpectCalls("CreateWebACL").
		ExpectCommandResult("public-api").
		ExpectRevert("delete webacl name=public-api").
		Run(t)

	in := mock.InputFor("CreateWebACL").(*wafv2.CreateWebACLInput)
	if got := awssdk.ToString(in.Name); got != "public-api" {
		t.Errorf("Name: got %q, want public-api", got)
	}
	// Scope defaults rather than being demanded.
	if in.Scope != wafv2types.ScopeRegional {
		t.Errorf("Scope: got %q, want REGIONAL", in.Scope)
	}
	if in.DefaultAction == nil || in.DefaultAction.Allow == nil {
		t.Error("DefaultAction.Allow was not decoded — an empty JSON object must still build the struct")
	}
	if in.VisibilityConfig == nil {
		t.Fatal("VisibilityConfig was not decoded")
	}
	if got := awssdk.ToString(in.VisibilityConfig.MetricName); got != "publicApi" {
		t.Errorf("VisibilityConfig.MetricName: got %q", got)
	}
	// A top-level JSON array decoding into a slice field.
	if len(in.Rules) != 1 {
		t.Fatalf("Rules: got %d, want 1", len(in.Rules))
	}
	if got := awssdk.ToString(in.Rules[0].Name); got != "RateLimit" {
		t.Errorf("Rules[0].Name: got %q, want RateLimit", got)
	}
	if in.Rules[0].Statement == nil || in.Rules[0].Statement.RateBasedStatement == nil {
		t.Error("the rule's rate-based statement was not decoded")
	}
}

func TestDeleteWebaclLooksUpTheLockToken(t *testing.T) {
	mock := NewMock().
		On("ListWebACLs", &wafv2.ListWebACLsOutput{
			WebACLs: []wafv2types.WebACLSummary{
				{Name: awssdk.String("other"), Id: awssdk.String("id-other"), LockToken: awssdk.String("tok-other")},
				{Name: awssdk.String("public-api"), Id: awssdk.String("id-1"), LockToken: awssdk.String("tok-1")},
			},
		}).
		On("DeleteWebACL", &wafv2.DeleteWebACLOutput{})

	Template("delete webacl name=public-api").
		Mock(mock).
		ExpectCalls("ListWebACLs", "DeleteWebACL").
		Run(t)

	in := mock.InputFor("DeleteWebACL").(*wafv2.DeleteWebACLInput)
	if got := awssdk.ToString(in.Id); got != "id-1" {
		t.Errorf("Id: got %q, want id-1 — the wrong ACL was matched", got)
	}
	if got := awssdk.ToString(in.LockToken); got != "tok-1" {
		t.Errorf("LockToken: got %q, want tok-1", got)
	}
}

func TestCreateRulegroup(t *testing.T) {
	visibility := docFile(t, "visibility.json",
		`{"sampledRequestsEnabled": true, "cloudWatchMetricsEnabled": true, "metricName": "rateLimits"}`)

	mock := NewMock().On("CreateRuleGroup", &wafv2.CreateRuleGroupOutput{
		Summary: &wafv2types.RuleGroupSummary{Name: awssdk.String("rate-limits")},
	})

	Template("create rulegroup name=rate-limits capacity=100 visibility-file=" + visibility).
		Mock(mock).
		ExpectCalls("CreateRuleGroup").
		ExpectCommandResult("rate-limits").
		ExpectRevert("delete rulegroup name=rate-limits").
		Run(t)

	in := mock.InputFor("CreateRuleGroup").(*wafv2.CreateRuleGroupInput)
	if got := awssdk.ToInt64(in.Capacity); got != 100 {
		t.Errorf("Capacity: got %d, want 100", got)
	}
	if in.VisibilityConfig == nil {
		t.Fatal("VisibilityConfig was not decoded")
	}
}

func TestDeleteRulegroup(t *testing.T) {
	mock := NewMock().
		On("ListRuleGroups", &wafv2.ListRuleGroupsOutput{
			RuleGroups: []wafv2types.RuleGroupSummary{
				{Name: awssdk.String("rate-limits"), Id: awssdk.String("id-1"), LockToken: awssdk.String("tok-1")},
			},
		}).
		On("DeleteRuleGroup", &wafv2.DeleteRuleGroupOutput{})

	Template("delete rulegroup name=rate-limits").
		Mock(mock).
		ExpectCalls("ListRuleGroups", "DeleteRuleGroup").
		Run(t)

	in := mock.InputFor("DeleteRuleGroup").(*wafv2.DeleteRuleGroupInput)
	if got := awssdk.ToString(in.LockToken); got != "tok-1" {
		t.Errorf("LockToken: got %q, want tok-1", got)
	}
}
