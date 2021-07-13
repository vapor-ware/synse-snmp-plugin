package devices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// TestUpsEstimatedChargeRemaining tests the enumeration reading value for upsEstimatedChargeRemaining.
// Regression test for https://vaporio.atlassian.net/browse/VIO-1093
func TestUpsEstimatedChargeRemaining(t *testing.T) {

	// Create a percentage device that will give us a good reading value.
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
			"column":     4,
			"oid":        ".1.3.6.1.2.1.33.1.2.4.0",
			"row":        0,
			"table_name": "UPS-MIB-UPS-Battery-Table",
		},
	}

	// Call SnmpPercentageRead.
	readings, err := SnmpPercentageRead(&device)

	// Verify we get a single good reading and no error.
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Equal(t, 100, readings[0].Value)
	assert.Equal(t, "percent", readings[0].Unit.Name)
	assert.Equal(t, "%", readings[0].Unit.Symbol)
}

// TestUpsOutputPercentLoad0 tests the enumeration reading value for upsOutputPercentLoad0.
// Regression test for https://vaporio.atlassian.net/browse/VIO-1095
func TestUpsOutputPercentLoad0(t *testing.T) {

	// Create a percentage device that will give us a good reading value.
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
			"base_oid":   ".1.3.6.1.2.1.33.1.4.4.1.5.%d",
			"column":     1,
			"oid":        ".1.3.6.1.2.1.33.1.4.4.1.5.1",
			"row":        0,
			"table_name": "UPS-MIB-UPS-Output-Table",
		},
	}

	// Call SnmpPercentageRead.
	readings, err := SnmpPercentageRead(&device)

	// Verify we get a single good reading and no error.
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Equal(t, 0, readings[0].Value)
	assert.Equal(t, "percent", readings[0].Unit.Name)
	assert.Equal(t, "%", readings[0].Unit.Symbol)
}
