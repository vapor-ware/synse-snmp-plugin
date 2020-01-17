package devices

import (
	"testing"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
	mibs "github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/mibs/ups_mib"
)

// Create Device creates the Device structure in test land for now.
func CreateDevices(proto *config.DeviceProto, handler *sdk.DeviceHandler) ([]*sdk.Device, error) { // nolint: gocyclo
	var devices []*sdk.Device

	for _, instance := range proto.Instances {

		device := &sdk.Device{
			Info: instance.Info,
			Data: instance.Data,
			Type: proto.Type,
			//Handler:  handler,
		}
		devices = append(devices, device)

	}
	return devices, nil
}

// Initial device test. Ensure we can register each type the ups mib supports
// and get a reading from each.
func TestDevices(t *testing.T) { // nolint: gocyclo
	t.Logf("TestDevices")

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

	t.Logf("securityParameters: %+v", securityParameters)

	// Create a snmp config.
	snmpConfig, err := core.NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public") //  Context name
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("snmpConfig: %+v", snmpConfig)

	// Create a client.
	client, err := core.NewSnmpClient(snmpConfig)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	t.Logf("client: %#v", client)

	// Create SnmpServerBase
	snmpServer, err := core.NewSnmpServerBase(
		client,
		snmpConfig)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Create the UpsMib.
	testUpsMib, err := mibs.NewUpsMib(snmpServer)
	if err != nil {
		t.Fatal(err) // Fail the test.
	}

	// Enumerate the mib. First few calls are testing bad parameters.
	_, err = testUpsMib.EnumerateDevices(nil)
	if err == nil {
		t.Fatal("err is nil")
	} else {
		if "data is nil" != err.Error() {
			t.Fatalf("Expected err: [data is nil], got [%v]", err.Error())
		}
	}

	// No rack.
	_, err = testUpsMib.EnumerateDevices(map[string]interface{}{})
	if err == nil {
		t.Fatal("err is nil")
	} else {
		if "rack is not in data" != err.Error() {
			t.Fatalf("Expected err: [rack is not in data], got [%v]", err.Error())
		}
	}

	// Rack is not a string.
	_, err = testUpsMib.EnumerateDevices(map[string]interface{}{"rack": 3})
	if err == nil {
		t.Fatal("err is nil")
	} else {
		if "rack is not a string, int" != err.Error() {
			t.Fatalf(
				"Expected err: [rack is not a string, int], got [%v]", err.Error())
		}
	}

	// No board.
	_, err = testUpsMib.EnumerateDevices(
		map[string]interface{}{"rack": "test_rack"})
	if err == nil {
		t.Fatal("err is nil")
	} else {
		if "board is not in data" != err.Error() {
			t.Fatalf("Expected err: [board is not in data], got [%v]", err.Error())
		}
	}

	// Board is not a string.
	_, err = testUpsMib.EnumerateDevices(
		map[string]interface{}{"rack": "test_rack", "board": -1})
	if err == nil {
		t.Fatal("err is nil")
	} else {
		if "board is not a string, int" != err.Error() {
			t.Fatalf(
				"Expected err: [board is not a string, int], got [%v]", err.Error())
		}
	}

	// This call uses valid parameters.
	snmpDevices, err := testUpsMib.EnumerateDevices(
		map[string]interface{}{"rack": "test_rack", "board": "test_board"})
	if err != nil {
		t.Fatal(err)
	}

	DumpDeviceConfigs(snmpDevices, "Devices from UPS-MIB")
	// Check the number of snmp device configs
	if len(snmpDevices) != 6 {
		t.Fatalf("Expected 6 snmp device configs, got %d.", len(snmpDevices))
	}

	// Get the number of snmp device kinds and instances across all configs
	protos := map[string]*config.DeviceProto{}
	instanceCount := 0
	for _, proto := range snmpDevices {
		protos[proto.Type] = proto
		instanceCount += len(proto.Instances)
	}
	// Check the total number of unique number of device kinds
	if len(protos) != 7 {
		t.Logf("found kinds: %v", protos)
		t.Fatalf("Expected 7 device kinds, got %d", len(protos))
	}

	// Check the total number of device instances
	if instanceCount != 40 {
		t.Fatalf("Expected 40 instances, got %d", instanceCount)
	}

	// Check the number of power instances
	powerInstanceCount := 0
	for _, proto := range snmpDevices {
		if proto.Type == "power" {
			powerInstanceCount += len(proto.Instances)
		}
	}
	if powerInstanceCount != 6 {
		t.Fatalf("Expected 6 power device configs, got %d", powerInstanceCount)
	}

	// FIXME (etd): disabled below for now. its not currently possible to manually set the
	//   device handler for a device, so the below will not work, as we can't read without
	//   a handler. we either need some more complicated set up to build the device, or to
	//   modify the SDK to allow changing handler. not clear which is the better route at
	//   the moment.
	/*
		// For each device config, create a device and perform a reading.
		var devices []*sdk.Device
		DumpDeviceConfigs(snmpDevices, "Second device dump:")

		for _, proto := range snmpDevices {
			var deviceHandler *sdk.DeviceHandler

			switch typ := proto.Type; typ {
			case "current":
				deviceHandler = &SnmpCurrent
			case "frequency":
				deviceHandler = &SnmpFrequency
			case "identity":
				deviceHandler = &SnmpIdentity
			case "power":
				deviceHandler = &SnmpPower
			case "status":
				deviceHandler = &SnmpStatus
			case "temperature":
				deviceHandler = &SnmpTemperature
			case "voltage":
				deviceHandler = &SnmpVoltage
			default:
				t.Fatalf("Unknown type: %v", typ)
			}

			devs, err := CreateDevices(proto, deviceHandler)
			if err != nil {
				t.Fatal(err)
			}
			devices = append(devices, devs...)
		}

		fmt.Printf("Dumping all devices\n")
		for i := 0; i < len(devices); i++ {
			fmt.Printf("device[%d]: %+v\n", i, devices[i])
		}

		// Read each device
		fmt.Printf("Reading each device.\n")
		for i := 0; i < len(devices); i++ {
			context, err := devices[i].Read() // Call Read through the device's function pointer.
			if err != nil {
				t.Fatal(err)
			}
			readings := context.Reading
			// Each device currently has one reading,
			if len(readings) != 1 {
				t.Fatalf("Expected 1 reading for device[%d], got %d", i, len(readings))
			}
			for j := 0; j < len(readings); j++ {
				fmt.Printf("Reading[%d][%d]: %T, %+v\n", i, j, readings[j], readings[j])
			}
		}
		fmt.Printf("Finished reading each device.\n")
	*/
}
