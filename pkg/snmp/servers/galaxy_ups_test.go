package servers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGalaxyUps is the first GalaxyUps test.
func TestGalaxyUps(t *testing.T) {
	t.Log("TestGalaxyUps start")
	t.Logf("t: %+v", t)

	data := make(map[string]interface{})
	data["contextName"] = "public"
	data["endpoint"] = "127.0.0.1"
	data["userName"] = "simulator"
	data["privacyProtocol"] = "AES"
	data["privacyPassphrase"] = "privatus"
	data["port"] = 1024
	data["authenticationProtocol"] = "SHA"
	data["authenticationPassphrase"] = "auctoritas"
	data["model"] = "Galaxy VM 180 kVA"
	data["version"] = "v3"

	galaxyUps, err := NewGalaxyUps(data)
	assert.NoError(t, err)

	// TODO: Need to do more with this, but at least exercises the code for now.
	t.Logf("galaxyUps: %+v", galaxyUps)
}

// TestGalaxyUpsInitializationFailure tests that we get an error when we cannot initialize the UPS MIB.
// This test uses a privacy failure to fail initialization.
func TestGalaxyUpsInitializationFailure(t *testing.T) {
	t.Log("TestGalaxyUpsInitializationFailure start")
	t.Logf("t: %+v", t)

	data := make(map[string]interface{})
	data["contextName"] = "public"
	data["endpoint"] = "127.0.0.1"
	data["userName"] = "simulator"
	data["privacyProtocol"] = "AES"
	data["privacyPassphrase"] = "incorrect_password"
	data["port"] = 1024
	data["authenticationProtocol"] = "SHA"
	data["authenticationPassphrase"] = "auctorias"
	data["model"] = "Galaxy VM 180 kVA"
	data["version"] = "v3"

	_, err := NewGalaxyUps(data)
	assert.Error(t, err)
	assert.Equal(t, "Incoming packet is not authentic, discarding", err.Error())
}
