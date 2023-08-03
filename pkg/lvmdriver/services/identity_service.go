package services

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	ctrl "sigs.k8s.io/controller-runtime"
)

// IdentityService handles requests from the container orchestrator
// to determine readiness, capabilities, and identify the driver
type IdentityService struct {
	csi.UnimplementedIdentityServer
	name         string
	version      string
	logger       logr.Logger
	capabilities []*csi.PluginCapability
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
		logger:  ctrl.Log.WithName(name).WithName("identity_service"),
		capabilities: []*csi.PluginCapability{
			{
				// This should be replaced as we add capabilities
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_UNKNOWN,
					},
				},
			},
		},
	}
}

func (s IdentityService) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	s.logger.Info("GetPluginInfo", "req", req.String())

	if s.name == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured")
	}

	if s.version == "" {
		return nil, status.Error(codes.Unavailable, "Driver is missing version")
	}

	return &csi.GetPluginInfoResponse{
		Name:          s.name,
		VendorVersion: s.version,
	}, nil
}

func (s IdentityService) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	s.logger.Info("GetPluginCapabilities", "req", req.String())
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: s.capabilities,
	}, nil
}

func (s IdentityService) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	s.logger.Info("Probe", "req", req.String())
	ok, err := s.ready()
	if err != nil {
		s.logger.Error(err, "probe failed")
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &csi.ProbeResponse{
		Ready: &wrapperspb.BoolValue{
			Value: ok,
		},
	}, nil
}
