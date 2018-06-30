package mibs

import (
	"fmt"
	"testing"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

func FindDeviceInstanceByInfo(devices []*sdk.DeviceConfig, info string) (
	*sdk.DeviceConfig, *sdk.DeviceKind, *sdk.DeviceInstance) {
	for _, device := range devices {
		for _, kind := range device.Devices {
			for _, instance := range kind.Instances {
				// This only works since info is coded to be unique.
				if instance.Info == info {
					return device, kind, instance
				}
			}
		}
	}
	return nil, nil, nil
}

// TestUpsMib
// Initial test creates all tables based on the UPS-MIB.
func TestUpsMib(t *testing.T) { // nolint: gocyclo
	// In order to create the table, we need to create an SNMP Server.
	// In order to create the SNMP server, we need to have an SnmpClient.

	// Create SecurityParameters for the config that should connect to the emulator.
	securityParameters, err := core.NewSecurityParameters(
		"simulator",  // User Name
		core.SHA,     // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		core.AES,     // Privacy Protocol
		"privatus")   // Privacy Passphrase
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Create a config.
	config, err := core.NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public") //  Context name
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Create a client.
	client, err := core.NewSnmpClient(config)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Create SnmpServerBase
	snmpServer, err := core.NewSnmpServerBase(
		client,
		config)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Create the UpsMib and dump it.
	testUpsMib, err := NewUpsMib(snmpServer)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}
	testUpsMib.Dump()

	// We should have 19 tables.
	tableCount := len(testUpsMib.Tables)
	if tableCount != 19 {
		t.Fatalf("testUpsMib: Expected 19 tables, got %d", tableCount)
	}

	// Get the ups identity data from the test MIB.
	upsIdentity := testUpsMib.UpsIdentityTable.UpsIdentity

	// Verify expected ups identity data from the test MIB.
	if upsIdentity == nil {
		t.Fatal("upsIdentity is nil")
	}

	if upsIdentity.Manufacturer != "Eaton Corporation" {
		t.Fatalf("Expected upsIdentity.Manufacturer [Eaton Corporation], got [%v]",
			upsIdentity.Manufacturer)
	}

	if upsIdentity.Model != "PXGMS UPS + EATON 93PM" {
		t.Fatalf("Expected upsIdentity.Model [PXGMS UPS + EATON 93PM], got [%v]",
			upsIdentity.Model)
	}

	if upsIdentity.UpsSoftwareVersion != "INV: 1.44.0000" {
		t.Fatalf("Expected upsIdentity.UpsSoftwareVersion [INV: 1.44.0000], got [%v]",
			upsIdentity.UpsSoftwareVersion)
	}

	if upsIdentity.AgentSoftwareVersion != "2.3.7" {
		t.Fatalf("Expected upsIdentity.AgentSoftwareVersion [2.3.7], got [%v]",
			upsIdentity.AgentSoftwareVersion)
	}

	if upsIdentity.Name != "ID: EM111UXX06, Msg: 9PL15N0000E40R2" {
		t.Fatalf("Expected upsIdentity.Name [ID: EM111UXX06, Msg: 9PL15N0000E40R2], got [%v]",
			upsIdentity.Name)
	}

	if upsIdentity.AttachedDevices != "Attached Devices not set" {
		t.Fatalf("Expected upsIdentity.AttachedDevices [Attached Devices not set], got [%v]",
			upsIdentity.AttachedDevices)
	}

	// Call the ups battery table device enumerator.
	upsBatteryTable := testUpsMib.UpsBatteryTable
	devices, err := upsBatteryTable.SnmpTable.DevEnumerator.DeviceEnumerator(
		map[string]interface{}{"rack": "my_pet_rack", "board": "my_pet_board"})

	// Ensure devices and no error.
	if err != nil {
		t.Fatal(err)
	}
	// Check the number of device instances that were created
	instanceCount := 0
	for _, cfg := range devices {
		for _, kind := range cfg.Devices {
			instanceCount += len(kind.Instances)
		}
	}
	if instanceCount != 7 {
		t.Fatalf("Expected 7 devices from the UpsBatteryTable, got %d", instanceCount)
	}

	// Enumerate UpsInputTable devices.
	upsInputTable := testUpsMib.UpsInputTable
	devices, err = upsInputTable.SnmpTable.DevEnumerator.DeviceEnumerator(
		map[string]interface{}{"rack": "my_pet_rack", "board": "my_pet_board"})
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		t.Fatalf("Expected devices, got none.\n")
	}

	// Check the number of device instances that were created
	instanceCount = 0
	for _, cfg := range devices {
		for _, kind := range cfg.Devices {
			instanceCount += len(kind.Instances)
		}
	}
	if instanceCount != 12 {
		t.Fatalf("Expected 12 devices from the UpsInputTable, got %d", instanceCount)
	}

	// Enumerate the mib.
	// Testing for bad parameters is in TestDevices.
	devices, err = testUpsMib.EnumerateDevices(
		map[string]interface{}{"rack": "test_rack", "board": "test_board"})
	if err != nil {
		t.Fatal(err)
	}
	// Check the number of device instances that were created
	instanceCount = 0
	for _, cfg := range devices {
		for _, kind := range cfg.Devices {
			instanceCount += len(kind.Instances)
		}
	}
	if instanceCount != 40 {
		t.Fatalf("Expected 40 devices, got %d", instanceCount)
	}

	fmt.Printf("Dumping devices enumerated from UPS-MIB\n")
	for _, device := range devices {
		for _, kind := range device.Devices {
			for _, instance := range kind.Instances {
				fmt.Printf("UPS-MIB device: %v %v %v %v %v row:%v column:%v\n",
					instance.Data["table_name"],
					kind.Name,
					instance.Info,
					instance.Data["oid"],
					instance.Data["base_oid"],
					instance.Data["row"],
					instance.Data["column"],
				)
			}
		}
	}

	// FIXME (etd) - Commenting out the below -- we can't really access devices via index anymore,
	// at least not with how the tests are current set up.

	manufacturerDevice, manufacturerDeviceKind, manufacturerInstance := FindDeviceInstanceByInfo(devices, "upsIdentManufacturer")
	fmt.Printf("manufacturerDevice: %+v\n", manufacturerDevice)
	fmt.Printf("manufacturerDeviceKind: %+v\n", manufacturerDeviceKind)
	fmt.Printf("manufacturerInstance: %+v\n", manufacturerInstance)

	if manufacturerInstance.Data["table_name"] != "UPS-MIB-UPS-Identity-Table" {
		t.Fatalf("Expected TableName == [UPS-MIB-UPS-Identity-Table], got [%v]", manufacturerInstance.Data["table_name"])
	}
	if manufacturerDeviceKind.Name != "identity" {
		t.Fatalf("Expected Name == [identity], got [%v]", manufacturerDeviceKind.Name)
	}
	if manufacturerInstance.Data["oid"] != ".1.3.6.1.2.1.33.1.1.1.0" {
		t.Fatalf("Expected oid == [.1.3.6.1.2.1.33.1.1.1.0], got [%v]", manufacturerInstance.Data["oid"])
	}

	powerDevice, powerDeviceKind, powerInstance := FindDeviceInstanceByInfo(devices, "upsInputTruePower1")
	fmt.Printf("powerDevice: %+v\n", powerDevice)
	fmt.Printf("powerDeviceKind: %+v\n", powerDeviceKind)
	fmt.Printf("powerInstance: %+v\n", powerInstance)

	if powerInstance.Data["table_name"] != "UPS-MIB-UPS-Input-Table" {
		t.Fatalf("Expected TableName == [UPS-MIB-UPS-Input-Table], got [%v]", powerInstance.Data["table_name"])
	}
	if powerDeviceKind.Name != "power" {
		t.Fatalf("Expected Name == [power], got [%v]", powerDeviceKind.Name)
	}
	if powerInstance.Data["oid"] != ".1.3.6.1.2.1.33.1.3.3.1.5.2" {
		t.Fatalf("Expected oid == [.1.3.6.1.2.1.33.1.3.3.1.5.2], got [%v]", powerInstance.Data["oid"])
	}
	t.Logf("TestUpsMib end")
}
