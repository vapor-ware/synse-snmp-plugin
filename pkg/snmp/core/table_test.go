package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTable
// Initial test creates a table based on the UPS-MIB, upsInput table.
func TestTable(t *testing.T) {

	// In order to create the table, we need to create an SNMP Server.
	// In order to create the SNMP server, we need to have an SnmpClient.

	// Create SecurityParameters for the config that should connect to the emulator.
	securityParameters, err := NewSecurityParameters(
		"simulator",  // User Name
		SHA,          // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		AES,          // Privacy Protocol
		"privatus",   // Privacy Passphrase
	)
	assert.NoError(t, err)

	// Create a config.
	config, err := NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public",   //  Context name
		[]string{}, // tags (none)
	)
	assert.NoError(t, err)

	// Create a client.
	client, err := NewSnmpClient(config)
	assert.NoError(t, err)

	// Create SnmpServerBase
	snmpServer, err := NewSnmpServerBase(client, config)
	assert.NoError(t, err)

	// Create SnmpTable similar to the table for the UPS input power.
	// The table here has an empty DeviceEnumerator.
	testUpsInputTable, err := NewSnmpTable(
		"fakeTestUpsInputTable", // Table name. Same as OID .1.3.6.1.2.1.33.1.3.3 (Walk OID)
		".1.3.6.1.2.1.33.1.3.3", // Walk OID
		[]string{ // Column names
			"upsInputLineIndex",
			"upsInputFrequency",
			"upsInputVoltage",
			"upsInputCurrent",
			"upsInputTruePower",
		},
		snmpServer, // server
		"1",        // rowBase
		"",         // indexColumn
		"2",        // readableColumn
		false,      // flattened table
	)
	assert.NoError(t, err)

	testUpsInputTable.Dump()

	// Call DeviceEnumerator for testUpsInputTable.
	// It is currently the default which does nothing, unlike the real table defined under ups_mib.
	devices, err := testUpsInputTable.DevEnumerator.DeviceEnumerator(nil)
	assert.NoError(t, err)
	assert.NotNil(t, devices)
	assert.Len(t, devices, 0)
}
