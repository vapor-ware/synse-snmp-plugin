package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// SnmpPercentage is the handler for the SNMP OIDs that report percentage.
var SnmpPercentage = sdk.DeviceHandler{
	Name: "percentage",
	Read: SnmpPercentageRead,
}

// SnmpPercentageRead is the read handler function for Synse SNMP devices that report percentage.
func SnmpPercentageRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Get the raw reading from the SNMP server.
	var result core.ReadResult
	result, err = getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Check for nil reading.
	var reading *output.Reading
	if result.Data == nil {
		reading, err = output.Percentage.MakeReading(nil)
		if err != nil {
			return nil, err
		}
		readings = []*output.Reading{reading}
		return
	}

	// Create the reading.
	reading, err = output.Percentage.MakeReading(result.Data)
	if err != nil {
		return nil, err
	}
	readings = []*output.Reading{reading}
	return
}
