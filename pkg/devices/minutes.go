package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// SnmpMinutes is the handler for the SNMP OIDs that report minutes.
var SnmpMinutes = sdk.DeviceHandler{
	Name: "minutes",
	Read: SnmpMinutesRead,
}

// SnmpMinutesRead is the read handler function for Synse SNMP devices that report minutes.
func SnmpMinutesRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Check for nil reading.
	if result.Data == nil {
		reading := output.Minutes.MakeReading(nil)
		readings = []*output.Reading{reading}
		return
	}

	// Create the reading.
	reading := output.Minutes.MakeReading(result.Data)
	readings = []*output.Reading{reading}
	return
}
