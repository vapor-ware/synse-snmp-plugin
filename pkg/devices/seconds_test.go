package devices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// TestUpsSecondsOnBattery tests the enumeration reading value for upsSecondsOnBattery.
// Regression test for https://vaporio.atlassian.net/browse/VIO-1096
func TestUpsSecondsOnBattery(t *testing.T) {

	// Create a seconds device that will give us a good reading value.
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
			"column":     2,
			"oid":        ".1.3.6.1.2.1.33.1.2.2.0",
			"row":        0,
			"table_name": "UPS-MIB-UPS-Battery-Table",
		},
	}

	// Call SnmpSecondsRead.
	readings, err := SnmpSecondsRead(&device)

	// Verify we get a single good reading and no error.
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Equal(t, 0, readings[0].Value)
	assert.Equal(t, "seconds", readings[0].Unit.Name)
	assert.Equal(t, "s", readings[0].Unit.Symbol)
}
