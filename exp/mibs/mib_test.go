package mibs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
)

func TestNewMIB(t *testing.T) {
	mib := NewMIB("name", &SnmpDevice{}, &SnmpDevice{}, &SnmpDevice{})

	assert.Equal(t, "name", mib.Name)
	assert.Len(t, mib.Devices, 3)
}

func TestMIB_String(t *testing.T) {
	mib := NewMIB("name", &SnmpDevice{})
	assert.Equal(t, "[MIB name]", mib.String())
}

func TestMIB_LoadDevices(t *testing.T) {
	m := MIB{
		Name: "test-mib",
		Devices: []*SnmpDevice{
			{
				OID:     "1.2.3.4",
				Info:    "test device 1",
				Type:    "temperature",
				Handler: "temperature",
				Output:  "temperature",
			},
			{
				OID:     "5.6.7.8",
				Info:    "test device 2",
				Type:    "state",
				Handler: "state",
				Output:  "state",
			},
		},
	}

	cfg := &core.SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "localhost",
	}
	devices, err := m.LoadDevices(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, devices)
	assert.Len(t, devices, 2)

	dev1 := devices[0]
	assert.Equal(t, "test device 1", dev1.Info)
	assert.Equal(t, "", dev1.Alias)
	assert.Equal(t, "temperature", dev1.Output)
	assert.Equal(t, "temperature", dev1.Handler)
	assert.Equal(t, "temperature", dev1.Type)
	assert.Equal(t, time.Duration(0), dev1.WriteTimeout)
	assert.Equal(t, map[string]string{
		"oid": "1.2.3.4",
	}, dev1.Context)
	assert.Equal(t, map[string]interface{}{
		"oid":        "1.2.3.4",
		"agent":      "localhost",
		"mib":        "test-mib",
		"target_cfg": cfg,
	}, dev1.Data)

	dev2 := devices[1]
	assert.Equal(t, "test device 2", dev2.Info)
	assert.Equal(t, "", dev2.Alias)
	assert.Equal(t, "state", dev2.Output)
	assert.Equal(t, "state", dev2.Handler)
	assert.Equal(t, "state", dev2.Type)
	assert.Equal(t, time.Duration(0), dev2.WriteTimeout)
	assert.Equal(t, map[string]string{
		"oid": "5.6.7.8",
	}, dev2.Context)
	assert.Equal(t, map[string]interface{}{
		"oid":        "5.6.7.8",
		"agent":      "localhost",
		"mib":        "test-mib",
		"target_cfg": cfg,
	}, dev2.Data)
}

func TestMIB_LoadDevices_nilConfig(t *testing.T) {
	m := MIB{
		Name: "test-mib",
		Devices: []*SnmpDevice{
			{
				OID:     "1.2.3.4",
				Info:    "test device 1",
				Type:    "temperature",
				Handler: "temperature",
				Output:  "temperature",
			},
		},
	}

	devices, err := m.LoadDevices(nil)
	assert.Error(t, err)
	assert.Nil(t, devices)
}

func TestMIB_LoadDevices_noDevices(t *testing.T) {
	m := MIB{
		Name:    "test-mib",
		Devices: []*SnmpDevice{},
	}

	devices, err := m.LoadDevices(&core.SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "localhost",
	})
	assert.NoError(t, err)
	assert.Empty(t, devices)
}

func TestMIB_LoadDevices_deviceError(t *testing.T) {
	m := MIB{
		Name: "test-mib",
		Devices: []*SnmpDevice{
			{
				OID:     "1.2.3.4",
				Info:    "test device 1",
				Type:    "temperature",
				Handler: "temperature",
				Output:  "nonexistent output type", // output does not exist
			},
		},
	}

	devices, err := m.LoadDevices(&core.SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "localhost",
	})
	assert.Error(t, err)
	assert.Nil(t, devices)
}
