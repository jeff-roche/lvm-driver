package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

type NodeService struct {
	csi.UnimplementedNodeServer
	mtx          sync.RWMutex // Need to handle concurrent system calls
	capabilities []csi.NodeServiceCapability_RPC_Type
	nodeId       string
	topologies   *csi.Topology
}

func NewNodeService(name string, nodeId string) csi.NodeServer {
	topologyKey := fmt.Sprintf("topology.%s/node", name)

	return &NodeService{
		nodeId: nodeId,
		capabilities: []csi.NodeServiceCapability_RPC_Type{
			csi.NodeServiceCapability_RPC_UNKNOWN,
		},
		topologies: &csi.Topology{
			Segments: map[string]string{
				topologyKey: nodeId,
			},
		},
	}
}

func (n *NodeService) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.V(2).Infof("received %#v", *req)

	return &csi.NodeGetInfoResponse{
		NodeId:             n.nodeId,
		AccessibleTopology: n.topologies,
	}, nil
}

func (n *NodeService) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	klog.V(2).Infof("received %#v", req)
	csiCapabilities := make([]*csi.NodeServiceCapability, 0, len(n.capabilities))

	for _, cap := range n.capabilities {
		csiCapabilities = append(csiCapabilities, &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: cap,
				},
			},
		})
	}

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: csiCapabilities,
	}, nil
}

func (n *NodeService) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.V(2).Infof("received NodePublishVolumeRequest: %v", req)
	n.mtx.Lock()
	defer n.mtx.Unlock()

	return nil, status.Error(codes.Unimplemented, "NodePublishVolume is not yet implemented")
}

func (n *NodeService) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.V(2).Infof("received NodeUnpublishVolumeRequest: %v", req)
	n.mtx.Lock()
	defer n.mtx.Unlock()

	return nil, status.Error(codes.Unimplemented, "NodeUnpublishVolume is not yet implemented")
}
