package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// SnmpStatusInt is the handler for the SNMP status-int devices.
var SnmpStatusInt = sdk.DeviceHandler{
	Name: "status-int",
	Read: SnmpStatusIntRead,
}

// SnmpStatusUint is the handler for the SNMP status-uint devices.
var SnmpStatusUint = sdk.DeviceHandler{
	Name: "status-uint",
	Read: SnmpStatusUintRead,
}

// SnmpStatusString is the handler for the SNMP status-string devices.
var SnmpStatusString = sdk.DeviceHandler{
	Name: "status-string",
	Read: SnmpStatusStringRead,
}

// SnmpStatusIntRead is the read handler function for SNMP status-int devices.
func SnmpStatusIntRead(device *sdk.Device) (readings []*sdk.Reading, err error) { // nolint: gocyclo

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Should be an int.
	var reading *sdk.Reading
	if result.Data != nil {
		var resultInt int
		resultInt, ok := result.Data.(int)
		if !ok {
			return nil, fmt.Errorf(
				"Expected int status reading, got type: %T, value: %v",
				result.Data, result.Data)
		}
		// Create the reading.
		reading, err = device.GetOutput("status-int").MakeReading(resultInt)
		if err != nil {
			return nil, err
		}
	} else {
		// Create the reading.
		reading, err = device.GetOutput("status-int").MakeReading(nil)
		if err != nil {
			return nil, err
		}
	}

	readings = []*sdk.Reading{reading}
	return readings, nil
}

// SnmpStatusUintRead is the read handler function for SNMP status-int devices.
func SnmpStatusUintRead(device *sdk.Device) (readings []*sdk.Reading, err error) { // nolint: gocyclo

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Should be an int.
	var reading *sdk.Reading
	if result.Data != nil {
		var resultUint uint
		resultUint, ok := result.Data.(uint)
		if !ok {
			return nil, fmt.Errorf(
				"Expected int status reading, got type: %T, value: %v",
				result.Data, result.Data)
		}
		// Create the reading.
		reading, err = device.GetOutput("status-uint").MakeReading(resultUint)
		if err != nil {
			return nil, err
		}
	} else {
		// Create the reading.
		reading, err = device.GetOutput("status-uint").MakeReading(nil)
		if err != nil {
			return nil, err
		}
	}

	readings = []*sdk.Reading{reading}
	return readings, nil
}

// SnmpStatusStringRead is the read handler function for SNMP status-string devices.
func SnmpStatusStringRead(device *sdk.Device) (readings []*sdk.Reading, err error) { // nolint: gocyclo

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Should be a string.
	resultString := "" // Default reading for nil. TODO: Ensure correctness here. Make sure we have a test for this.
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
			// An Int has to be be an enumeration.
			// TODO: Logic around here probably needs work.
			if IsEnumeration(device.Data) {
				resultString, err = TranslateEnumeration(result, device.Data)
				if err != nil {
					return nil, err
				}
			} else {
				resultString = fmt.Sprintf("%d", resultInt)
			}
		}
	}
	// Create the reading.
	reading, err := device.GetOutput("status-string").MakeReading(resultString)
	if err != nil {
		return nil, err
	}
	readings = []*sdk.Reading{reading}
	return readings, nil
}
