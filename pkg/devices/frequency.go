package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// SnmpFrequency is the handler for the SNMP OIDs that report frequency.
var SnmpFrequency = sdk.DeviceHandler{
	Name: "frequency",
	Read: SnmpFrequencyRead,
}

// SnmpFrequencyRead is the read handler function for synse SNMP devices that report frequency.
func SnmpFrequencyRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Check for nil reading.
	if result.Data == nil {
		reading := output.Frequency.MakeReading(nil)
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
	reading := output.Frequency.MakeReading(resultFloat)

	readings = []*output.Reading{reading}
	return readings, nil
}
