package core

import (
	"testing"
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
		"privatus")   // Privacy Passphrase
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Create a config.
	config, err := NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public") //  Context name
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Create a client.
	client, err := NewSnmpClient(config)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Create SnmpServerBase
	snmpServer, err := NewSnmpServerBase(
		client,
		config)
	if err != nil {
		t.Error(err) // Fail the test.
	}

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
		false)      // flattened table
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	testUpsInputTable.Dump()

	// Call DeviceEnumerator for testUpsInputTable.
	// It is currently the default which does nothing, unlike the real table defined under ups_mib.
	devices, err := testUpsInputTable.DevEnumerator.DeviceEnumerator(nil)
	if devices == nil {
		t.Fatal("devices is nil")
	}
	if len(devices) != 0 {
		t.Fatalf("Should have zero devices enumerated, got %d", len(devices))
	}
	if err != nil {
		t.Fatal(err)
	}
}
