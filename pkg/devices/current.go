package devices

import (
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

	// Get the raw reading from the SNMP server.
	var result core.ReadResult
	result, err = getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Check for nil reading.
	var reading *output.Reading
	if result.Data == nil {
		reading, err = output.ElectricCurrent.MakeReading(nil)
		if err != nil {
			return nil, err
		}
		readings = []*output.Reading{reading}
		return readings, nil
	}

	// Account for a multiplier if any and convert to float.
	var resultFloat float32
	resultFloat, err = MultiplyReading(result, device.Data)
	if err != nil {
		return nil, err
	}

	// Create the reading.
	reading, err = output.ElectricCurrent.MakeReading(resultFloat)
	if err != nil {
		return nil, err
	}

	readings = []*output.Reading{reading}
	return readings, nil
}
