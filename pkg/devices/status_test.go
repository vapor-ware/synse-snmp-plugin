package devices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// TestUpsBatteryStatus tests the enumeration reading value for upsBatteryStatus.
// Regression test for https://vaporio.atlassian.net/browse/VIO-1092
func TestUpsBatteryStatus(t *testing.T) {

	// Create a status device that will give us a good reading value.
	device := sdk.Device{
		Data: map[string]interface{}{
			"contextName":              "public",
			"endpoint":                 "127.0.0.1",
			"userName":                 "simulator",
			"privacyProtocol":          "AES",
			"privacyPassphrase":        "privatus",
			"port":                     1024,
			"authenticationProtocol":   "SHA",
			"authenticationPassphrase": "auctoritas",
			"model":                    "Galaxy VM 180 kVA",
			"version":                  "v3",

			////
			"base_oid":   ".1.3.6.1.2.1.33.1.2.%d.0",
			"column":     1,
			"oid":        ".1.3.6.1.2.1.33.1.2.1.0",
			"row":        0,
			"table_name": "UPS-MIB-UPS-Battery-Table",
			// This is an enumeration. We need to translate the integer we read to a string.
			"enumeration":  "true", // Defines that this is an enumeration.
			"enumeration1": "unknown",
			"enumeration2": "batteryNormal",
			"enumeration3": "batteryLow",
			"enumeration4": "batteryDepleted",
		},
	}

	// Call SnmpStatusRead.
	readings, err := SnmpStatusRead(&device)

	// Verify we get a single good reading and no error.
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Equal(t, "batteryNormal", readings[0].Value)
}
