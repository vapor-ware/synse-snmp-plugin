package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// IsEnumeration returns whether or not the SNMP device reading represents an
// enumeration for a synse device.
// data is the map associated with a synse device.
func IsEnumeration(data map[string]interface{}) (yes bool) {
	// Find the enumeration key in the map and check that it is set to true.
	setting, ok := data["enumeration"]
	if !ok {
		return false
	}

	if setting == "true" {
		return true
	}
	return false
}

// TranslateEnumeration translates a read result to a string based on the
// enumeration. The caller should call IsEnumeration first for this translation
// to make sense.
func TranslateEnumeration(result core.ReadResult, data map[string]interface{}) (string, error) {
	// Raw SNMP reading should be an int.
	resultInt, ok := result.Data.(int)
	if !ok {
		if result.Data == nil {
			return "", nil // Nil reading data. Return empty string.
		}
		return "", fmt.Errorf(
			"expected int reading, got type: %T, value: %v",
			result.Data, result.Data)
	}

	// Key lookup to find the enumeration.
	key := fmt.Sprintf("enumeration%d", resultInt)
	translation, ok := data[key]
	if !ok {
		// Not found. Return something with the raw int.
		translation = fmt.Sprintf("undefined%d", resultInt)
	}
	return fmt.Sprint(translation), nil
}
