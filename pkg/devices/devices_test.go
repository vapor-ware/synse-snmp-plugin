package devices

import (
	"github.com/gosnmp/gosnmp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
	mibs "github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/mibs/ups_mib"
)

// Create Device creates the Device structure in test land for now.
func CreateDevices(proto *config.DeviceProto) ([]*sdk.Device, error) {
	var devices []*sdk.Device

	deviceMap := map[string]*sdk.DeviceHandler{}
	for _, handler := range SNMPDeviceHandlers {
		deviceMap[handler.Name] = handler
	}

	for _, instance := range proto.Instances {
		device, err := sdk.NewDeviceFromConfig(proto, instance, deviceMap)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

// log information about the provided device prototypes, aiding in test debug.
func logDeviceProtos(t *testing.T, protos []*config.DeviceProto, header string) {
	t.Log("== Dumping Device Prototypes ==")
	t.Log(header)

	if protos == nil {
		t.Log(" <nil>")
		return
	}

	t.Logf(". Count device proto: %d", len(protos))
	for _, proto := range protos {
		t.Logf(".. [proto: %s] Count device instances: %d", proto.Type, len(proto.Instances))
		for _, instance := range proto.Instances {
			t.Logf("     device: %v %v %v %v %v row:%v column:%v, tags: %v\n",
				instance.Data["table_name"],
				proto.Type,
				instance.Info,
				instance.Data["oid"],
				instance.Data["base_oid"],
				instance.Data["row"],
				instance.Data["column"],
				proto.Tags,
			)
		}

	}
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
		"public",                           //  Context name
		[]string{"serverName:customerUps"}, // tags
	)
	assert.NoError(t, err)

	snmpConfig.MsgFlag = gosnmp.AuthPriv

	// Create a client.
	client, err := core.NewSnmpClient(snmpConfig)
	assert.NoError(t, err)

	// Create SnmpServerBase
	snmpServer, err := core.NewSnmpServerBase(client, snmpConfig)
	assert.NoError(t, err)

	// Create the UpsMib.
	testUpsMib, err := mibs.NewUpsMib(snmpServer)
	assert.NoError(t, err)

	// This call uses valid parameters.
	//
	// NOTE: (etd) Enumerating devices now (v3) returns a []DeviceProto, whereas in v2 it was a
	//	[]DeviceConfig, which was a higher-level aggregate of DeviceKind (the v3 DeviceProto is
	//	effectively the same thing as the v2 DeviceKind here). Since they are not yet aggregated
	// 	on return here, there will be a larger number of initial "snmpDevices" because there are
	//  device protos of the same type being generated from different tables. This initial check
	//  verifies that we get an expected number of protos from all tables. A subsequent check gets
	//  the number of unique proto types from all tables and verifies.
	//
	//	Note also that the "rack" and "board" are no longer required in v3 as there is no global
	//  "location" associated with a device configuration. Any locational information should be
	//  contained in the device Context field. Location tags may also be generated to allow filtering
	//  by location
	snmpDevices, err := testUpsMib.EnumerateDevices(map[string]interface{}{})
	assert.NoError(t, err)
	assert.Len(t, snmpDevices, 24) // all DeviceProtos from all tables.

	logDeviceProtos(t, snmpDevices, "Devices from UPS-MIB")

	// Check the number of unique DeviceProto types enumerated. Also aggregate the total
	// number of instances for all protos (e.g. the total number of devices).
	protos := map[string]int{}
	instanceCount := 0
	for _, proto := range snmpDevices {
		instanceCount += len(proto.Instances)

		// If the proto Type is already in the map, add the number of new instances,
		// otherwise initialize the entry.
		_, ok := protos[proto.Type]
		if !ok {
			protos[proto.Type] = len(proto.Instances)
		} else {
			protos[proto.Type] += len(proto.Instances)
		}
	}
	// Check the total number of unique number of device proto types
	assert.Len(t, protos, 10, protos)
	// Check the total number of device instances
	assert.Equal(t, 54, instanceCount)

	// Check the number of device instances for each device prototype.
	t.Logf("device prototype map: %#v", protos)
	assert.Equal(t, 9, protos["power"])
	assert.Equal(t, 6, protos["identity"])
	assert.Equal(t, 8, protos["status"])
	assert.Equal(t, 10, protos["voltage"])
	assert.Equal(t, 10, protos["current"])
	assert.Equal(t, 1, protos["temperature"])
	assert.Equal(t, 4, protos["frequency"])
	assert.Equal(t, 4, protos["percentage"])
	assert.Equal(t, 1, protos["minutes"])
	assert.Equal(t, 1, protos["seconds"])

	logDeviceProtos(t, snmpDevices, "Second device dump:")

	// For each device config, create a device and perform a reading.
	var devices []*sdk.Device
	for _, proto := range snmpDevices {
		devs, err := CreateDevices(proto)
		assert.NoError(t, err)

		devices = append(devices, devs...)
	}

	t.Log("Dumping all devices")
	for i := 0; i < len(devices); i++ {
		t.Logf("device[%d]: %+v", i, devices[i])
	}

	// Read each device. Check device tags.
	t.Logf("Reading each device.")
	for i := 0; i < len(devices); i++ {
		context, err := devices[i].Read() // Call Read through the device's function pointer.
		assert.NoError(t, err)

		readings := context.Reading
		// Each device currently has one reading,
		assert.Len(t, readings, 1)
		for j := 0; j < len(readings); j++ {
			t.Logf("Reading[%d][%d]: %T, %+v", i, j, readings[j], readings[j])
		}

		// Check device tags.
		assert.Len(t, devices[i].Tags, 1)
		assert.Equal(t, "default", devices[i].Tags[0].Namespace)
		assert.Equal(t, "serverName", devices[i].Tags[0].Annotation)
		assert.Equal(t, "customerUps", devices[i].Tags[0].Label)
	}
	t.Log("Finished reading each device.")
}
