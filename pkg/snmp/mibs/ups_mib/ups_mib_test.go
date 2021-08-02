package mibs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

func findDeviceInstanceByInfo(devices []*config.DeviceProto, info string) (
	*config.DeviceProto, *config.DeviceInstance) {

	for _, proto := range devices {
		for _, instance := range proto.Instances {
			// This only works since info is coded to be unique.
			if instance.Info == info {
				return proto, instance
			}
		}
	}
	return nil, nil
}

// TestUpsMib
// Initial test creates all tables based on the UPS-MIB.
func TestUpsMib(t *testing.T) { // nolint: gocyclo
	// In order to create the table, we need to create an SNMP Server.
	// In order to create the SNMP server, we need to have an SnmpClient.

	// Create SecurityParameters for the config that should connect to the emulator.
	// TODO (etd): This setup bit is used in a few places, could create a test util for it.
	securityParameters, err := core.NewSecurityParameters(
		"simulator",  // User Name
		core.SHA,     // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		core.AES,     // Privacy Protocol
		"privatus",   // Privacy Passphrase
	)
	assert.NoError(t, err)

	// Create a config.
	cfg, err := core.NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public", //  Context name
	)
	assert.NoError(t, err)

	// Create a client.
	client, err := core.NewSnmpClient(cfg)
	assert.NoError(t, err)

	// Create SnmpServerBase
	snmpServer, err := core.NewSnmpServerBase(client, cfg)
	assert.NoError(t, err)

	// Create the UpsMib and dump it.
	testUpsMib, err := NewUpsMib(snmpServer)
	assert.NoError(t, err)

	testUpsMib.Dump()

	// We should have 19 tables.
	assert.Len(t, testUpsMib.Tables, 19)

	// Get the ups identity data from the test MIB.
	upsIdentity := testUpsMib.UpsIdentityTable.UpsIdentity

	// Verify expected ups identity data from the test MIB.
	assert.NotNil(t, upsIdentity)
	assert.Equal(t, "Eaton Corporation", upsIdentity.Manufacturer)
	assert.Equal(t, "PXGMS UPS + EATON 93PM", upsIdentity.Model)
	assert.Equal(t, "INV: 1.44.0000", upsIdentity.UpsSoftwareVersion)
	assert.Equal(t, "2.3.7", upsIdentity.AgentSoftwareVersion)
	assert.Equal(t, "ID: EM111UXX06, Msg: 9PL15N0000E40R2", upsIdentity.Name)
	assert.Equal(t, "Attached Devices not set", upsIdentity.AttachedDevices)

	// Call the ups battery table device enumerator.
	upsBatteryTable := testUpsMib.UpsBatteryTable
	devices, err := upsBatteryTable.SnmpTable.DevEnumerator.DeviceEnumerator(
		map[string]interface{}{"rack": "my_pet_rack", "board": "my_pet_board"},
	)
	assert.NoError(t, err)
	assert.Greater(t, len(devices), 0)

	// Check the number of device instances that were created
	instanceCount := 0
	for _, proto := range devices {
		instanceCount += len(proto.Instances)
	}
	assert.Equal(t, 7, instanceCount, "upsBatteryTable")

	// Enumerate UpsInputTable devices.
	upsInputTable := testUpsMib.UpsInputTable
	devices, err = upsInputTable.SnmpTable.DevEnumerator.DeviceEnumerator(
		map[string]interface{}{"rack": "my_pet_rack", "board": "my_pet_board"},
	)
	assert.NoError(t, err)
	assert.Greater(t, len(devices), 0)

	// Check the number of device instances that were created
	instanceCount = 0
	for _, proto := range devices {
		instanceCount += len(proto.Instances)
	}
	assert.Equal(t, 12, instanceCount, "upsInputTable")

	// Enumerate the UpsAlarmsHeadersTable devices.
	upsAlarmsHeadersTable := testUpsMib.UpsAlarmsHeadersTable
	devices, err = upsAlarmsHeadersTable.SnmpTable.DevEnumerator.DeviceEnumerator(
		map[string]interface{}{"rack": "my_pet_rack", "board": "my_pet_board"},
	)
	assert.NoError(t, err)
	assert.Greater(t, len(devices), 0)

	// Check the number of device instances that were created
	instanceCount = 0
	for _, cfg := range devices {
		instanceCount += len(cfg.Instances)
	}
	assert.Equal(t, 1, instanceCount, "upsAlarmsHeaderTable")

	// Enumerate the mib.
	// Testing for bad parameters is in TestDevices.
	devices, err = testUpsMib.EnumerateDevices(
		map[string]interface{}{"rack": "test_rack", "board": "test_board"},
	)
	assert.NoError(t, err)
	assert.Greater(t, len(devices), 0)

	// Check the number of device instances that were created
	instanceCount = 0
	for _, proto := range devices {
		instanceCount += len(proto.Instances)
	}
	assert.Equal(t, 54, instanceCount, "devices")

	t.Log("Dumping devices enumerated from UPS-MIB")
	for _, proto := range devices {
		for _, instance := range proto.Instances {
			t.Logf("UPS-MIB device: %v %v %v %v %v row:%v column:%v",
				instance.Data["table_name"],
				proto.Type,
				instance.Info,
				instance.Data["oid"],
				instance.Data["base_oid"],
				instance.Data["row"],
				instance.Data["column"],
			)
		}
	}

	manufacturerProto, manufacturerInstance := findDeviceInstanceByInfo(devices, "upsIdentManufacturer")
	t.Logf("manufacturerDeviceProto: %+v", manufacturerProto)
	t.Logf("manufacturerInstance: %+v", manufacturerInstance)

	assert.Equal(t, "UPS-MIB-UPS-Identity-Table", manufacturerInstance.Data["table_name"])
	assert.Equal(t, "identity", manufacturerProto.Type)
	assert.Equal(t, ".1.3.6.1.2.1.33.1.1.1.0", manufacturerInstance.Data["oid"])

	powerProto, powerInstance := findDeviceInstanceByInfo(devices, "upsInputTruePower1")
	t.Logf("powerDeviceProto: %+v", powerProto)
	t.Logf("powerInstance: %+v", powerInstance)

	assert.Equal(t, "UPS-MIB-UPS-Input-Table", powerInstance.Data["table_name"])
	assert.Equal(t, "power", powerProto.Type)
	assert.Equal(t, ".1.3.6.1.2.1.33.1.3.3.1.5.2", powerInstance.Data["oid"])
}
