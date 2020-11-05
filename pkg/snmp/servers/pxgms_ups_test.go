package servers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPxgmsUps is the first PxgmsUps test.
func TestPxgmsUps(t *testing.T) {
	t.Log("TestPxgmUps start")
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
	data["model"] = "PXGMS UPS + EATON 93PM"
	data["version"] = "v3"

	pxgmsUps, err := NewPxgmsUps(data)
	assert.NoError(t, err)

	// TODO: Need to do more with this, but at least exercises the code for now.
	t.Logf("pxgmsUps: %+v", pxgmsUps)
}

// TestPxgmsUpsInitializationFailure tests that we get an error when we cannot initialize the UPS MIB.
// This test uses an auth failure to fail initialization.
func TestPxgmsUpsInitializationFailure(t *testing.T) {
	t.Log("TestPxgmUpsInitializationFailure start")
	t.Logf("t: %+v", t)

	data := make(map[string]interface{})
	data["contextName"] = "public"
	data["endpoint"] = "127.0.0.1"
	data["userName"] = "simulator"
	data["privacyProtocol"] = "AES"
	data["privacyPassphrase"] = "incorrect_password"
	data["port"] = 1024
	data["authenticationProtocol"] = "SHA"
	data["authenticationPassphrase"] = "incorrect_authentication"
	data["model"] = "PXGMS UPS + EATON 93PM"
	data["version"] = "v3"

	_, err := NewPxgmsUps(data)
	assert.Error(t, err)
	assert.Equal(t, "Incoming packet is not authentic, discarding", err.Error())
}
