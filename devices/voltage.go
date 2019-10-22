package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
)

// SnmpVoltage is the handler for the SNMP OIDs that report voltage.
var SnmpVoltage = sdk.DeviceHandler{
	Name: "voltage",
	Read: SnmpVoltageRead,
}

// SnmpVoltageRead is the read handler function for synse SNMP devices that report voltage.
func SnmpVoltageRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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
	reading, err := device.GetOutput("voltage").MakeReading(resultFloat)
	if err != nil {
		return nil, err
	}

	readings = []*sdk.Reading{reading}
	return readings, nil
}
