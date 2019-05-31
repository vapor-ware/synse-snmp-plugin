package devices

// This file contains device utility functions.

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
)

// DumpDeviceConfigs utility function dumps a slice of DeviceConfig to the
// console with a header.
func DumpDeviceConfigs(devices []*config.DeviceProto, header string) {
	fmt.Printf("Dumping Devices. ")
	fmt.Print(header)

	if devices == nil {
		fmt.Printf(" <nil>\n")
		return
	}

	fmt.Printf(". Count device config: %d\n", len(devices))

	fmt.Printf(".. Count device proto: %d\n", len(devices))
	for _, proto := range devices {

		fmt.Printf("... Count device instances: %d\n", len(proto.Instances))
		for _, instance := range proto.Instances {
			fmt.Printf("device: %v %v %v %v %v row:%v column:%v\n",
				instance.Data["table_name"],
				proto.Type,
				instance.Info,
				instance.Data["oid"],
				instance.Data["base_oid"],
				instance.Data["row"],
				instance.Data["column"],
			)
		}

	}
}
