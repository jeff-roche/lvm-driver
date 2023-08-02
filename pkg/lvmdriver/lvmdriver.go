/*
Copyright © 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lvmdriver

import (
	"runtime"

	"k8s.io/klog/v2"
)

type LvmDriverOptions struct {
	NodeID     string
	DriverName string
	Endpoint   string
}

type LvmDriver struct {
	name     string
	nodeID   string
	endpoint string
	version  string
}

func NewLvmDriver(options *LvmDriverOptions) *LvmDriver {
	klog.V(1).Infof("Driver: %v version :%v", options.DriverName, driverVersion)

	lvmd := &LvmDriver{
		name:     options.DriverName,
		version:  driverVersion,
		nodeID:   options.NodeID,
		endpoint: options.Endpoint,
	}

	return lvmd
}

func (driver *LvmDriver) Run() {
	versionInfo, err := GetVersionYAML(driver.name)
	if err != nil {
		klog.Fatalf("%v", err)
	}
	klog.V(1).Infof("\nDRIVER INFORMATION:\n-------------------\n%s\n\nStreaming logs below:", versionInfo)

	// TODO: remove this once we have meaningful servers
	// Wait forever, letting other work commence if needed
	for {
		runtime.Gosched()
	}
}
