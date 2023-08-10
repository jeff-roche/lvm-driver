package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEndpoint(t *testing.T) {
	tests := []struct {
		desc      string
		endpoint  string
		protocol  string
		addr      string
		expectErr bool
	}{
		{
			desc:      "tcp address",
			endpoint:  "tcp://128.0.0.1",
			protocol:  "tcp",
			addr:      "128.0.0.1",
			expectErr: false,
		},
		{
			desc:      "unix socket address",
			endpoint:  "unix://tmp/foobar",
			protocol:  "unix",
			addr:      "tmp/foobar",
			expectErr: false,
		},
		{
			desc:      "unsupported address",
			endpoint:  "foo://bar",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p, a, err := ParseEndpoint(test.endpoint)

			if err != nil {
				if test.expectErr {
					return
				}

				t.Errorf("detected an error when one was not expected: %s", err)
			}

			if err == nil && test.expectErr {
				t.Error("no error detected when one was expected")
			}

			if p != test.protocol {
				t.Errorf("wrong protocol returned. expected %s, got %s", test.protocol, p)
			}

			if a != test.addr {
				t.Errorf("wrong address returned. expected %s, got %s", test.addr, a)
			}
		})
	}
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		desc          string
		method        string
		expectedLevel int32
	}{
		{
			desc:          "identity probe log level",
			method:        "/csi.v1.Identity/Probe",
			expectedLevel: 8,
		},
		{
			desc:          "node get capabilities log level",
			method:        "/csi.v1.Node/NodeGetCapabilities",
			expectedLevel: 8,
		},
		{
			desc:          "node get volume stats log level",
			method:        "/csi.v1.Node/NodeGetVolumeStats",
			expectedLevel: 8,
		},
		{
			desc:          "standard log level",
			method:        "/csi.v1.foobar",
			expectedLevel: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			level := getLogLevel(test.method)
			assert.Equal(t, test.expectedLevel, level, "incorrect log level returned")
		})
	}
}
