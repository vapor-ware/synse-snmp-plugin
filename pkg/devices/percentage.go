package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// SnmpPercentage is the handler for the SNMP OIDs that report percentage.
var SnmpPercentage = sdk.DeviceHandler{
	Name: "percentage",
	Read: SnmpPercentageRead,
}

// SnmpPercentageRead is the read handler function for Synse SNMP devices that report percentage.
func SnmpPercentageRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Check for nil reading.
	if result.Data == nil {
		reading := output.Percentage.MakeReading(nil)
		readings = []*output.Reading{reading}
		return
	}

	// Create the reading.
	reading := output.Percentage.MakeReading(result.Data)
	readings = []*output.Reading{reading}
	return
}
