/*
Copyright 2020-2021 The Kubernetes Authors.

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

package pci

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"k8s.io/klog/v2"

	"sigs.k8s.io/node-feature-discovery/pkg/api/feature"
	"sigs.k8s.io/node-feature-discovery/source"
)

// MaxLen is the maximums length showed in node labels
const (
	classMaxLen  = 4
	vendorMaxLen = 10
	deviceMaxLen = 35
)

var mandatoryDevAttrs = []string{"class", "vendor", "device", "subsystem_vendor", "subsystem_device"}
var optionalDevAttrs = []string{"sriov_totalvfs"}

// HiddenSuffixVendor is a list of meaningless vendor suffixes
var hiddenSuffixVendor = []string{
	"Corporation", "Corp.", "Corp.,", "Corp", "corp.", "Co.", "Co", "co.", "co.,", "CO.,", "Co.,", "Co.,Ltd", "Co.,LTD.", "Co.,Ltd.",
	"INC.", "INC", "Inc.", "Inc", "Inc,", "inc.",
	"Ltd.", "Ltd", "LTD.", "ltd.",
	"Technologies", "Technologies,", "Technology", "Technology,",
	"Information",
	"Company",
	"Group",
	"LLC", "LLC.",
}

// HiddenSuffixDevice is a list of meaningless device suffixes
var hiddenSuffixDevice = []string{
	"Processor",
	"Controller",
	"Adapter",
	"Integrated",
	"Technology",
	"Graphics",
	"Display",
	"PCI", "PCIe", "PCI-e", "PCI-to-PCI",
}

// Read a single PCI device attribute
// A PCI attribute in this context, maps to the corresponding sysfs file
func readSinglePciAttribute(devPath string, attrName string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(devPath, attrName))
	if err != nil {
		return "", fmt.Errorf("failed to read device attribute %s: %v", attrName, err)
	}
	// Strip whitespace and '0x' prefix
	attrVal := strings.TrimSpace(strings.TrimPrefix(string(data), "0x"))

	if attrName == "class" && len(attrVal) > 4 {
		// Take four first characters, so that the programming
		// interface identifier gets stripped from the raw class code
		attrVal = attrVal[0:4]
	}
	return attrVal, nil
}

// Read information of one PCI device
func readPciDevInfo(devPath string) (*feature.InstanceFeature, error) {
	attrs := make(map[string]string)
	for _, attr := range mandatoryDevAttrs {
		attrVal, err := readSinglePciAttribute(devPath, attr)
		if err != nil {
			return nil, fmt.Errorf("failed to read device %s: %s", attr, err)
		}
		attrs[attr] = attrVal
	}
	for _, attr := range optionalDevAttrs {
		attrVal, err := readSinglePciAttribute(devPath, attr)
		if err == nil {
			attrs[attr] = attrVal
		}
	}
	return feature.NewInstanceFeature(attrs), nil
}

// detectPci detects available PCI devices and retrieves their device attributes.
// An error is returned if reading any of the mandatory attributes fails.
func detectPci() ([]feature.InstanceFeature, error) {
	sysfsBasePath := source.SysfsDir.Path("bus/pci/devices")

	devices, err := ioutil.ReadDir(sysfsBasePath)
	if err != nil {
		return nil, err
	}

	// Iterate over devices
	devInfo := make([]feature.InstanceFeature, 0, len(devices))
	for _, device := range devices {
		info, err := readPciDevInfo(filepath.Join(sysfsBasePath, device.Name()))
		if err != nil {
			klog.Error(err)
			continue
		}
		devInfo = append(devInfo, *info)
	}

	return devInfo, nil
}

// Get pci-oriented human-readable class name
func getReadableClass(s string) string {
	classSplitList := strings.Fields(s)
	classResult := classSplitList[0]

	if len(classResult) > classMaxLen {
		classResult = classResult[:classMaxLen]
	}

	return classResult
}

// Get pci-oriented human-readable vendor name
func getReadableVendor(s string) string {
	vendorSplitList := strings.Fields(s)
	if len(vendorSplitList) == 0 {
		return ""
	}

	vendorLastItem := vendorSplitList[len(vendorSplitList)-1]
	vendorSplitList[0] = strings.Trim(vendorSplitList[0], ".,?")

	var vendorResult string
	if vendorLastItem == "ID)" {
		vendorResult = "Wrong-ID"
	} else if strings.HasPrefix(vendorLastItem, "[") {
		vendorResult = vendorLastItem[1 : len(vendorLastItem)-1]
	} else if len(vendorSplitList[0]) > vendorMaxLen {
		vendorResult = vendorSplitList[0][:vendorMaxLen]
	} else if len(vendorSplitList) == 1 || inArray(vendorSplitList[1], hiddenSuffixVendor) {
		vendorResult = vendorSplitList[0]
	} else {
		vendorSplitList[1] = strings.Trim(vendorSplitList[1], ".,?")
		vendorPreResult := vendorSplitList[0] + "-" + vendorSplitList[1]

		if len(vendorPreResult) > vendorMaxLen {
			vendorResult = vendorPreResult[:vendorMaxLen]
		} else {
			vendorResult = vendorPreResult
		}
	}

	vendorResult = strings.Trim(vendorResult, "-/&")
	vendorResult = strings.Replace(vendorResult, "/", "-", -1)

	return vendorResult
}

// Get pci-oriented human-readable device name
func getReadableDevice(s string) string {
	resultsDevice := s

	startIndex := strings.Index(s, "[")
	if startIndex != -1 {
		endIndex := strings.Index(s, "]")
		resultsDevice = s[startIndex+1 : endIndex]
	}

	deviceSplitList := strings.Fields(resultsDevice)
	if len(deviceSplitList) == 0 {
		return ""
	}

	resultsDevice = ""
	for _, device := range deviceSplitList {
		if !inArray(device, hiddenSuffixDevice) {
			device = strings.Trim(device, "-/&()")
			if len(device) != 0 {
				resultsDevice += device + "-"
			}
		}
	}

	resultsDevice = strings.Replace(resultsDevice, "(", "", -1)
	resultsDevice = strings.Replace(resultsDevice, ")", "", -1)

	if len(resultsDevice) > deviceMaxLen {
		resultsDevice = resultsDevice[:deviceMaxLen-1]
	}

	resultsDevice = strings.Replace(resultsDevice, "/", "-", -1)
	resultsDevice = strings.Replace(resultsDevice, "&", "-", -1)

	resultsDevice = strings.Trim(resultsDevice, "-")

	return resultsDevice
}

// Find target string in a given array
func inArray(target string, strArray []string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)

	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}
