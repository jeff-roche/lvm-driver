package utils_test

import (
	"testing"

	"github.com/openshift/lvm-driver/pkg/lvmdriver/utils"
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
			p, a, err := utils.ParseEndpoint(test.endpoint)

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
