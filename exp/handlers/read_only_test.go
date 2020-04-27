package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
)

func TestReadOnlyReadHandler_NilDevice(t *testing.T) {
	readings, err := readOnlyReadHandler(nil)

	assert.Error(t, err)
	assert.Nil(t, readings)
}

func TestReadOnlyReadHandler_NoAgent(t *testing.T) {
	readings, err := readOnlyReadHandler(&sdk.Device{
		Data: map[string]interface{}{},
	})

	assert.Error(t, err)
	assert.Nil(t, readings)
}

func TestReadOnlyReadHandler_NoOid(t *testing.T) {
	readings, err := readOnlyReadHandler(&sdk.Device{
		Data: map[string]interface{}{
			"agent": "udp://localhost:1024",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, readings)
}

func TestReadOnlyReadHandler_NoTargetCfg(t *testing.T) {
	readings, err := readOnlyReadHandler(&sdk.Device{
		Data: map[string]interface{}{
			"agent": "udp://localhost:1024",
			"oid":   "1.2.3.4",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, readings)
}

func TestReadOnlyReadHandler_FailedClientInit(t *testing.T) {
	cfg := &core.SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v2",
		Agent:   "this:is-not@a:valid_agent+string",
	}

	readings, err := readOnlyReadHandler(&sdk.Device{
		Data: map[string]interface{}{
			"agent":      "udp://localhost:1024",
			"oid":        "1.2.3.4",
			"target_cfg": cfg,
		},
	})

	assert.Error(t, err)
	assert.Nil(t, readings)
}
