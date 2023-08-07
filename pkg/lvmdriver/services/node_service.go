package services

import (
	"context"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	ctrl "sigs.k8s.io/controller-runtime"
)

type NodeService struct {
	csi.UnimplementedNodeServer
	logger       logr.Logger
	mtx          sync.RWMutex // Need to handle concurrent system calls
	capabilities []csi.NodeServiceCapability_RPC_Type
}

func NewNodeService(name string) csi.NodeServer {
	return &NodeService{
		logger: ctrl.Log.WithName(name).WithName("node_service"),
		capabilities: []csi.NodeServiceCapability_RPC_Type{
			csi.NodeServiceCapability_RPC_UNKNOWN,
		},
	}
}

func (n *NodeService) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	n.logger.Info("GetNodeCapabilities", "req", req.String())
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
	n.logger.Info("NodePublishVolume", "req", req.String())
	n.mtx.Lock()
	defer n.mtx.Unlock()

	return nil, status.Error(codes.Unimplemented, "NodePublishVolume is not yet implemented")
}

func (n *NodeService) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	n.logger.Info("NodeUnpublishVolume", "req", req.String())
	n.mtx.Lock()
	defer n.mtx.Unlock()

	return nil, status.Error(codes.Unimplemented, "NodeUnpublishVolume is not yet implemented")
}
