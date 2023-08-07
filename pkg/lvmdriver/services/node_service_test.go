package services_test

import (
	"context"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/openshift/lvm-driver/pkg/lvmdriver/services"
	"github.com/stretchr/testify/assert"
)

func TestNodeGetCapabilites(t *testing.T) {
	validCapabilities := []csi.NodeServiceCapability_RPC_Type{
		csi.NodeServiceCapability_RPC_UNKNOWN,
	}

	nodeSvc := services.NewNodeService("foobar")
	req := &csi.NodeGetCapabilitiesRequest{}

	resp, err := nodeSvc.NodeGetCapabilities(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Capabilities)
	assert.Equal(t, len(validCapabilities), len(resp.Capabilities))

	returnedCapabilities := make([]csi.NodeServiceCapability_RPC_Type, 0, len(resp.Capabilities))

	for _, cap := range resp.Capabilities {
		returnedCapabilities = append(returnedCapabilities, cap.GetRpc().GetType())
	}

	// Make sure all valid and only vlaid capabilities were returned
	assert.ElementsMatch(t, returnedCapabilities, validCapabilities)
}

func TestNodePublishVolume(t *testing.T) {

}

func TestNodeUnpublishVolume(t *testing.T) {

}
