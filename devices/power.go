package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
)

// SnmpPower is the handler for SNMP OIDs that report power.
var SnmpPower = sdk.DeviceHandler{
	Name: "power",
	Read: SnmpPowerRead,
}

// SnmpPowerRead is the read handler function for synse SNMP devices that report power.
func SnmpPowerRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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
	// FIXME (etd): differentiate between watts/VA
	reading, err := device.GetOutput("watts.power").MakeReading(resultFloat)
	if err != nil {
		return nil, err
	}

	readings = []*sdk.Reading{reading}
	return readings, nil
}
