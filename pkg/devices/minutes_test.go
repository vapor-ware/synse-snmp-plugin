package devices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// TestEstimatedMinutesRemaining tests the enumeration reading value for upsEstimatedMinutesRemaining.
// Regression test for https://vaporio.atlassian.net/browse/VIO-1094
func TestUpsEstimatedMinutesRemaining(t *testing.T) {

	// Create a minutes device that will give us a good reading value.
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
			"column":     3,
			"oid":        ".1.3.6.1.2.1.33.1.2.3.0",
			"row":        0,
			"table_name": "UPS-MIB-UPS-Battery-Table",
		},
	}

	// Call SnmpMinutesRead.
	readings, err := SnmpMinutesRead(&device)

	// Verify we get a single good reading and no error.
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Equal(t, 1092, readings[0].Value)
	assert.Equal(t, "minutes", readings[0].Unit.Name)
	assert.Equal(t, "min", readings[0].Unit.Symbol)
}
