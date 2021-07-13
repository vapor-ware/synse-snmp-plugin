package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// SnmpSeconds is the handler for the SNMP OIDs that report seconds.
var SnmpSeconds = sdk.DeviceHandler{
	Name: "seconds",
	Read: SnmpSecondsRead,
}

// SnmpSecondsRead is the read handler function for Synse SNMP devices that report seconds.
func SnmpSecondsRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Check for nil reading.
	if result.Data == nil {
		reading := output.Seconds.MakeReading(nil)
		readings = []*output.Reading{reading}
		return
	}

	// Create the reading.
	reading := output.Seconds.MakeReading(result.Data)
	readings = []*output.Reading{reading}
	return
}
