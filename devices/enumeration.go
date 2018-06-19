package devices

import (
	"fmt"

	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// IsEnumeration returns whether or not the SNMP device reading represents an
// enumeration for a sysne device.
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
func TranslateEnumeration(result core.ReadResult, data map[string]string) (translation string, err error) {
	// Raw SNMP reading should be an int.
	resultInt, ok := result.Data.(int)
	if !ok {
		return "", fmt.Errorf(
			"Expected int reading, got type: %T, value: %v",
			result.Data, result.Data)
	}

	key := fmt.Sprintf("enumeration%d", resultInt)
	translation, ok = data[key]
	if !ok {
		translation = "undefined" // No translation found. unknown is acually used in SNMP.
	}
	return translation, nil
}
