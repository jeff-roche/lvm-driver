package services

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/klog/v2"
)

// IdentityService handles requests from the container orchestrator
// to determine readiness, capabilities, and identify the driver
type IdentityService struct {
	csi.UnimplementedIdentityServer
	name         string
	version      string
	capabilities []csi.PluginCapability_Service_Type
	ready        func() (bool, error)
}

// NewIdentityService returns a new IdentityService.
//
// ready is a function to check the plugin status.
// It should return non-nil error if the plugin is not healthy.
// If the plugin is not yet ready, it should return (false, nil).
// Otherwise, return (true, nil).
func NewIdentityService(name string, version string, ready func() (bool, error)) csi.IdentityServer {
	return &IdentityService{
		ready:   ready,
		name:    name,
		version: version,
		capabilities: []csi.PluginCapability_Service_Type{
			csi.PluginCapability_Service_UNKNOWN,
		},
	}
}

func (s IdentityService) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	klog.V(2).Info("received PluginInfoRequest")
	if s.name == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured")
	}

	if s.version == "" {
		return nil, status.Error(codes.Unavailable, "Driver is missing version")
	}

	resp := &csi.GetPluginInfoResponse{
		Name:          s.name,
		VendorVersion: s.version,
	}

	return resp, nil
}

func (s IdentityService) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	klog.V(2).Info("received GetPluginCapabilitiesRequest")
	capabilities := make([]*csi.PluginCapability, 0, len(s.capabilities))

	for _, cap := range s.capabilities {
		capabilities = append(capabilities, &csi.PluginCapability{
			// This should be replaced as we add capabilities
			Type: &csi.PluginCapability_Service_{
				Service: &csi.PluginCapability_Service{
					Type: cap,
				},
			},
		})
	}

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: capabilities,
	}, nil
}

func (s IdentityService) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	ok, err := s.ready()
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &csi.ProbeResponse{
		Ready: &wrapperspb.BoolValue{
			Value: ok,
		},
	}, nil
}
