package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpVoltage is the handler for the SNMP OIDs that report voltage.
var SnmpVoltage = sdk.DeviceHandler{
	Type:  "voltage",
	Model: "PXGMS UPS + EATON 93PM",

	Read:     SnmpVoltageRead,
	Write:    nil, // NYI for V1
	BulkRead: nil,
}

// SnmpVoltageRead is the read handler function for synse SNMP devices that report voltage.
func SnmpVoltageRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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

	// Account for a multiplier if any and convert to float.
	var resultFloat float32
	resultFloat, err = MultiplyReading(result, data)
	if err != nil {
		return nil, err
	}
	resultString := fmt.Sprintf("%.1f", resultFloat)

	// Create the reading.
	readings = []*sdk.Reading{
		sdk.NewReading("voltage", resultString),
	}
	return readings, nil
}
