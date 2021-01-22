package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/outputs"
)

// SnmpIdentity is the handler for the snmp-identity device.
var SnmpIdentity = sdk.DeviceHandler{
	Name: "identity",
	Read: SnmpIdentityRead,
}

// SnmpIdentityRead is the read handler function for snmp-identity devices.
func SnmpIdentityRead(device *sdk.Device) (readings []*output.Reading, err error) {

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Should be a string.
	resultData, ok := result.Data.(string)
	if !ok {
		if result.Data == nil {
			// We got nil, so create a nil reading.
			reading := outputs.Identity.MakeReading(nil)
			readings = []*output.Reading{reading}
			return readings, nil
		}
		return nil, fmt.Errorf(
			"expected string identity reading, got type: %T, value: %v",
			result.Data, result.Data)
	}

	// Create the reading.
	reading := outputs.Identity.MakeReading(resultData)
	readings = []*output.Reading{reading}
	return readings, nil
}
