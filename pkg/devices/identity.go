package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/outputs"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
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
	resultString, ok := result.Data.(string)
	if !ok {
		return nil, fmt.Errorf(
			"Expected int identity reading, got type: %T, value: %v",
			result.Data, result.Data)
	}

	// Create the reading.
	reading := outputs.Identity.MakeReading(resultString)

	readings = []*output.Reading{reading}
	return readings, nil
}
