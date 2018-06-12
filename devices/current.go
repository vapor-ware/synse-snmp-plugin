package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpCurrent is the handler for the SNMP OIDs that report current.
var SnmpCurrent = sdk.DeviceHandler{
	Type:  "current",
	Model: "PXGMS UPS + EATON 93PM", // TODO: This needs to be here to match the device handler. Want it more generic.

	Read:     SnmpCurrentRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpCurrentRead is the read handler function for synse SNNP devices that report current.
func SnmpCurrentRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

	// Arg checks.
	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	// Get the SNMP device config from the strings in the data.
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

	// Account for a multiplier if any and convert to float.
	var resultFloat float32
	resultFloat, err = MultiplyReading(result, data)
	if err != nil {
		return nil, err
	}
	resultString := fmt.Sprintf("%.1f", resultFloat)

	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("current", resultString),
	}
	return readings, nil
}
