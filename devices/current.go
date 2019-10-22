package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
)

// SnmpCurrent is the handler for the SNMP OIDs that report current.
var SnmpCurrent = sdk.DeviceHandler{
	Name: "current",
	Read: SnmpCurrentRead,
}

// SnmpCurrentRead is the read handler function for synse SNNP devices that report current.
func SnmpCurrentRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Account for a multiplier if any and convert to float.
	var resultFloat float32
	resultFloat, err = MultiplyReading(result, device.Data)
	if err != nil {
		return nil, err
	}

	// Create the reading.
	reading, err := device.GetOutput("current").MakeReading(resultFloat)
	if err != nil {
		return nil, err
	}

	readings = []*sdk.Reading{reading}
	return readings, nil
}
