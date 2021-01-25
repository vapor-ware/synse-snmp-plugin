package devices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// TestCurrentNilReadingValue tests that we get a nil reading value for a
// current reading where the OID is not served by the emulator. See
// https://github.com/vapor-ware/synse-snmp-plugin/issues/61
func TestCurrentNilReadingValue(t *testing.T) {
	// Create a current device that will give us a nil reading value.
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
			"base_oid":   ".1.3.6.1.2.1.32.1.2.%d.0", // ups mib is at .33, not .32
			"column":     6,
			"oid":        ".1.3.6.1.2.1.32.1.2.6.0", // ups mib at .33, not .32
			"row":        0,
			"table_name": "UPS-MIB-UPS-Battery-Table",
			"multiplier": float32(0.1),
		},
	}
	// Call SnmpCurrentRead.
	readings, err := SnmpCurrentRead(&device)
	// Verify we get a nil reading and no error.
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Nil(t, readings[0].Value)
}
