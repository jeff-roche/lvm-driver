package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/openshift/lvm-driver/pkg/lvmdriver/services"
	"github.com/stretchr/testify/assert"
)

func TestNodeGetInfo(t *testing.T) {
	driverName := "foo"
	nodeName := "bar"
	topologyKey := fmt.Sprintf("topology.%s/node", driverName)

	nodeSvc := services.NewNodeService(driverName, nodeName)
	req := &csi.NodeGetInfoRequest{}

	resp, err := nodeSvc.NodeGetInfo(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.NodeId, nodeName)
	assert.Equal(t, len(resp.AccessibleTopology.Segments), 1)
	assert.Contains(t, resp.AccessibleTopology.Segments, topologyKey)
	assert.Equal(t, resp.AccessibleTopology.Segments[topologyKey], nodeName)
}

func TestNodeGetCapabilites(t *testing.T) {
	validCapabilities := []csi.NodeServiceCapability_RPC_Type{
		csi.NodeServiceCapability_RPC_UNKNOWN,
	}

	nodeSvc := services.NewNodeService("NodeGetCapabilitiesSvc", "node_001")
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
	nodeSvc := services.NewNodeService("NodePublishVolumeSvc", "node_001")
	req := &csi.NodePublishVolumeRequest{}

	resp, err := nodeSvc.NodePublishVolume(context.Background(), req)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "Unimplemented")
}

func TestNodeUnpublishVolume(t *testing.T) {
	nodeSvc := services.NewNodeService("NodeUnpublishVolumeSvc", "node_001")
	req := &csi.NodeUnpublishVolumeRequest{}

	resp, err := nodeSvc.NodeUnpublishVolume(context.Background(), req)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "Unimplemented")
}
