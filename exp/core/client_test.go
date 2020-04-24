package core

import (
	"testing"
	"time"

	"github.com/soniah/gosnmp"
	"github.com/stretchr/testify/assert"
)

func TestClient_Close(t *testing.T) {
	c := Client{
		GoSNMP: &gosnmp.GoSNMP{},
	}
	c.Close()
}

func TestNewClient(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "AuthPriv",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, "localhost", client.Target)
	assert.Equal(t, uint16(1024), client.Port)
	assert.Equal(t, "udp", client.Transport)
	assert.Equal(t, "", client.Community)
	assert.Equal(t, gosnmp.Version3, client.Version)
	assert.Equal(t, 1*time.Second, client.Timeout)
	assert.Equal(t, 1, client.Retries)
}

func TestNewClient2(t *testing.T) {
	// Same test, different configuration
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v2",
		Agent:   "tcp://localhost",
		Timeout: 1 * time.Second,
		Retries: 1,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, "localhost", client.Target)
	assert.Equal(t, uint16(161), client.Port)
	assert.Equal(t, "tcp", client.Transport)
	assert.Equal(t, "", client.Community)
	assert.Equal(t, gosnmp.Version2c, client.Version)
	assert.Equal(t, 1*time.Second, client.Timeout)
	assert.Equal(t, 1, client.Retries)
}

func TestNewClient3(t *testing.T) {
	// Same test, different configuration
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v1",
		Agent:   "localhost:4321",
		Timeout: 1 * time.Second,
		Retries: 1,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, "localhost", client.Target)
	assert.Equal(t, uint16(4321), client.Port)
	assert.Equal(t, "udp", client.Transport)
	assert.Equal(t, "", client.Community)
	assert.Equal(t, gosnmp.Version1, client.Version)
	assert.Equal(t, 1*time.Second, client.Timeout)
	assert.Equal(t, 1, client.Retries)
}

func TestNewClient_NoSecurity(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, "localhost", client.Target)
	assert.Equal(t, uint16(1024), client.Port)
	assert.Equal(t, "udp", client.Transport)
	assert.Equal(t, "", client.Community)
	assert.Equal(t, gosnmp.Version3, client.Version)
	assert.Equal(t, 1*time.Second, client.Timeout)
	assert.Equal(t, 1, client.Retries)
}

func TestNewClient_BadVersion(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "invalid-version",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_BadAgent(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "this is not a valid agent string!",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_BadTransport(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "something://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_BadPort(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:not-a-port",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_BadPort2(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:999999999999",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_v2WithSecurity(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v2",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "AuthPriv",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrNonV3SecurityParams, err)
}

func TestNewClient_BadSecLevel(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "invalid-level",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrInvalidMessageFlag, err)
}

func TestNewClient_BadAuthProtocol(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "AuthPriv",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "invalid-protocol",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrInvalidAuthProtocol, err)
}

func TestNewClient_BadPrivProtocol(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "AuthPriv",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "invalid-protocol",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrInvalidPrivProtocol, err)
}
