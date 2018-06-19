package devices

import (
	"fmt"
	"os"
	"testing"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/mibs/ups_mib"
)

// ParseProtoypeConfigs is a wrapper around config.ParsePrototyeConfig() that
// takes a directory parameter for sanity.
func ParsePrototypeConfigs(prototypeDirectory string) (
	prototypeConfigs []*sdk.PrototypeConfig, err error) {

	// Set EnvProtoPath.
	err = os.Setenv(sdk.EnvProtoPath, prototypeDirectory)
	if err != nil {
		return nil, err
	}
	// Unset env on exit.
	defer func() {
		_ = os.Unsetenv(sdk.EnvProtoPath)
	}()

	// Parse the Protoype configuration.
	prototypeConfigs, err = sdk.ParsePrototypeConfig()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(prototypeConfigs); i++ {
		fmt.Printf("prototypeConfigs[%d]: %+v\n", i, prototypeConfigs[i])
	}
	return prototypeConfigs, nil
}

// FindPrototypeConfigByType finds a prototype config in the given set where
// Type matches t or nil if not found.
func FindPrototypeConfigByType(
	prototypeConfigs []*sdk.PrototypeConfig, t string) (
	prototypeConfig *sdk.PrototypeConfig) {
	if prototypeConfigs == nil {
		return nil
	}
	for i := 0; i < len(prototypeConfigs); i++ {
		if prototypeConfigs[i].Type == t {
			return prototypeConfigs[i]
		}
	}
	return nil
}

// Create Device creates the Device structure in test land for now.
func CreateDevice(
	deviceConfig *sdk.DeviceConfig,
	prototypeConfig *sdk.PrototypeConfig,
	deviceHandler *sdk.DeviceHandler,
	plugin *sdk.Plugin) (device *sdk.Device, err error) {

	return sdk.NewDevice(
		prototypeConfig,
		deviceConfig,
		deviceHandler,
		plugin)
}

// testDeviceIdentifier is here so that we can create a plugin.
func testDeviceIdentifier(x map[string]string) string {
	return ""
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

	t.Logf("client: %+v", client)

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

	// Parse out all of the prototype configs. Expect 7.
	prototypeConfigs, err := ParsePrototypeConfigs("../config/proto")
	if err != nil {
		t.Fatal(err)
	}
	if len(prototypeConfigs) != 7 {
		t.Fatalf("Expected 7 prototype configs, got %d", len(prototypeConfigs))
	}

	// Find the power prototype config.
	powerPrototypeConfig := FindPrototypeConfigByType(prototypeConfigs, "power")
	fmt.Printf("powerPrototypeConfig: %+v\n", powerPrototypeConfig)

	// Setup the power device handler.
	powerDeviceHandler := &SnmpPower
	fmt.Printf("powerDeviceHandler: %+v\n", powerDeviceHandler)

	// Need handlers to create a plugin.
	handlers, err := sdk.NewHandlers(
		testDeviceIdentifier, testUpsMib.EnumerateDevices)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("handlers: %+v\n", handlers)

	// Need a plugin config to create a plugin.
	pluginConfig := sdk.PluginConfig{
		Name:    "test config",
		Version: "test config v1",
		Network: sdk.NetworkSettings{
			Type:    "tcp",
			Address: "test_config",
		},
		Settings: sdk.Settings{
			Read:        sdk.ReadSettings{Buffer: 1024},
			Write:       sdk.WriteSettings{Buffer: 1024},
			Transaction: sdk.TransactionSettings{TTL: "2s"},
		},
	}
	fmt.Printf("pluginConfig: %+v\n", pluginConfig)

	// Create a plugin.
	plugin, err := sdk.NewPlugin(handlers, &pluginConfig)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("plugin: %+v\n", plugin)

	// At long last we should be able to create the Device structure.
	powerDevice, err := CreateDevice(
		powerDeviceConfigs[0], powerPrototypeConfig, powerDeviceHandler, plugin)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("powerDevice: %+v\n", powerDevice)

	// Get the first reading.
	context, err := powerDevice.Read() // Call Read through the device's function pointer.
	if err != nil {
		t.Fatal(err)
	}
	readings := context.Reading
	for i := 0; i < len(readings); i++ {
		fmt.Printf("Reading[%d]: %T, %+v\n", i, readings[i], readings[i])
	}

	// Get the rest of the prototype configs and DeviceHandlers.
	currentPrototypeConfig := FindPrototypeConfigByType(
		prototypeConfigs, "current")
	fmt.Printf("currentPrototypeConfig: %+v\n", currentPrototypeConfig)

	currentDeviceHandler := &SnmpCurrent
	fmt.Printf("currentDeviceHandler: %+v\n", currentDeviceHandler)

	frequencyPrototypeConfig := FindPrototypeConfigByType(
		prototypeConfigs, "frequency")
	fmt.Printf("frequencyPrototypeConfig: %+v\n", frequencyPrototypeConfig)

	frequencyDeviceHandler := &SnmpFrequency
	fmt.Printf("frequencyDeviceHandler: %+v\n", frequencyDeviceHandler)

	identityPrototypeConfig := FindPrototypeConfigByType(
		prototypeConfigs, "identity")
	fmt.Printf("identityPrototypeConfig: %+v\n", identityPrototypeConfig)

	identityDeviceHandler := &SnmpIdentity
	fmt.Printf("identityDeviceHandler: %+v\n", identityDeviceHandler)

	statusPrototypeConfig := FindPrototypeConfigByType(prototypeConfigs, "status")
	fmt.Printf("statusPrototypeConfig: %+v\n", statusPrototypeConfig)

	statusDeviceHandler := &SnmpStatus
	fmt.Printf("statusDeviceHandler: %+v\n", statusDeviceHandler)

	temperaturePrototypeConfig := FindPrototypeConfigByType(
		prototypeConfigs, "temperature")
	fmt.Printf("temperaturePrototypeConfig: %+v\n", temperaturePrototypeConfig)

	temperatureDeviceHandler := &SnmpTemperature
	fmt.Printf("temperatureDeviceHandler: %+v\n", temperatureDeviceHandler)

	voltagePrototypeConfig := FindPrototypeConfigByType(
		prototypeConfigs, "voltage")
	fmt.Printf("voltagePrototypeConfig: %+v\n", voltagePrototypeConfig)

	voltageDeviceHandler := &SnmpVoltage
	fmt.Printf("voltageDeviceHandler: %+v\n", voltageDeviceHandler)

	// For each device config, create a device and perform a reading.
	var devices []*sdk.Device
	DumpDeviceConfigs(snmpDevices, "Second device dump:")

	for i := 0; i < len(snmpDevices); i++ {
		fmt.Printf("snmpDevice[%d]: %+v\n", i, snmpDevices[i])

		var protoConfig *sdk.PrototypeConfig
		var deviceHandler *sdk.DeviceHandler

		switch typ := snmpDevices[i].Type; typ {
		case "current":
			protoConfig = currentPrototypeConfig
			deviceHandler = currentDeviceHandler
		case "frequency":
			protoConfig = frequencyPrototypeConfig
			deviceHandler = frequencyDeviceHandler
		case "identity":
			protoConfig = identityPrototypeConfig
			deviceHandler = identityDeviceHandler
		case "power":
			protoConfig = powerPrototypeConfig
			deviceHandler = powerDeviceHandler
		case "status":
			protoConfig = statusPrototypeConfig
			deviceHandler = statusDeviceHandler
		case "temperature":
			protoConfig = temperaturePrototypeConfig
			deviceHandler = temperatureDeviceHandler
		case "voltage":
			protoConfig = voltagePrototypeConfig
			deviceHandler = voltageDeviceHandler
		default:
			t.Fatalf("Unknown type: %v", typ)
		}

		device, err := CreateDevice(
			snmpDevices[i], protoConfig, deviceHandler, plugin)
		if err != nil {
			t.Fatal(err)
		}
		devices = append(devices, device)
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
