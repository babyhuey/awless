package awsat

import (
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

// IAM commands. Each asserts that the template params reach the AWS input under the
// right field, which is the mapping the awsName tags describe and the compiler cannot
// check.

func TestCreateAccesskey(t *testing.T) {
	mock := NewMock().On("CreateAccessKey", &iam.CreateAccessKeyOutput{
		AccessKey: &iamtypes.AccessKey{
			AccessKeyId:     awssdk.String("AKIAIOSFODNN7EXAMPLE"),
			SecretAccessKey: awssdk.String("wJalrXUtnFEMI"),
		},
	})

	// no-prompt keeps the command from asking whether to write the new key into
	// ~/.aws/credentials, which would read from stdin during the test.
	Template("create accesskey user=jsmith no-prompt=true").
		Mock(mock).
		ExpectCalls("CreateAccessKey").
		ExpectCommandResult("AKIAIOSFODNN7EXAMPLE").
		Run(t)

	in := mock.InputFor("CreateAccessKey").(*iam.CreateAccessKeyInput)
	if got := awssdk.ToString(in.UserName); got != "jsmith" {
		t.Errorf("UserName: got %q, want jsmith", got)
	}
}

func TestDeleteAccesskey(t *testing.T) {
	mock := NewMock().On("DeleteAccessKey", &iam.DeleteAccessKeyOutput{})

	Template("delete accesskey id=AKIAIOSFODNN7EXAMPLE user=jsmith").
		Mock(mock).
		ExpectCalls("DeleteAccessKey").
		Run(t)

	in := mock.InputFor("DeleteAccessKey").(*iam.DeleteAccessKeyInput)
	if got := awssdk.ToString(in.AccessKeyId); got != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("AccessKeyId: got %q", got)
	}
}

func TestDeleteGroup(t *testing.T) {
	mock := NewMock().On("DeleteGroup", &iam.DeleteGroupOutput{})

	Template("delete group name=devs").
		Mock(mock).
		ExpectCalls("DeleteGroup").
		Run(t)

	in := mock.InputFor("DeleteGroup").(*iam.DeleteGroupInput)
	if got := awssdk.ToString(in.GroupName); got != "devs" {
		t.Errorf("GroupName: got %q, want devs", got)
	}
}

func TestAttachAndDetachUser(t *testing.T) {
	attach := NewMock().On("AddUserToGroup", &iam.AddUserToGroupOutput{})
	Template("attach user group=devs name=jsmith").
		Mock(attach).
		ExpectCalls("AddUserToGroup").
		ExpectRevert("detach user group=devs name=jsmith").
		Run(t)

	in := attach.InputFor("AddUserToGroup").(*iam.AddUserToGroupInput)
	if got := awssdk.ToString(in.GroupName); got != "devs" {
		t.Errorf("GroupName: got %q, want devs", got)
	}
	if got := awssdk.ToString(in.UserName); got != "jsmith" {
		t.Errorf("UserName: got %q, want jsmith", got)
	}

	detach := NewMock().On("RemoveUserFromGroup", &iam.RemoveUserFromGroupOutput{})
	Template("detach user group=devs name=jsmith").
		Mock(detach).
		ExpectCalls("RemoveUserFromGroup").
		Run(t)
}

func TestCreateLoginprofile(t *testing.T) {
	mock := NewMock().On("CreateLoginProfile", &iam.CreateLoginProfileOutput{
		LoginProfile: &iamtypes.LoginProfile{UserName: awssdk.String("jsmith")},
	})

	Template("create loginprofile username=jsmith password=S3curePassw0rd").
		Mock(mock).
		ExpectCalls("CreateLoginProfile").
		Run(t)

	in := mock.InputFor("CreateLoginProfile").(*iam.CreateLoginProfileInput)
	if got := awssdk.ToString(in.UserName); got != "jsmith" {
		t.Errorf("UserName: got %q, want jsmith", got)
	}
	if got := awssdk.ToString(in.Password); got != "S3curePassw0rd" {
		t.Errorf("Password did not reach the input: got %q", got)
	}
}

func TestUpdateAndDeleteLoginprofile(t *testing.T) {
	update := NewMock().On("UpdateLoginProfile", &iam.UpdateLoginProfileOutput{})
	Template("update loginprofile username=jsmith password=N3wPassw0rd").
		Mock(update).
		ExpectCalls("UpdateLoginProfile").
		Run(t)

	del := NewMock().On("DeleteLoginProfile", &iam.DeleteLoginProfileOutput{})
	Template("delete loginprofile username=jsmith").
		Mock(del).
		ExpectCalls("DeleteLoginProfile").
		Run(t)

	in := del.InputFor("DeleteLoginProfile").(*iam.DeleteLoginProfileInput)
	if got := awssdk.ToString(in.UserName); got != "jsmith" {
		t.Errorf("UserName: got %q, want jsmith", got)
	}
}

func TestCreateInstanceprofile(t *testing.T) {
	mock := NewMock().On("CreateInstanceProfile", &iam.CreateInstanceProfileOutput{
		InstanceProfile: &iamtypes.InstanceProfile{
			InstanceProfileName: awssdk.String("my-profile"),
		},
	})

	Template("create instanceprofile name=my-profile").
		Mock(mock).
		ExpectCalls("CreateInstanceProfile").
		ExpectRevert("delete instanceprofile name=my-profile").
		Run(t)
}

func TestDeleteInstanceprofile(t *testing.T) {
	mock := NewMock().On("DeleteInstanceProfile", &iam.DeleteInstanceProfileOutput{})

	Template("delete instanceprofile name=my-profile").
		Mock(mock).
		ExpectCalls("DeleteInstanceProfile").
		Run(t)
}

// attach policy accepts either an ARN or an access+service pair, and resolves the
// latter to a managed policy ARN, so both branches are worth covering.
func TestAttachPolicyToUserByArn(t *testing.T) {
	mock := NewMock().On("AttachUserPolicy", &iam.AttachUserPolicyOutput{})

	Template("attach policy user=jsmith arn=arn:aws:iam::aws:policy/ReadOnlyAccess").
		Mock(mock).
		ExpectCalls("AttachUserPolicy").
		Run(t)

	in := mock.InputFor("AttachUserPolicy").(*iam.AttachUserPolicyInput)
	if got := awssdk.ToString(in.PolicyArn); got != "arn:aws:iam::aws:policy/ReadOnlyAccess" {
		t.Errorf("PolicyArn: got %q", got)
	}
	if got := awssdk.ToString(in.UserName); got != "jsmith" {
		t.Errorf("UserName: got %q, want jsmith", got)
	}
}

func TestAttachPolicyToRoleAndGroup(t *testing.T) {
	role := NewMock().On("AttachRolePolicy", &iam.AttachRolePolicyOutput{})
	Template("attach policy role=my-role arn=arn:aws:iam::aws:policy/ReadOnlyAccess").
		Mock(role).
		ExpectCalls("AttachRolePolicy").
		Run(t)

	group := NewMock().On("AttachGroupPolicy", &iam.AttachGroupPolicyOutput{})
	Template("attach policy group=devs arn=arn:aws:iam::aws:policy/ReadOnlyAccess").
		Mock(group).
		ExpectCalls("AttachGroupPolicy").
		Run(t)

	in := group.InputFor("AttachGroupPolicy").(*iam.AttachGroupPolicyInput)
	if got := awssdk.ToString(in.GroupName); got != "devs" {
		t.Errorf("GroupName: got %q, want devs", got)
	}
}

func TestDetachPolicyFromUser(t *testing.T) {
	mock := NewMock().On("DetachUserPolicy", &iam.DetachUserPolicyOutput{})

	Template("detach policy user=jsmith arn=arn:aws:iam::aws:policy/ReadOnlyAccess").
		Mock(mock).
		ExpectCalls("DetachUserPolicy").
		Run(t)
}

func TestCreatePolicyBuildsDocument(t *testing.T) {
	mock := NewMock().On("CreatePolicy", &iam.CreatePolicyOutput{
		Policy: &iamtypes.Policy{Arn: awssdk.String("arn:aws:iam::123456789012:policy/my-policy")},
	})

	Template("create policy name=my-policy effect=Allow action=s3:GetObject resource=arn:aws:s3:::my-bucket/*").
		Mock(mock).
		ExpectCalls("CreatePolicy").
		ExpectCommandResult("arn:aws:iam::123456789012:policy/my-policy").
		Run(t)

	in := mock.InputFor("CreatePolicy").(*iam.CreatePolicyInput)
	doc := awssdk.ToString(in.PolicyDocument)
	// The document is assembled as JSON; IAM requires these exact keys.
	for _, want := range []string{`"Version"`, `"Statement"`, `"Allow"`, "s3:GetObject", "my-bucket"} {
		if !strings.Contains(doc, want) {
			t.Errorf("policy document missing %s\n  got: %s", want, doc)
		}
	}
}

func TestUpdatePolicyCreatesNewVersion(t *testing.T) {
	mock := NewMock().
		// IAM allows at most five versions per policy, so the command lists the
		// existing ones before adding another.
		On("ListPolicyVersions", &iam.ListPolicyVersionsOutput{
			Versions: []iamtypes.PolicyVersion{{VersionId: awssdk.String("v1"), IsDefaultVersion: true}},
		}).
		// The existing document is fetched so the new version can be built from it.
		On("GetPolicyVersion", &iam.GetPolicyVersionOutput{
			PolicyVersion: &iamtypes.PolicyVersion{
				VersionId: awssdk.String("v1"),
				Document:  awssdk.String("%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Action%22%3A%22s3%3AGetObject%22%2C%22Resource%22%3A%22arn%3Aaws%3As3%3A%3A%3Amy-bucket/%2A%22%7D%5D%7D"),
			},
		}).
		On("CreatePolicyVersion", &iam.CreatePolicyVersionOutput{
			PolicyVersion: &iamtypes.PolicyVersion{VersionId: awssdk.String("v2")},
		})

	Template("update policy arn=arn:aws:iam::123456789012:policy/my-policy effect=Deny action=s3:PutObject resource=arn:aws:s3:::my-bucket/*").
		Mock(mock).
		ExpectCalls("ListPolicyVersions", "GetPolicyVersion", "CreatePolicyVersion").
		Run(t)

	in := mock.InputFor("CreatePolicyVersion").(*iam.CreatePolicyVersionInput)
	// A new version must be made the default, or the update has no effect.
	if !in.SetAsDefault {
		t.Error("expected SetAsDefault, otherwise the new version is inert")
	}
}

func TestDeletePolicy(t *testing.T) {
	mock := NewMock().On("DeletePolicy", &iam.DeletePolicyOutput{})

	Template("delete policy arn=arn:aws:iam::123456789012:policy/my-policy").
		Mock(mock).
		ExpectCalls("DeletePolicy").
		Run(t)
}

func TestAttachAndDetachMfadevice(t *testing.T) {
	attach := NewMock().On("EnableMFADevice", &iam.EnableMFADeviceOutput{})
	Template("attach mfadevice id=arn:aws:iam::123456789012:mfa/jsmith user=jsmith mfa-code-1=123456 mfa-code-2=654321 no-prompt=true").
		Mock(attach).
		ExpectCalls("EnableMFADevice").
		Run(t)

	in := attach.InputFor("EnableMFADevice").(*iam.EnableMFADeviceInput)
	if got := awssdk.ToString(in.AuthenticationCode1); got != "123456" {
		t.Errorf("AuthenticationCode1: got %q, want 123456", got)
	}
	if got := awssdk.ToString(in.AuthenticationCode2); got != "654321" {
		t.Errorf("AuthenticationCode2: got %q, want 654321", got)
	}

	detach := NewMock().On("DeactivateMFADevice", &iam.DeactivateMFADeviceOutput{})
	Template("detach mfadevice id=arn:aws:iam::123456789012:mfa/jsmith user=jsmith").
		Mock(detach).
		ExpectCalls("DeactivateMFADevice").
		Run(t)
}

func TestDeleteMfadevice(t *testing.T) {
	mock := NewMock().On("DeleteVirtualMFADevice", &iam.DeleteVirtualMFADeviceOutput{})

	Template("delete mfadevice id=arn:aws:iam::123456789012:mfa/jsmith").
		Mock(mock).
		ExpectCalls("DeleteVirtualMFADevice").
		Run(t)
}

// attach instanceprofile associates a profile with a running instance, which requires
// discovering any existing association first.
func TestAttachInstanceprofile(t *testing.T) {
	mock := NewMock().
		On("AssociateIamInstanceProfile", &ec2.AssociateIamInstanceProfileOutput{
			IamInstanceProfileAssociation: &ec2types.IamInstanceProfileAssociation{
				AssociationId: awssdk.String("iip-assoc-1234"),
			},
		})

	Template("attach instanceprofile instance=i-1234 name=my-profile").
		Mock(mock).
		ExpectCalls("AssociateIamInstanceProfile").
		Run(t)

	in := mock.InputFor("AssociateIamInstanceProfile").(*ec2.AssociateIamInstanceProfileInput)
	if got := awssdk.ToString(in.InstanceId); got != "i-1234" {
		t.Errorf("InstanceId: got %q, want i-1234", got)
	}
	if in.IamInstanceProfile == nil || awssdk.ToString(in.IamInstanceProfile.Name) != "my-profile" {
		t.Errorf("IamInstanceProfile: got %#v, want name my-profile", in.IamInstanceProfile)
	}
}
