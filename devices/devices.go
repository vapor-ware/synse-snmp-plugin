package devices

// This file contains device utility functions.

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
)

// FindDeviceConfigsByType returns all elements in a DeviceConfig array where
// the Type is t.
// TODO: Could some of these be SDK helper functions? Maybe?
func FindDeviceConfigsByType(devices []*config.DeviceConfig, t string) (
	matches []*config.DeviceConfig, err error) {
	if devices == nil {
		return nil, fmt.Errorf("devices is nil")
	}

	for i := 0; i < len(devices); i++ {
		if devices[i].Type == t {
			matches = append(matches, devices[i])
		}
	}
	return matches, err
}

// FindDeviceConfigsByModel returns all elements in a DeviceConfig array where
// the Model is model.
func FindDeviceConfigsByModel(devices []*config.DeviceConfig, model string) (
	matches []*config.DeviceConfig, err error) {
	if devices == nil {
		return nil, fmt.Errorf("devices is nil")
	}

	for i := 0; i < len(devices); i++ {
		if devices[i].Model == model {
			matches = append(matches, devices[i])
		}
	}
	return matches, err
}

// DumpDeviceConfigs utility function dumps a slice of DeviceConfig to the
// console with a header.
func DumpDeviceConfigs(devices []*config.DeviceConfig, header string) {
	fmt.Printf("Dumping Devices. ")
	fmt.Print(header)

	if devices == nil {
		fmt.Printf(" <nil>\n")
		return
	}

	fmt.Printf(". Count: %d\n", len(devices))

	for i := 0; i < len(devices); i++ {
		fmt.Printf("device[%d]: %v %v %v %v %v row:%v column:%v\n", i,
			devices[i].Data["table_name"],
			devices[i].Type,
			devices[i].Data["info"],
			devices[i].Data["oid"],
			devices[i].Data["base_oid"],
			devices[i].Data["row"],
			devices[i].Data["column"])
	}
}
