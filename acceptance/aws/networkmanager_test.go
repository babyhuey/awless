package awsat

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/networkmanager"
	networkmanagertypes "github.com/aws/aws-sdk-go-v2/service/networkmanager/types"
)

func TestCreateGlobalnetwork(t *testing.T) {
	mock := NewMock().On("CreateGlobalNetwork", &networkmanager.CreateGlobalNetworkOutput{
		GlobalNetwork: &networkmanagertypes.GlobalNetwork{
			GlobalNetworkId: awssdk.String("global-network-1234"),
		},
	})

	Template("create globalnetwork description='Production WAN'").
		Mock(mock).
		ExpectCalls("CreateGlobalNetwork").
		ExpectCommandResult("global-network-1234").
		ExpectRevert("delete globalnetwork id=global-network-1234").
		Run(t)

	in := mock.InputFor("CreateGlobalNetwork").(*networkmanager.CreateGlobalNetworkInput)
	if got := awssdk.ToString(in.Description); got != "Production WAN" {
		t.Errorf("Description: got %q", got)
	}
}

func TestDeleteGlobalnetwork(t *testing.T) {
	mock := NewMock().On("DeleteGlobalNetwork", &networkmanager.DeleteGlobalNetworkOutput{})

	Template("delete globalnetwork id=global-network-1234").
		Mock(mock).
		ExpectCalls("DeleteGlobalNetwork").
		Run(t)

	in := mock.InputFor("DeleteGlobalNetwork").(*networkmanager.DeleteGlobalNetworkInput)
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
}

func TestCreateCorenetwork(t *testing.T) {
	mock := NewMock().On("CreateCoreNetwork", &networkmanager.CreateCoreNetworkOutput{
		CoreNetwork: &networkmanagertypes.CoreNetwork{
			CoreNetworkId: awssdk.String("core-network-5678"),
		},
	})

	Template("create corenetwork global-network=global-network-1234 description='Core'").
		Mock(mock).
		ExpectCalls("CreateCoreNetwork").
		ExpectCommandResult("core-network-5678").
		ExpectRevert("delete corenetwork id=core-network-5678").
		Run(t)

	in := mock.InputFor("CreateCoreNetwork").(*networkmanager.CreateCoreNetworkInput)
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
	if got := awssdk.ToString(in.Description); got != "Core" {
		t.Errorf("Description: got %q", got)
	}
}

func TestDeleteCorenetwork(t *testing.T) {
	mock := NewMock().On("DeleteCoreNetwork", &networkmanager.DeleteCoreNetworkOutput{})

	Template("delete corenetwork id=core-network-5678").
		Mock(mock).
		ExpectCalls("DeleteCoreNetwork").
		Run(t)

	in := mock.InputFor("DeleteCoreNetwork").(*networkmanager.DeleteCoreNetworkInput)
	if got := awssdk.ToString(in.CoreNetworkId); got != "core-network-5678" {
		t.Errorf("CoreNetworkId: got %q", got)
	}
}

func TestCreateNetworkmanagersite(t *testing.T) {
	mock := NewMock().On("CreateSite", &networkmanager.CreateSiteOutput{
		Site: &networkmanagertypes.Site{
			SiteId: awssdk.String("site-1234"),
		},
	})

	Template("create networkmanagersite global-network=global-network-1234 description='London office'").
		Mock(mock).
		ExpectCalls("CreateSite").
		ExpectCommandResult("site-1234").
		ExpectRevert("delete networkmanagersite global-network=global-network-1234 id=site-1234").
		Run(t)

	in := mock.InputFor("CreateSite").(*networkmanager.CreateSiteInput)
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
	if got := awssdk.ToString(in.Description); got != "London office" {
		t.Errorf("Description: got %q", got)
	}
}

func TestDeleteNetworkmanagersite(t *testing.T) {
	mock := NewMock().On("DeleteSite", &networkmanager.DeleteSiteOutput{})

	Template("delete networkmanagersite id=site-1234 global-network=global-network-1234").
		Mock(mock).
		ExpectCalls("DeleteSite").
		Run(t)

	in := mock.InputFor("DeleteSite").(*networkmanager.DeleteSiteInput)
	if got := awssdk.ToString(in.SiteId); got != "site-1234" {
		t.Errorf("SiteId: got %q", got)
	}
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
}

func TestCreateNetworkmanagerdevice(t *testing.T) {
	mock := NewMock().On("CreateDevice", &networkmanager.CreateDeviceOutput{
		Device: &networkmanagertypes.Device{
			DeviceId: awssdk.String("device-1234"),
		},
	})

	Template("create networkmanagerdevice global-network=global-network-1234 model='ISR 4451' vendor=Cisco site=site-1234").
		Mock(mock).
		ExpectCalls("CreateDevice").
		ExpectCommandResult("device-1234").
		ExpectRevert("delete networkmanagerdevice global-network=global-network-1234 id=device-1234").
		Run(t)

	in := mock.InputFor("CreateDevice").(*networkmanager.CreateDeviceInput)
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
	if got := awssdk.ToString(in.Model); got != "ISR 4451" {
		t.Errorf("Model: got %q", got)
	}
	if got := awssdk.ToString(in.Vendor); got != "Cisco" {
		t.Errorf("Vendor: got %q", got)
	}
	if got := awssdk.ToString(in.SiteId); got != "site-1234" {
		t.Errorf("SiteId: got %q", got)
	}
}

func TestDeleteNetworkmanagerdevice(t *testing.T) {
	mock := NewMock().On("DeleteDevice", &networkmanager.DeleteDeviceOutput{})

	Template("delete networkmanagerdevice id=device-1234 global-network=global-network-1234").
		Mock(mock).
		ExpectCalls("DeleteDevice").
		Run(t)

	in := mock.InputFor("DeleteDevice").(*networkmanager.DeleteDeviceInput)
	if got := awssdk.ToString(in.DeviceId); got != "device-1234" {
		t.Errorf("DeviceId: got %q", got)
	}
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
}

func TestCreateNetworkmanagerlink(t *testing.T) {
	mock := NewMock().On("CreateLink", &networkmanager.CreateLinkOutput{
		Link: &networkmanagertypes.Link{
			LinkId: awssdk.String("link-1234"),
		},
	})

	Template("create networkmanagerlink global-network=global-network-1234 site=site-1234 type=broadband").
		Mock(mock).
		ExpectCalls("CreateLink").
		ExpectCommandResult("link-1234").
		ExpectRevert("delete networkmanagerlink global-network=global-network-1234 id=link-1234").
		Run(t)

	in := mock.InputFor("CreateLink").(*networkmanager.CreateLinkInput)
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
	if got := awssdk.ToString(in.SiteId); got != "site-1234" {
		t.Errorf("SiteId: got %q", got)
	}
	if got := awssdk.ToString(in.Type); got != "broadband" {
		t.Errorf("Type: got %q", got)
	}
}

func TestDeleteNetworkmanagerlink(t *testing.T) {
	mock := NewMock().On("DeleteLink", &networkmanager.DeleteLinkOutput{})

	Template("delete networkmanagerlink id=link-1234 global-network=global-network-1234").
		Mock(mock).
		ExpectCalls("DeleteLink").
		Run(t)

	in := mock.InputFor("DeleteLink").(*networkmanager.DeleteLinkInput)
	if got := awssdk.ToString(in.LinkId); got != "link-1234" {
		t.Errorf("LinkId: got %q", got)
	}
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
}

func TestCreateNetworkmanagerconnection(t *testing.T) {
	mock := NewMock().On("CreateConnection", &networkmanager.CreateConnectionOutput{
		Connection: &networkmanagertypes.Connection{
			ConnectionId: awssdk.String("connection-1234"),
		},
	})

	Template("create networkmanagerconnection global-network=global-network-1234 device=device-1234 connected-device=device-5678").
		Mock(mock).
		ExpectCalls("CreateConnection").
		ExpectCommandResult("connection-1234").
		ExpectRevert("delete networkmanagerconnection global-network=global-network-1234 id=connection-1234").
		Run(t)

	in := mock.InputFor("CreateConnection").(*networkmanager.CreateConnectionInput)
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
	if got := awssdk.ToString(in.DeviceId); got != "device-1234" {
		t.Errorf("DeviceId: got %q", got)
	}
	if got := awssdk.ToString(in.ConnectedDeviceId); got != "device-5678" {
		t.Errorf("ConnectedDeviceId: got %q", got)
	}
}

func TestDeleteNetworkmanagerconnection(t *testing.T) {
	mock := NewMock().On("DeleteConnection", &networkmanager.DeleteConnectionOutput{})

	Template("delete networkmanagerconnection id=connection-1234 global-network=global-network-1234").
		Mock(mock).
		ExpectCalls("DeleteConnection").
		Run(t)

	in := mock.InputFor("DeleteConnection").(*networkmanager.DeleteConnectionInput)
	if got := awssdk.ToString(in.ConnectionId); got != "connection-1234" {
		t.Errorf("ConnectionId: got %q", got)
	}
	if got := awssdk.ToString(in.GlobalNetworkId); got != "global-network-1234" {
		t.Errorf("GlobalNetworkId: got %q", got)
	}
}

// Empty-response cases.

func TestCreateGlobalnetworkEmptyResponse(t *testing.T) {
	mock := NewMock().On("CreateGlobalNetwork", &networkmanager.CreateGlobalNetworkOutput{})

	Template("create globalnetwork description='test'").
		Mock(mock).
		ExpectCalls("CreateGlobalNetwork").
		ExpectCommandResult("").
		Run(t)
}

func TestCreateCorenetworkEmptyResponse(t *testing.T) {
	mock := NewMock().On("CreateCoreNetwork", &networkmanager.CreateCoreNetworkOutput{})

	Template("create corenetwork global-network=global-network-1234").
		Mock(mock).
		ExpectCalls("CreateCoreNetwork").
		ExpectCommandResult("").
		Run(t)
}
