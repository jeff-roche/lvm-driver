package services_test

import (
	"context"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/openshift/lvm-driver/pkg/lvmdriver/services"
)

func TestGetPluginInfo(t *testing.T) {
	req := csi.GetPluginInfoRequest{}

	readyFunc := func() (bool, error) {
		return true, nil
	}

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

			if err != nil && !test.expectErr {
				t.Errorf("detected an error when one was not expected: %s", err)
			}

			if err == nil && test.expectErr {
				t.Error("no error detected when one was expected")
			}
		})

	}
}
