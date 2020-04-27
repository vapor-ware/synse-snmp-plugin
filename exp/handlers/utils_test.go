package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
)

func TestGetAgent(t *testing.T) {
	data := map[string]interface{}{
		"agent": "udp://localhost:1024",
	}

	agent, err := getAgent(data)
	assert.NoError(t, err)
	assert.Equal(t, "udp://localhost:1024", agent)
}

func TestGetAgent_NotExist(t *testing.T) {
	data := map[string]interface{}{}

	agent, err := getAgent(data)
	assert.Error(t, err)
	assert.Equal(t, "", agent)
}

func TestGetAgent_BadType(t *testing.T) {
	data := map[string]interface{}{
		"agent": 1234,
	}

	agent, err := getAgent(data)
	assert.Error(t, err)
	assert.Equal(t, "", agent)
}

func TestGetOid(t *testing.T) {
	data := map[string]interface{}{
		"oid": "1.2.3.4",
	}

	oid, err := getOid(data)
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3.4", oid)
}

func TestGetOid_NotExist(t *testing.T) {
	data := map[string]interface{}{}

	oid, err := getOid(data)
	assert.Error(t, err)
	assert.Equal(t, "", oid)
}

func TestGetOid_BadType(t *testing.T) {
	data := map[string]interface{}{
		"oid": 1234,
	}

	oid, err := getOid(data)
	assert.Error(t, err)
	assert.Equal(t, "", oid)
}

func TestGetTargetConfig(t *testing.T) {
	cfg := core.SnmpTargetConfiguration{}
	data := map[string]interface{}{
		"target_cfg": &cfg,
	}

	targetCfg, err := getTargetConfig(data)
	assert.NoError(t, err)
	assert.Equal(t, &cfg, targetCfg)
}

func TestGetTargetConfig_NotExist(t *testing.T) {
	data := map[string]interface{}{}

	targetCfg, err := getTargetConfig(data)
	assert.Error(t, err)
	assert.Nil(t, targetCfg)
}

func TestGetTargetConfig_BadType(t *testing.T) {
	data := map[string]interface{}{
		"target_cfg": 1234,
	}

	targetCfg, err := getTargetConfig(data)
	assert.Error(t, err)
	assert.Nil(t, targetCfg)
}
