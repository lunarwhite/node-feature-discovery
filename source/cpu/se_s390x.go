//go:build s390x
// +build s390x

/*
Copyright 2022 The Kubernetes Authors.

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

package cpu

import (
	"os"

	"sigs.k8s.io/node-feature-discovery/source"
)

func discoverSE() map[string]string {
	se := make(map[string]string)
	// This file is available in kernels >=5.12 + backports. Skip specifically
	// checking facilities and kernel command lines and just assume Secure
	// Execution to be unavailable or disabled if the file is not present.
	protVirtHost := source.SysfsDir.Path("firmware/uv/prot_virt_host")
	if content, err := os.ReadFile(protVirtHost); err == nil {
		if string(content) == "1\n" {
			se["enabled"] = "true"
		}
	}
	return se
}
