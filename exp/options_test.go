package exp

import (
	"testing"

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
