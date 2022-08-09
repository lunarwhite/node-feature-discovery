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

package usb

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"k8s.io/klog/v2"

	"sigs.k8s.io/node-feature-discovery/pkg/api/feature"
)

// MaxLen is the maximums length showed in node labels
const (
	vendorMaxLen = 10
	deviceMaxLen = 33
)

var devAttrs = []string{"class", "vendor", "device", "serial"}

// The USB device sysfs files do not have terribly user friendly names, map
// these for consistency with the PCI matcher.
var devAttrFileMap = map[string]string{
	"class":  "bDeviceClass",
	"device": "idProduct",
	"vendor": "idVendor",
	"serial": "serial",
}

// HiddenSuffixVendor is a list of meaningless vendor suffixes
var hiddenSuffixVendor = []string{
	"Corporation", "Corp.", "Corp.,", "Corp", "corp.", "Co.", "Co", "co.", "co.,", "CO.,", "Co.,", "Co.,Ltd", "Co.,LTD.", "Co.,Ltd.",
	"INC.", "INC", "Inc.", "Inc", "Inc,", "inc.",
	"Ltd.", "Ltd", "LTD.", "ltd.",
	"Technologies", "Technologies,", "Technology", "Technology,",
	"Information",
	"Electronics", "ELECTRONICS ", "Electric", "ELECTRIC",
	"Company",
	"Group",
	"LLC", "LLC.",
}

func readSingleUsbSysfsAttribute(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read device attribute %s: %v", filepath.Base(path), err)
	}

	attrVal := strings.TrimSpace(string(data))

	return attrVal, nil
}

// Read a single USB device attribute
// A USB attribute in this context, maps to the corresponding sysfs file
func readSingleUsbAttribute(devPath string, attrName string) (string, error) {
	return readSingleUsbSysfsAttribute(path.Join(devPath, devAttrFileMap[attrName]))
}

// Read information of one USB device
func readUsbDevInfo(devPath string) ([]feature.InstanceFeature, error) {
	instances := make([]feature.InstanceFeature, 0)
	attrs := make(map[string]string)

	for _, attr := range devAttrs {
		attrVal, _ := readSingleUsbAttribute(devPath, attr)
		if len(attrVal) > 0 {
			attrs[attr] = attrVal
		}
	}

	// USB devices encode their class information either at the device or the interface level. If the device class
	// is set, return as-is.
	if attrs["class"] != "00" {
		instances = append(instances, *feature.NewInstanceFeature(attrs))
	} else {
		// Otherwise, if a 00 is presented at the device level, descend to the interface level.
		interfaces, err := filepath.Glob(devPath + "/*/bInterfaceClass")
		if err != nil {
			return nil, err
		}

		// A device may, notably, have multiple interfaces with mixed classes, so we create a unique device for each
		// unique interface class.
		for _, intf := range interfaces {
			// Determine the interface class
			attrVal, err := readSingleUsbSysfsAttribute(intf)
			if err != nil {
				return nil, err
			}

			subdevAttrs := make(map[string]string, len(attrs))
			for k, v := range attrs {
				subdevAttrs[k] = v
			}
			subdevAttrs["class"] = attrVal

			instances = append(instances, *feature.NewInstanceFeature(subdevAttrs))
		}
	}

	return instances, nil
}

// detectUsb detects available USB devices and retrieves their device attributes.
func detectUsb() ([]feature.InstanceFeature, error) {
	// Unlike PCI, the USB sysfs interface includes entries not just for
	// devices. We work around this by globbing anything that includes a
	// valid product ID.
	const devPathGlob = "/sys/bus/usb/devices/*/idProduct"
	devPaths, err := filepath.Glob(devPathGlob)
	if err != nil {
		return nil, err
	}

	// Iterate over devices
	devInfo := make([]feature.InstanceFeature, 0)
	for _, devPath := range devPaths {
		devs, err := readUsbDevInfo(filepath.Dir(devPath))
		if err != nil {
			klog.Error(err)
			continue
		}

		devInfo = append(devInfo, devs...)
	}

	return devInfo, nil
}

// Get usb-oriented human-readable vendor name
func getReadableVendor(s string) string {
	vendorSplitList := strings.Fields(s)
	if len(vendorSplitList) == 0 {
		return ""
	}
	
	vendorLastItem := vendorSplitList[len(vendorSplitList)-1]
	vendorSplitList[0] = strings.Trim(vendorSplitList[0], ".,?[]()")

	var vendorResult string
	if vendorLastItem == "ID)" {
		vendorResult = "Wrong-ID"
	} else if len(vendorSplitList[0]) > vendorMaxLen {
		vendorResult = vendorSplitList[0][:vendorMaxLen]
	} else if len(vendorSplitList) == 1 || inArray(vendorSplitList[1], hiddenSuffixVendor) {
		vendorResult = vendorSplitList[0]
	} else {
		vendorSplitList[1] = strings.Trim(vendorSplitList[1], ".,?[]()")
		vendorPreResult := vendorSplitList[0] + "-" + vendorSplitList[1]

		if len(vendorPreResult) > vendorMaxLen {
			vendorResult = vendorPreResult[:vendorMaxLen]
		} else {
			vendorResult = vendorPreResult
		}
	}

	vendorResult = strings.Trim(vendorResult, "-/&+")
	vendorResult = strings.Replace(vendorResult, "/", "-", -1)
	vendorResult = strings.Replace(vendorResult, ".", "", -1)
	return vendorResult
}

// Get usb-oriented human-readable device name
func getReadableDevice(s string) string {
	deviceSplitList := strings.Fields(s)
	if len(deviceSplitList) == 0 {
		return ""
	}

	resultsDevice := ""
	for _, device := range deviceSplitList {
		device = strings.Trim(device, "-/&*+[]()~.,#®\"")
		if len(device) != 0 {
			resultsDevice += device + "-"
		}
	}

	if len(resultsDevice) > deviceMaxLen {
		resultsDevice = resultsDevice[:deviceMaxLen-1]
	}

	resultsDevice = strings.Replace(resultsDevice, "/", "-", -1)
	resultsDevice = strings.Replace(resultsDevice, "&", "-", -1)
	resultsDevice = strings.Replace(resultsDevice, "*", "-", -1)
	resultsDevice = strings.Replace(resultsDevice, "+", "-", -1)
	resultsDevice = strings.Replace(resultsDevice, "^", "-", -1)
	resultsDevice = strings.Replace(resultsDevice, "\\", "-", -1)
	resultsDevice = strings.Replace(resultsDevice, "(", "-", -1)
	resultsDevice = strings.Replace(resultsDevice, ")", "-", -1)
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
