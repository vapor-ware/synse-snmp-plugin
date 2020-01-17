package devices

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// SnmpStatus is the handler for the snmp-status device.
var SnmpStatus = sdk.DeviceHandler{
	Name: "status",
	Read: SnmpStatusRead,
}

// SnmpStatusRead is the read handler function for snmp-status devices.
func SnmpStatusRead(device *sdk.Device) (readings []*output.Reading, err error) { // nolint: gocyclo

	// Get the raw reading from the SNMP server.
	result, err := getRawReading(device)
	if err != nil {
		return nil, err
	}

	// Generally, the value we get back as the reading should be the value we return.
	// It could be a string, int, uint, etc. The "status" output does not place any
	// kind of restrictions on its returned value.
	//
	// The value we get back may need additional processing if the device reading is
	// an enumeration. Here, we check if the device data has an enumeration defined,
	// and if so, translate the enumeration. Otherwise, we return whatever the raw
	// reading is.
	var value interface{}
	if result.Data != nil {
		if IsEnumeration(device.Data) {
			value, err = TranslateEnumeration(result, device.Data)
			if err != nil {
				return nil, err
			}
		} else {
			value = result.Data
		}
	}

	return []*output.Reading{
		output.Status.MakeReading(value),
	}, nil
}
