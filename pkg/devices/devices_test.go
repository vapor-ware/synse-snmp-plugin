package devices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
	mibs "github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/mibs/ups_mib"
)

// Create Device creates the Device structure in test land for now.
func CreateDevices(proto *config.DeviceProto, handler string) ([]*sdk.Device, error) { // nolint: gocyclo
	var devices []*sdk.Device

	for _, instance := range proto.Instances {

		device := &sdk.Device{
			Info:    instance.Info,
			Data:    instance.Data,
			Type:    proto.Type,
			Handler: handler,
		}
		devices = append(devices, device)

	}
	return devices, nil
}

// Initial device test. Ensure we can register each type the ups mib supports
// and get a reading from each.
func TestDevices(t *testing.T) { // nolint: gocyclo

	// Create SecurityParameters for the config that should connect to the emulator.
	securityParameters, err := core.NewSecurityParameters(
		"simulator",  // User Name
		core.SHA,     // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		core.AES,     // Privacy Protocol
		"privatus",   // Privacy Passphrase
	)
	assert.NoError(t, err)

	// Create a snmp config.
	snmpConfig, err := core.NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public", //  Context name
	)
	assert.NoError(t, err)

	// Create a client.
	client, err := core.NewSnmpClient(snmpConfig)
	assert.NoError(t, err)

	// Create SnmpServerBase
	snmpServer, err := core.NewSnmpServerBase(client, snmpConfig)
	assert.NoError(t, err)

	// Create the UpsMib.
	testUpsMib, err := mibs.NewUpsMib(snmpServer)
	assert.NoError(t, err)

	// Enumerate the mib. First few calls are testing bad parameters.
	_, err = testUpsMib.EnumerateDevices(nil)
	assert.Error(t, err)
	assert.Equal(t, "data is nil", err.Error())

	// No rack.
	_, err = testUpsMib.EnumerateDevices(map[string]interface{}{})
	assert.Error(t, err)
	assert.Equal(t, "rack is not in data", err.Error())

	// Rack is not a string.
	_, err = testUpsMib.EnumerateDevices(map[string]interface{}{"rack": 3})
	assert.Error(t, err)
	assert.Equal(t, "rack is not a string, int", err.Error())

	// No board.
	_, err = testUpsMib.EnumerateDevices(map[string]interface{}{"rack": "test_rack"})
	assert.Error(t, err)
	assert.Equal(t, "board is not in data", err.Error())

	// Board is not a string.
	_, err = testUpsMib.EnumerateDevices(map[string]interface{}{"rack": "test_rack", "board": -1})
	assert.Error(t, err)
	assert.Equal(t, "board is not a string, int", err.Error())

	// This call uses valid parameters.
	snmpDevices, err := testUpsMib.EnumerateDevices(map[string]interface{}{
		"rack":  "test_rack",
		"board": "test_board",
	})
	assert.NoError(t, err)
	assert.Len(t, snmpDevices, 8)

	// FIXME (etd): just log these out in the test, don't need to fmt.print..
	DumpDeviceConfigs(snmpDevices, "Devices from UPS-MIB")

	// Get the number of snmp device kinds and instances across all configs
	protos := map[string]*config.DeviceProto{}
	instanceCount := 0
	for _, proto := range snmpDevices {
		protos[proto.Type] = proto
		instanceCount += len(proto.Instances)
	}
	// Check the total number of unique number of device kinds
	assert.Len(t, protos, 9, protos)
	assert.Equal(t, 45, instanceCount)

	// Check the number of power instances
	powerInstanceCount := 0
	for _, proto := range snmpDevices {
		if proto.Type == "power" {
			powerInstanceCount += len(proto.Instances)
		}
	}
	assert.Equal(t, 6, powerInstanceCount)

	// For each device config, create a device and perform a reading.
	//var devices []*sdk.Device
	DumpDeviceConfigs(snmpDevices, "Second device dump:")

	//for _, proto := range snmpDevices {
	//	var deviceHandler *sdk.DeviceHandler
	//
	//	switch typ := proto.Type; typ {
	//	case "current":
	//		deviceHandler = &SnmpCurrent
	//	case "frequency":
	//		deviceHandler = &SnmpFrequency
	//	case "identity":
	//		deviceHandler = &SnmpIdentity
	//	case "power":
	//		deviceHandler = &SnmpPower
	//	case "status":
	//		deviceHandler = &SnmpStatus
	//	case "temperature":
	//		deviceHandler = &SnmpTemperature
	//	case "voltage":
	//		deviceHandler = &SnmpVoltage
	//	default:
	//		t.Fatalf("Unknown type: %v", typ)
	//	}
	//
	//	devs, err := CreateDevices(proto, deviceHandler)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	devices = append(devices, devs...)
	//}
	//
	//fmt.Printf("Dumping all devices\n")
	//for i := 0; i < len(devices); i++ {
	//	fmt.Printf("device[%d]: %+v\n", i, devices[i])
	//}
	//
	//// Read each device
	//fmt.Printf("Reading each device.\n")
	//for i := 0; i < len(devices); i++ {
	//	context, err := devices[i].Read() // Call Read through the device's function pointer.
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	readings := context.Reading
	//	// Each device currently has one reading,
	//	if len(readings) != 1 {
	//		t.Fatalf("Expected 1 reading for device[%d], got %d", i, len(readings))
	//	}
	//	for j := 0; j < len(readings); j++ {
	//		fmt.Printf("Reading[%d][%d]: %T, %+v\n", i, j, readings[j], readings[j])
	//	}
	//}
	//fmt.Printf("Finished reading each device.\n")
}
