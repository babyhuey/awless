package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	idptypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func TestCreateUserpool(t *testing.T) {
	mock := NewMock().On("CreateUserPool", &cognitoidentityprovider.CreateUserPoolOutput{
		UserPool: &idptypes.UserPoolType{Id: awssdk.String("us-west-2_abcDEF123")},
	})

	Template("create userpool name=customers deletion-protection=ACTIVE auto-verified=email").
		Mock(mock).
		ExpectCalls("CreateUserPool").
		// The id is what delete takes, since pool names are not unique.
		ExpectCommandResult("us-west-2_abcDEF123").
		ExpectRevert("delete userpool id=us-west-2_abcDEF123").
		Run(t)

	in := mock.InputFor("CreateUserPool").(*cognitoidentityprovider.CreateUserPoolInput)
	if got := awssdk.ToString(in.PoolName); got != "customers" {
		t.Errorf("PoolName: got %q, want customers", got)
	}
	if got := string(in.DeletionProtection); got != "ACTIVE" {
		t.Errorf("DeletionProtection: got %q, want ACTIVE", got)
	}
	if len(in.AutoVerifiedAttributes) != 1 || string(in.AutoVerifiedAttributes[0]) != "email" {
		t.Errorf("AutoVerifiedAttributes: got %v", in.AutoVerifiedAttributes)
	}
}

// The password policy is a nested document.
func TestCreateUserpoolWithPolicies(t *testing.T) {
	policies := docFile(t, "policies.json",
		`{"passwordPolicy": {"minimumLength": 12, "requireNumbers": true, "requireSymbols": true}}`)

	mock := NewMock().On("CreateUserPool", &cognitoidentityprovider.CreateUserPoolOutput{})

	Template("create userpool name=customers policies-file=" + policies).
		Mock(mock).
		ExpectCalls("CreateUserPool").
		ExpectCommandResult("customers").
		Run(t)

	in := mock.InputFor("CreateUserPool").(*cognitoidentityprovider.CreateUserPoolInput)
	if in.Policies == nil || in.Policies.PasswordPolicy == nil {
		t.Fatal("the password policy was not decoded")
	}
	if got := awssdk.ToInt32(in.Policies.PasswordPolicy.MinimumLength); got != 12 {
		t.Errorf("MinimumLength: got %d, want 12", got)
	}
	if !in.Policies.PasswordPolicy.RequireNumbers {
		t.Error("expected RequireNumbers to be true")
	}
}

func TestDeleteUserpool(t *testing.T) {
	mock := NewMock().On("DeleteUserPool", &cognitoidentityprovider.DeleteUserPoolOutput{})

	Template("delete userpool id=us-west-2_abcDEF123").
		Mock(mock).
		ExpectCalls("DeleteUserPool").
		Run(t)

	in := mock.InputFor("DeleteUserPool").(*cognitoidentityprovider.DeleteUserPoolInput)
	if got := awssdk.ToString(in.UserPoolId); got != "us-west-2_abcDEF123" {
		t.Errorf("UserPoolId: got %q", got)
	}
}

// Unauthenticated identities hand AWS credentials to anyone who asks, so the flag is
// required rather than defaulted either way.
func TestCreateIdentitypoolRequiresTheUnauthenticatedDecision(t *testing.T) {
	t.Run("omitting it is rejected", func(t *testing.T) {
		err := Template("create identitypool name=customers-identities").
			Mock(NewMock()).
			RunExpectingError(t)
		if err == nil {
			t.Fatal("expected allow-unauthenticated to be required")
		}
	})

	t.Run("false is passed through", func(t *testing.T) {
		mock := NewMock().On("CreateIdentityPool", &cognitoidentity.CreateIdentityPoolOutput{
			IdentityPoolId: awssdk.String("us-west-2:1111"),
		})
		Template("create identitypool name=customers-identities allow-unauthenticated=false").
			Mock(mock).
			ExpectCalls("CreateIdentityPool").
			ExpectCommandResult("us-west-2:1111").
			ExpectRevert("delete identitypool id=us-west-2:1111").
			Run(t)

		in := mock.InputFor("CreateIdentityPool").(*cognitoidentity.CreateIdentityPoolInput)
		if in.AllowUnauthenticatedIdentities {
			t.Error("expected AllowUnauthenticatedIdentities to be false")
		}
		if got := awssdk.ToString(in.IdentityPoolName); got != "customers-identities" {
			t.Errorf("IdentityPoolName: got %q", got)
		}
	})
}

func TestDeleteIdentitypool(t *testing.T) {
	mock := NewMock().On("DeleteIdentityPool", &cognitoidentity.DeleteIdentityPoolOutput{})

	Template("delete identitypool id=us-west-2:1111").
		Mock(mock).
		ExpectCalls("DeleteIdentityPool").
		Run(t)

	in := mock.InputFor("DeleteIdentityPool").(*cognitoidentity.DeleteIdentityPoolInput)
	if got := awssdk.ToString(in.IdentityPoolId); got != "us-west-2:1111" {
		t.Errorf("IdentityPoolId: got %q", got)
	}
}
