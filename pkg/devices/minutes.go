package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// SnmpMinutes is the handler for the SNMP OIDs that report minutes.
var SnmpMinutes = sdk.DeviceHandler{
	Name: "minutes",
	Read: SnmpMinutesRead,
}

// SnmpMinutesRead is the read handler function for Synse SNMP devices that report minutes.
func SnmpMinutesRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Get the raw reading from the SNMP server.
	var result core.ReadResult
	result, err = getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Check for nil reading.
	var reading *output.Reading
	if result.Data == nil {
		reading, err = output.Minutes.MakeReading(nil)
		if err != nil {
			return nil, err
		}
		readings = []*output.Reading{reading}
		return
	}

	// Create the reading.
	reading, err = output.Minutes.MakeReading(result.Data)
	if err != nil {
		return nil, err
	}
	readings = []*output.Reading{reading}
	return
}
