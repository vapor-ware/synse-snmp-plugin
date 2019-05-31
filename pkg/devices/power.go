package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/outputs"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// SnmpPower is the handler for SNMP OIDs that report power.
var SnmpPower = sdk.DeviceHandler{
	Name: "power",
	Read: SnmpPowerRead,
}

// SnmpPowerRead is the read handler function for synse SNMP devices that report power.
func SnmpPowerRead(device *sdk.Device) (readings []*output.Reading, err error) {

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
	result, err := snmpClient.Get(fmt.Sprint(data["oid"]))
	if err != nil {
		return nil, err
	}

	// Account for a multiplier if any and convert to float.
	var resultFloat float32
	resultFloat, err = MultiplyReading(result, data)
	if err != nil {
		return nil, err
	}

	// Create the reading.
	// FIXME (etd): differentiate between watts/VA
	reading := outputs.WattsPower.MakeReading(resultFloat)

	readings = []*output.Reading{reading}
	return readings, nil
}
