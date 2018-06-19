package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// SnmpStatus is the handler for the snmp-status device.
var SnmpStatus = sdk.DeviceHandler{
	Name: "status",
	Read: SnmpStatusRead,
}

// SnmpStatusRead is the read handler function for snmp-status devices.
func SnmpStatusRead(device *sdk.Device) (readings []*sdk.Reading, err error) {

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

	// Should be a string.
	resultString := "" // Default reading for nil.
	if result.Data != nil {
		var ok bool
		resultString, ok = result.Data.(string)
		if !ok {
			// Could be an int as well.
			var resultInt int
			resultInt, ok = result.Data.(int)
			if !ok {
				return nil, fmt.Errorf(
					"Expected string or int status reading, got type: %T, value: %v",
					result.Data, result.Data)
			}
			// An Int could be an enumeration.
			if IsEnumeration(data) {
				resultString, err = TranslateEnumeration(result, data)
				if err != nil {
					return nil, err
				}
			} else {
				resultString = fmt.Sprintf("%d", resultInt)
			}
		}
	}
	// Create the reading.
	readings = []*sdk.Reading{
		device.GetOutput("status").MakeReading(resultString),
	}
	return readings, nil
}
