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
	var result core.ReadResult
	result, err = getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Should be a string.
	var reading *output.Reading
	resultData, ok := result.Data.(string)
	if !ok {
		if result.Data == nil {
			// We got nil, so create a nil reading.
			reading, err = outputs.Identity.MakeReading(nil)
			if err != nil {
				return nil, err
			}
			readings = []*output.Reading{reading}
			return readings, nil
		}
		return nil, fmt.Errorf(
			"expected string identity reading, got type: %T, value: %v",
			result.Data, result.Data)
	}

	// Create the reading.
	reading, err = outputs.Identity.MakeReading(resultData)
	if err != nil {
		return nil, err
	}
	readings = []*output.Reading{reading}
	return readings, nil
}
