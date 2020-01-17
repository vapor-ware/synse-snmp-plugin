package devices

// This file contains device utility functions and common device functions.

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
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

// Get the raw reading from the SNMP server with error checks.
// Factors out common code.
func getRawReading(device *sdk.Device) (result core.ReadResult, err error) {
	// Arg checks.
	if device == nil {
		return result, fmt.Errorf("device is nil")
	}

	// Get the SNMP device config from the strings in data.
	data := device.Data
	snmpConfig, err := core.GetDeviceConfig(data)
	if err != nil {
		return result, err
	}

	// Create SnmpClient.
	snmpClient, err := core.NewSnmpClient(snmpConfig)
	if err != nil {
		return result, err
	}

	// Read the SNMP OID in the device config.
	return snmpClient.Get(fmt.Sprint(data["oid"]))
}
