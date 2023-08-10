package services_test

import (
	"context"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/openshift/lvm-driver/pkg/lvmdriver/services"
	"github.com/stretchr/testify/assert"
)

var readyFunc = func() (bool, error) {
	return true, nil
}

func TestIdentityGetPluginInfo(t *testing.T) {
	req := csi.GetPluginInfoRequest{}

	tests := []struct {
		desc      string
		svc       csi.IdentityServer
		expectErr bool
	}{
		{
			desc:      "successful request",
			svc:       services.NewIdentityService("driverName", "driverVersion", readyFunc),
			expectErr: false,
		},
		{
			desc:      "driver name missing",
			svc:       services.NewIdentityService("", "driverVersion", readyFunc),
			expectErr: true,
		},
		{
			desc:      "driver version missing",
			svc:       services.NewIdentityService("driverName", "", readyFunc),
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			_, err := test.svc.GetPluginInfo(context.Background(), &req)

			if test.expectErr {
				assert.Error(t, err, "no error detected when one was expected")
			} else {
				assert.NoErrorf(t, err, "detected an error when one was not expected: %s", err)
			}
		})

	}
}

func TestIdentityProbe(t *testing.T) {
	idSvc := services.NewIdentityService("foo", "unix://bar", readyFunc)
	req := &csi.ProbeRequest{}

	resp, err := idSvc.Probe(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Ready.Value, true)
}

func TestIdentityGetPluginCapabilities(t *testing.T) {
	validCapabilities := []csi.PluginCapability_Service_Type{
		csi.PluginCapability_Service_UNKNOWN,
	}

	idSvc := services.NewIdentityService("foo", "unix://bar", readyFunc)
	req := &csi.GetPluginCapabilitiesRequest{}

	resp, err := idSvc.GetPluginCapabilities(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Capabilities)
	assert.Equal(t, len(validCapabilities), len(resp.Capabilities))

	returnedCapabilities := make([]csi.PluginCapability_Service_Type, 0, len(resp.Capabilities))

	for _, cap := range resp.Capabilities {
		returnedCapabilities = append(returnedCapabilities, cap.GetService().Type)
	}

	// Make sure all valid and only valid capabilities were returned
	assert.ElementsMatch(t, returnedCapabilities, validCapabilities)
}
