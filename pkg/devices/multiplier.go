package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// MultiplyReading is a helper method to multiply out raw readings
// appropriately.
func MultiplyReading(result core.ReadResult, data map[string]interface{}) (resultFloat float32, err error) {

	// Raw SNMP reading should be an int.
	resultInt, ok := result.Data.(int)
	if !ok {
		return 0.0, fmt.Errorf(
			"expected int reading, got type: %T, value: %v",
			result.Data, result.Data)
	}

	// Account for a multiplier if any, otherwise just convert to float32.
	multiplier, ok := data["multiplier"]
	if ok {
		multiplierFloat, isOk := multiplier.(float32)
		if !isOk {
			return 0.0, fmt.Errorf(
				"expected float multiplier, got type: %T, value: %v", multiplier, multiplier,
			)
		}
		resultFloat = float32(resultInt) * multiplierFloat
	} else {
		resultFloat = float32(resultInt)
	}

	return resultFloat, err
}
