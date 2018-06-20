package devices

import (
	"fmt"
	"testing"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/mibs/ups_mib"
)

// Create Device creates the Device structure in test land for now.
func CreateDevices(config *sdk.DeviceConfig, handler *sdk.DeviceHandler) []*sdk.Device {
	var devices []*sdk.Device

	for _, device := range config.Devices {
		for _, instance := range device.Instances {
			device := &sdk.Device{
				Info:     instance.Info,
				Data:     instance.Data,
				Kind:     device.Name,
				Location: &sdk.Location{Rack: "rack", Board: "board"},
				Handler:  handler,
			}
			devices = append(devices, device)
		}
	}
	return devices
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
		if "data is nil" != err.Error() {
			t.Fatalf("Expected err: [data is nil], got [%v]", err.Error())
		}
	}

	// No rack.
	_, err = testUpsMib.EnumerateDevices(map[string]interface{}{})
	if err == nil {
		if "rack is not in data" != err.Error() {
			t.Fatalf("Expected err: [rack is not in data], got [%v]", err.Error())
		}
	}

	// Rack is not a string.
	_, err = testUpsMib.EnumerateDevices(map[string]interface{}{"rack": 3})
	if err == nil {
		if "rack is not a string, int" != err.Error() {
			t.Fatalf(
				"Expected err: [rack is not a string, int], got [%v]", err.Error())
		}
	}

	// No board.
	_, err = testUpsMib.EnumerateDevices(
		map[string]interface{}{"rack": "test_rack"})
	if err == nil {
		if "board is not in data" != err.Error() {
			t.Fatalf("Expected err: [board is not in data], got [%v]", err.Error())
		}
	}

	// Board is not a string.
	_, err = testUpsMib.EnumerateDevices(
		map[string]interface{}{"rack": "test_rack", "board": -1})
	if err == nil {
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

	//func DumpDeviceConfigs(devices []*config.DeviceConfig, header string) {
	DumpDeviceConfigs(snmpDevices, "Devices from UPS-MIB")
	if len(snmpDevices) != 40 {
		t.Fatalf("Expected 40 snmp devices, got %d.", len(snmpDevices))
	}

	// Find power device configs in the UPS MIB. There should be six.
	powerDeviceConfigs, err := FindDeviceConfigsByType(snmpDevices, "power")
	if err != nil {
		t.Fatal(err)
	}
	DumpDeviceConfigs(powerDeviceConfigs, "Power device configs")
	if len(powerDeviceConfigs) != 6 {
		t.Fatalf("Expected 6 power device configs, got %d", len(powerDeviceConfigs))
	}

	// At long last we should be able to create the Device structure.
	powerDevices := CreateDevices(powerDeviceConfigs[0], &SnmpPower)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("powerDevice: %+v\n", powerDevices)

	if len(powerDevices) != 1 {
		t.Fatalf("Expected 1 power device, got %d.", len(powerDevices))
	}

	powerDevice := powerDevices[0]
	// Get the first reading.
	context, err := powerDevice.Read() // Call Read through the device's function pointer.
	if err != nil {
		t.Fatal(err)
	}
	readings := context.Reading
	for i := 0; i < len(readings); i++ {
		fmt.Printf("Reading[%d]: %T, %+v\n", i, readings[i], readings[i])
	}

	// For each device config, create a device and perform a reading.
	var devices []*sdk.Device
	DumpDeviceConfigs(snmpDevices, "Second device dump:")

	for _, snmpDevice := range snmpDevices {
		for _, kind := range snmpDevice.Devices {
			var deviceHandler *sdk.DeviceHandler

			switch typ := kind.Name; typ {
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

			// This is not a good way to do this. This is just done here to
			// get the tests building again. Tests will need to be refactored.
			tmpConfig := &sdk.DeviceConfig{
				SchemeVersion: sdk.SchemeVersion{Version: "1.0"},
				Locations:     []*sdk.LocationConfig{},
				Devices:       []*sdk.DeviceKind{kind},
			}

			devs := CreateDevices(tmpConfig, deviceHandler)
			devices = append(devices, devs...)
		}
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
}
