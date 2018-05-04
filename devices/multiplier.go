package devices

import (
	"fmt"
	"strconv"

	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// MultiplyReading is a helper method to multiply out raw readings
// appropriately.
func MultiplyReading(result core.ReadResult, data map[string]string) (resultFloat float32, err error) {

	// Raw SNMP reading should be an int.
	resultInt, ok := result.Data.(int)
	if !ok {
		return 0.0, fmt.Errorf(
			"Expected int reading, got type: %T, value: %v",
			result.Data, result.Data)
	}

	// Account for a multiplier if any, otherwise just convert to float32.
	multiplierString, ok := data["multiplier"]
	if ok {
		var multiplierFloat float64
		multiplierFloat, err = strconv.ParseFloat(multiplierString, 32)
		if err != nil {
			return 0.0, err
		}
		resultFloat = float32(resultInt) * float32(multiplierFloat)
	} else {
		resultFloat = float32(resultInt)
	}

	return resultFloat, err
}
