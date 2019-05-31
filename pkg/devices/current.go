package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// SnmpCurrent is the handler for the SNMP OIDs that report current.
var SnmpCurrent = sdk.DeviceHandler{
	Name: "current",
	Read: SnmpCurrentRead,
}

// SnmpCurrentRead is the read handler function for Synse SNMP devices that report current.
func SnmpCurrentRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Arg checks.
	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	// Get the SNMP device config from the strings in the data.
	data := device.Data
	snmpConfig, err := core.GetDeviceConfig(data)
	if err != nil {
		return nil, err
	}

	// Create SnmpClient.
	snmpClient, err := core.NewSnmpClient(snmpConfig)
	if err != nil {
		return nil, err
	}

	// Read the SNMP OID in the device config.
	result, err := snmpClient.Get(fmt.Sprint(data["oid"]))
	if err != nil {
		return nil, err
	}

	// Account for a multiplier if any and convert to float.
	var resultFloat float32
	resultFloat, err = MultiplyReading(result, data)
	if err != nil {
		return nil, err
	}

	// Create the reading.
	reading := output.ElectricCurrent.MakeReading(resultFloat)

	readings = []*output.Reading{reading}
	return readings, nil
}
