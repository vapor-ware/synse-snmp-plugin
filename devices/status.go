package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
)

// SnmpStatus is the handler for the snmp-status device.
var SnmpStatus = sdk.DeviceHandler{
	Name: "status",
	Read: SnmpStatusRead,
}

// SnmpStatusRead is the read handler function for snmp-status devices.
func SnmpStatusRead(device *sdk.Device) (readings []*sdk.Reading, err error) { // nolint: gocyclo

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
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
	reading, err := device.GetOutput("status").MakeReading(resultString)
	if err != nil {
		return nil, err
	}
	readings = []*sdk.Reading{reading}
	return readings, nil
}
