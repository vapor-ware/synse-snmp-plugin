package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// SnmpPower is the handler for SNMP OIDs that report power.
var SnmpPower = sdk.DeviceHandler{
	Name: "power",
	Read: SnmpPowerRead,
}

// SnmpPowerRead is the read handler function for synse SNMP devices that report power.
func SnmpPowerRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Get the raw reading from the SNMP server.
	var result core.ReadResult
	result, err = getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Check for nil reading.
	var reading *output.Reading
	if result.Data == nil {
		reading, err = output.Watt.MakeReading(nil)
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
	// FIXME (etd): differentiate between watts/VA
	reading, err = output.Watt.MakeReading(resultFloat)

	readings = []*output.Reading{reading}
	return readings, nil
}
