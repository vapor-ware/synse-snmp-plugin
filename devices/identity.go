package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpIdentity is the handler for the snmp-identity device.
var SnmpIdentity = sdk.DeviceHandler{
	Type:  "identity",
	Model: "PXGMS UPS + EATON 93PM",

	Read:     SnmpIdentityRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpIdentityRead is the read handler function for snmp-identity devices.
func SnmpIdentityRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

	// Arg checks.
	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	// Get the SNMP device config from the strings in data.
	data := device.Data
	snmpConfig, err := core.GetDeviceConfig(data)
	if err != nil {
		return nil, err
	}

	// Create SnmpClient.
	snmpClient, err := core.NewSnmpClient(snmpConfig)
	if err != nil {
		return nil, err
	}

	// Read the SNMP OID in the device config.
	result, err := snmpClient.Get(data["oid"])
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
	readings = []*sdk.Reading{
		sdk.NewReading("identity", resultString),
	}
	return readings, nil
}
