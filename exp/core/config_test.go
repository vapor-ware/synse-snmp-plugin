package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadTargetConfiguration(t *testing.T) {
	c := map[string]interface{}{
		"mib":     "test-mib",
		"version": "v3",
		"agent":   "udp://localhost:1024",
		"security": map[string]interface{}{
			"level":    "authPriv",
			"context":  "public",
			"username": "test",
			"authentication": map[string]interface{}{
				"protocol":   "SHA",
				"passphrase": "test",
			},
			"privacy": map[string]interface{}{
				"protocol":   "AES",
				"passphrase": "test",
			},
		},
	}

	cfg, err := LoadTargetConfiguration(c)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "test-mib", cfg.MIB)
	assert.Equal(t, "v3", cfg.Version)
	assert.Equal(t, "udp://localhost:1024", cfg.Agent)
	assert.Equal(t, "", cfg.Community)
	assert.Equal(t, 3*time.Second, cfg.Timeout)
	assert.Equal(t, 1, cfg.Retries)

	security := cfg.Security
	assert.NotNil(t, security)
	assert.Equal(t, "authPriv", security.Level)
	assert.Equal(t, "public", security.Context)
	assert.Equal(t, "test", security.Username)

	authentication := security.Authentication
	assert.NotNil(t, authentication)
	assert.Equal(t, "SHA", authentication.Protocol)
	assert.Equal(t, "test", authentication.Passphrase)

	privacy := security.Privacy
	assert.NotNil(t, privacy)
	assert.Equal(t, "AES", privacy.Protocol)
	assert.Equal(t, "test", privacy.Passphrase)
}

func TestLoadTargetConfiguration_BadData(t *testing.T) {
	c := map[string]interface{}{
		"version": 3,
		"agent":   []string{"localhost", "1024"},
		"security": map[string]interface{}{
			"level":    "authPriv",
			"context":  1 * time.Millisecond,
			"username": "test",
			"authentication": map[string]interface{}{
				"protocol":   "SHA",
				"passphrase": "test",
			},
			"privacy": map[string]interface{}{
				"protocol":   "AES",
				"passphrase": "test",
			},
		},
	}

	cfg, err := LoadTargetConfiguration(c)
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
