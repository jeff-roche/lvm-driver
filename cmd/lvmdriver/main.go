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

package main

import (
	"flag"
	"os"

	"github.com/openshift/lvm-driver/pkg/lvmdriver"
	"k8s.io/klog/v2"
)

var (
	endpoint   = flag.String("endpoint", "unix://tmp/csi.sock", "CSI Endpoint")
	nodeID     = flag.String("nodeid", "", "node id")
	driverName = flag.String("drivername", "lvm.redhat.com", "name of the driver")
)

func init() {
	flag.Set("logtostderr", "true")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	if *nodeID == "" {
		klog.Warning("nodeid is empty")
	}

	driverInit()

	os.Exit(0)
}

func driverInit() {
	opts := lvmdriver.LvmDriverOptions{
		NodeID:     *nodeID,
		DriverName: *driverName,
		Endpoint:   *endpoint,
	}

	driver := lvmdriver.NewLvmDriver(&opts)
	driver.Run()
}
