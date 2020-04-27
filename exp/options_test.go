package exp

import (
	"testing"

	"github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

	"github.com/stretchr/testify/assert"
)

func TestSnmpDeviceIdentifier(t *testing.T) {
	data := map[string]interface{}{
		"oid":   "1.2.3.4.5.6",
		"mib":   "test-mib",
		"agent": "localhost:1234",
	}

	identifier := SnmpDeviceIdentifier(data)
	assert.Equal(t, "localhost:1234-test-mib:1.2.3.4.5.6", identifier)
}

func TestSnmpDeviceIdentifier_NoOid(t *testing.T) {
	data := map[string]interface{}{
		"mib":   "test-mib",
		"agent": "localhost:1234",
	}

	assert.Panics(t, func() {
		_ = SnmpDeviceIdentifier(data)
	})
}

func TestSnmpDeviceIdentifier_NoMib(t *testing.T) {
	data := map[string]interface{}{
		"oid":   "1.2.3.4.5.6",
		"agent": "localhost:1234",
	}

	assert.Panics(t, func() {
		_ = SnmpDeviceIdentifier(data)
	})
}

func TestSnmpDeviceIdentifier_NoAgent(t *testing.T) {
	data := map[string]interface{}{
		"oid": "1.2.3.4.5.6",
		"mib": "test-mib",
	}

	assert.Panics(t, func() {
		_ = SnmpDeviceIdentifier(data)
	})
}

func TestSnmpDeviceRegistrar(t *testing.T) {
	defer mibs.Clear()

	err := mibs.Register(&mibs.MIB{
		Name: "test-mib",
		Devices: []*mibs.SnmpDevice{
			{
				OID:     "1.2.3.4",
				Info:    "test device",
				Handler: "read-only",
				Type:    "temperature",
				Output:  "temperature",
			},
		},
	})
	assert.NoError(t, err)

	devices, err := SnmpDeviceRegistrar(map[string]interface{}{
		"mib":     "test-mib",
		"version": "v3",
		"agent":   "localhost",
		"timeout": "1s",
	})
	assert.NoError(t, err)
	assert.Len(t, devices, 1)
}

func TestSnmpDeviceRegistrar_FailedConfigLoad(t *testing.T) {
	devices, err := SnmpDeviceRegistrar(map[string]interface{}{
		"mib":     "test-mib",
		"version": "v3",
		"agent":   "localhost",
		"timeout": "not a valid timeout",
	})

	assert.Error(t, err)
	assert.Nil(t, devices)
}

func TestSnmpDeviceRegistrar_NoMib(t *testing.T) {
	devices, err := SnmpDeviceRegistrar(map[string]interface{}{
		// no mib defined here
		"version": "v3",
		"agent":   "localhost",
		"timeout": "1s",
	})

	assert.Error(t, err)
	assert.Nil(t, devices)
}

func TestSnmpDeviceRegistrar_FailedFindMIB(t *testing.T) {
	devices, err := SnmpDeviceRegistrar(map[string]interface{}{
		"mib":     "test-mib", // mib defined, but not registered
		"version": "v3",
		"agent":   "localhost",
		"timeout": "1s",
	})

	assert.Error(t, err)
	assert.Nil(t, devices)
}

func TestSnmpDeviceRegistrar_FailedDeviceLoad(t *testing.T) {
	defer mibs.Clear()

	err := mibs.Register(&mibs.MIB{
		Name: "test-mib",
		Devices: []*mibs.SnmpDevice{
			{
				OID:     "1.2.3.4",
				Info:    "test device",
				Handler: "read-only",
				Type:    "temperature",
				Output:  "no-such-output", // load will fail since output doesn't exist
			},
		},
	})
	assert.NoError(t, err)

	devices, err := SnmpDeviceRegistrar(map[string]interface{}{
		"mib":     "test-mib",
		"version": "v3",
		"agent":   "localhost",
		"timeout": "1s",
	})
	assert.Error(t, err)
	assert.Nil(t, devices)
}
