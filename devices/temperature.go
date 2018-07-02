package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpTemperature is the handler for the SNMP OIDs that report temperature.
var SnmpTemperature = sdk.DeviceHandler{
	Name: "temperature",
	Read: SnmpTemperatureRead,
}

// SnmpTemperatureRead is the read handler function for synse SNMP devices that report temperature.
func SnmpTemperatureRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

	// Arg checks.
	if device == nil {
		return nil, fmt.Errorf("device is nil")
	}

	// Get the device config from the strings in data.
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
	readings = []*sdk.Reading{
		device.GetOutput("temperature").MakeReading(resultFloat),
	}
	return readings, nil
}
