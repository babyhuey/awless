package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	sdtypes "github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"
)

// A namespace is created by one of two different API calls depending on whether a VPC is
// given, so the command dispatches rather than mapping to a single call. These two tests are
// the reason that dispatch exists.
func TestCreateNamespacePrivateDNS(t *testing.T) {
	mock := NewMock().On("CreatePrivateDnsNamespace", &servicediscovery.CreatePrivateDnsNamespaceOutput{
		OperationId: awssdk.String("op-1234"),
	})

	Template("create namespace name=internal vpc=vpc-1234 description=Internal").
		Mock(mock).
		// Only the private-DNS call should be made; the HTTP one is a different namespace
		// type entirely.
		ExpectCalls("CreatePrivateDnsNamespace").
		// Creation is asynchronous, so what comes back is an operation id.
		ExpectCommandResult("op-1234").
		Run(t)

	in := mock.InputFor("CreatePrivateDnsNamespace").(*servicediscovery.CreatePrivateDnsNamespaceInput)
	if got := awssdk.ToString(in.Name); got != "internal" {
		t.Errorf("Name: got %q, want internal", got)
	}
	if got := awssdk.ToString(in.Vpc); got != "vpc-1234" {
		t.Errorf("Vpc: got %q, want vpc-1234", got)
	}
}

func TestCreateNamespaceHTTPOnly(t *testing.T) {
	mock := NewMock().On("CreateHttpNamespace", &servicediscovery.CreateHttpNamespaceOutput{
		OperationId: awssdk.String("op-5678"),
	})

	Template("create namespace name=http-only description=HTTP").
		Mock(mock).
		ExpectCalls("CreateHttpNamespace").
		ExpectCommandResult("op-5678").
		Run(t)

	in := mock.InputFor("CreateHttpNamespace").(*servicediscovery.CreateHttpNamespaceInput)
	if got := awssdk.ToString(in.Name); got != "http-only" {
		t.Errorf("Name: got %q, want http-only", got)
	}
}

func TestDeleteNamespace(t *testing.T) {
	mock := NewMock().On("DeleteNamespace", &servicediscovery.DeleteNamespaceOutput{})

	Template("delete namespace id=ns-1234abcd").
		Mock(mock).
		ExpectCalls("DeleteNamespace").
		Run(t)

	in := mock.InputFor("DeleteNamespace").(*servicediscovery.DeleteNamespaceInput)
	if got := awssdk.ToString(in.Id); got != "ns-1234abcd" {
		t.Errorf("Id: got %q, want ns-1234abcd", got)
	}
}

// The DNS configuration decides which records the service registers, and is a document.
func TestCreateDiscoveryservice(t *testing.T) {
	dns := docFile(t, "dns.json",
		`{"routingPolicy": "MULTIVALUE", "dnsRecords": [{"type": "A", "ttl": 60}]}`)

	mock := NewMock().On("CreateService", &servicediscovery.CreateServiceOutput{
		Service: &sdtypes.Service{Id: awssdk.String("srv-1234abcd")},
	})

	Template("create discoveryservice name=api namespace=ns-1234abcd dns-file=" + dns).
		Mock(mock).
		ExpectCalls("CreateService").
		ExpectCommandResult("srv-1234abcd").
		ExpectRevert("delete discoveryservice id=srv-1234abcd").
		Run(t)

	in := mock.InputFor("CreateService").(*servicediscovery.CreateServiceInput)
	if got := awssdk.ToString(in.Name); got != "api" {
		t.Errorf("Name: got %q, want api", got)
	}
	if got := awssdk.ToString(in.NamespaceId); got != "ns-1234abcd" {
		t.Errorf("NamespaceId: got %q", got)
	}
	if in.DnsConfig == nil {
		t.Fatal("DnsConfig was not decoded")
	}
	if got := string(in.DnsConfig.RoutingPolicy); got != "MULTIVALUE" {
		t.Errorf("DnsConfig.RoutingPolicy: got %q", got)
	}
	if len(in.DnsConfig.DnsRecords) != 1 {
		t.Fatalf("DnsRecords: got %d, want 1", len(in.DnsConfig.DnsRecords))
	}
	if got := awssdk.ToInt64(in.DnsConfig.DnsRecords[0].TTL); got != 60 {
		t.Errorf("DnsRecords[0].TTL: got %d, want 60", got)
	}
}

func TestDeleteDiscoveryservice(t *testing.T) {
	mock := NewMock().On("DeleteService", &servicediscovery.DeleteServiceOutput{})

	Template("delete discoveryservice id=srv-1234abcd").
		Mock(mock).
		ExpectCalls("DeleteService").
		Run(t)

	in := mock.InputFor("DeleteService").(*servicediscovery.DeleteServiceInput)
	if got := awssdk.ToString(in.Id); got != "srv-1234abcd" {
		t.Errorf("Id: got %q", got)
	}
}
