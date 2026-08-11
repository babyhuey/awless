package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	wafv2types "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func TestCreateIpset(t *testing.T) {
	mock := NewMock().On("CreateIPSet", &wafv2.CreateIPSetOutput{
		Summary: &wafv2types.IPSetSummary{Name: awssdk.String("blocklist")},
	})

	Template("create ipset name=blocklist addresses=203.0.113.0/24,198.51.100.7/32 description=Blocked").
		Mock(mock).
		ExpectCalls("CreateIPSet").
		ExpectCommandResult("blocklist").
		ExpectRevert("delete ipset name=blocklist").
		Run(t)

	in := mock.InputFor("CreateIPSet").(*wafv2.CreateIPSetInput)
	if got := awssdk.ToString(in.Name); got != "blocklist" {
		t.Errorf("Name: got %q, want blocklist", got)
	}
	// Scope defaults to REGIONAL, which is the only one usable outside us-east-1.
	if in.Scope != wafv2types.ScopeRegional {
		t.Errorf("Scope: got %q, want REGIONAL", in.Scope)
	}
	if len(in.Addresses) != 2 || in.Addresses[0] != "203.0.113.0/24" {
		t.Errorf("Addresses: got %v", in.Addresses)
	}
	// A set holds one address family, and the version is inferred when not given.
	if in.IPAddressVersion != wafv2types.IPAddressVersionIpv4 {
		t.Errorf("IPAddressVersion: got %q, want IPV4", in.IPAddressVersion)
	}
}

// The address family is derived from the addresses, because AWS rejects a set that mixes
// them and requiring the user to restate it is a trap.
func TestCreateIpsetInfersIPv6(t *testing.T) {
	mock := NewMock().On("CreateIPSet", &wafv2.CreateIPSetOutput{})

	Template("create ipset name=v6block addresses=2001:db8::/32").
		Mock(mock).
		ExpectCalls("CreateIPSet").
		ExpectCommandResult("v6block").
		Run(t)

	in := mock.InputFor("CreateIPSet").(*wafv2.CreateIPSetInput)
	if in.IPAddressVersion != wafv2types.IPAddressVersionIpv6 {
		t.Errorf("IPAddressVersion: got %q, want IPV6", in.IPAddressVersion)
	}
}

func TestCreateIpsetWithExplicitScope(t *testing.T) {
	mock := NewMock().On("CreateIPSet", &wafv2.CreateIPSetOutput{})

	Template("create ipset name=edge addresses=203.0.113.0/24 scope=CLOUDFRONT ip-version=IPV4").
		Mock(mock).
		ExpectCalls("CreateIPSet").
		Run(t)

	in := mock.InputFor("CreateIPSet").(*wafv2.CreateIPSetInput)
	if in.Scope != wafv2types.ScopeCloudfront {
		t.Errorf("Scope: got %q, want CLOUDFRONT", in.Scope)
	}
}

// Deleting needs an id and a LockToken, neither of which the user has. Both are looked up
// by listing first, so the command makes two calls.
func TestDeleteIpsetLooksUpTheLockToken(t *testing.T) {
	mock := NewMock().
		On("ListIPSets", &wafv2.ListIPSetsOutput{
			IPSets: []wafv2types.IPSetSummary{
				{Name: awssdk.String("other"), Id: awssdk.String("id-other"), LockToken: awssdk.String("tok-other")},
				{Name: awssdk.String("blocklist"), Id: awssdk.String("id-1"), LockToken: awssdk.String("tok-1")},
			},
		}).
		On("DeleteIPSet", &wafv2.DeleteIPSetOutput{})

	Template("delete ipset name=blocklist").
		Mock(mock).
		ExpectCalls("ListIPSets", "DeleteIPSet").
		Run(t)

	in := mock.InputFor("DeleteIPSet").(*wafv2.DeleteIPSetInput)
	if got := awssdk.ToString(in.Id); got != "id-1" {
		t.Errorf("Id: got %q, want id-1 — the wrong set was matched", got)
	}
	if got := awssdk.ToString(in.LockToken); got != "tok-1" {
		t.Errorf("LockToken: got %q, want tok-1", got)
	}
	if got := awssdk.ToString(in.Name); got != "blocklist" {
		t.Errorf("Name: got %q, want blocklist", got)
	}
}

// A name that does not exist must fail with something a user can act on, rather than
// calling delete with an empty id.
func TestDeleteIpsetUnknownName(t *testing.T) {
	mock := NewMock().On("ListIPSets", &wafv2.ListIPSetsOutput{
		IPSets: []wafv2types.IPSetSummary{
			{Name: awssdk.String("other"), Id: awssdk.String("id-other")},
		},
	})

	err := Template("delete ipset name=missing").
		Mock(mock).
		RunExpectingError(t)
	if err == nil {
		t.Fatal("expected an error for an ip set that does not exist")
	}
}

func TestUpdateIpsetReplacesAddresses(t *testing.T) {
	mock := NewMock().
		On("ListIPSets", &wafv2.ListIPSetsOutput{
			IPSets: []wafv2types.IPSetSummary{
				{Name: awssdk.String("blocklist"), Id: awssdk.String("id-1"), LockToken: awssdk.String("tok-1")},
			},
		}).
		On("UpdateIPSet", &wafv2.UpdateIPSetOutput{})

	Template("update ipset name=blocklist addresses=203.0.113.0/24,192.0.2.0/24 description=Updated").
		Mock(mock).
		ExpectCalls("ListIPSets", "UpdateIPSet").
		Run(t)

	in := mock.InputFor("UpdateIPSet").(*wafv2.UpdateIPSetInput)
	if len(in.Addresses) != 2 || in.Addresses[1] != "192.0.2.0/24" {
		t.Errorf("Addresses: got %v", in.Addresses)
	}
	if got := awssdk.ToString(in.LockToken); got != "tok-1" {
		t.Errorf("LockToken: got %q, want tok-1", got)
	}
	if got := awssdk.ToString(in.Description); got != "Updated" {
		t.Errorf("Description: got %q, want Updated", got)
	}
}
