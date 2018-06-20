package devices

// This file contains device utility functions.

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
)

//// FindDeviceConfigsByType returns all elements in a DeviceConfig array where
//// the Type is t.
//// TODO: Could some of these be SDK helper functions? Maybe?
//// FIXME: this is only used in tests - move to a test utility?
//func FindDeviceConfigsByType(devices []*sdk.DeviceConfig, t string) (
//	matches []*sdk.DeviceConfig, err error) {
//	if devices == nil {
//		return nil, fmt.Errorf("devices is nil")
//	}
//
//	for _, device := range devices {
//		for _, kind := range device.Devices {
//			if kind.Name == t {
//				matches = append(matches, device)
//			}
//		}
//	}
//	return matches, err
//}

// DumpDeviceConfigs utility function dumps a slice of DeviceConfig to the
// console with a header.
func DumpDeviceConfigs(devices []*sdk.DeviceConfig, header string) {
	fmt.Printf("Dumping Devices. ")
	fmt.Print(header)

	if devices == nil {
		fmt.Printf(" <nil>\n")
		return
	}

	fmt.Printf(". Count device config: %d\n", len(devices))

	for _, device := range devices {
		fmt.Printf(".. Count device kind: %d\n", len(device.Devices))
		for _, kind := range device.Devices {
			fmt.Printf("... Count device instances: %d\n", len(kind.Instances))
			for _, instance := range kind.Instances {
				fmt.Printf("device: %v %v %v %v %v row:%v column:%v\n",
					instance.Data["table_name"],
					kind.Name,
					instance.Data["info"],
					instance.Data["oid"],
					instance.Data["base_oid"],
					instance.Data["row"],
					instance.Data["column"],
				)
			}
		}
	}
}
