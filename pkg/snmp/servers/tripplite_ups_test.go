package servers

/*
import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// TestTrippliteUps is the first TrippliteUps test.
func TestTrippliteUps(t *testing.T) {
	t.Log("TestTrippliteUps start")
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
	data["model"] = "SU10000RT3UPM"
	data["version"] = "v3"

	trippliteUps, err := NewTrippliteUps(data)
	assert.NoError(t, err)
	assert.NotNil(t, trippliteUps)
	assert.NotNil(t, trippliteUps.SnmpServer)
	assert.NotNil(t, trippliteUps.SnmpServer.SnmpServerBase)
	assert.NotNil(t, trippliteUps.SnmpServer.SnmpServerBase.SnmpClient)
	assert.NotNil(t, trippliteUps.SnmpServer.SnmpServerBase.SnmpClient.DeviceConfig)

	clientDeviceConfig := trippliteUps.SnmpServer.SnmpServerBase.SnmpClient.DeviceConfig
	assert.Equal(t, clientDeviceConfig.Version, "V3")
	assert.Equal(t, clientDeviceConfig.Endpoint, "127.0.0.1")
	assert.Equal(t, clientDeviceConfig.ContextName, "public")
	thirtySeconds, _ := time.ParseDuration("30s")
	assert.Equal(t, clientDeviceConfig.Timeout, thirtySeconds)
	assert.NotNil(t, clientDeviceConfig.SecurityParameters)
	assert.Equal(t, clientDeviceConfig.SecurityParameters.AuthenticationProtocol, core.AuthenticationProtocol(3))
	assert.Equal(t, clientDeviceConfig.SecurityParameters.PrivacyProtocol, core.PrivacyProtocol(3))
	assert.Equal(t, clientDeviceConfig.SecurityParameters.UserName, "simulator")
	assert.Equal(t, clientDeviceConfig.SecurityParameters.AuthenticationPassphrase, "auctoritas")
	assert.Equal(t, clientDeviceConfig.SecurityParameters.PrivacyPassphrase, "privatus")
	assert.Equal(t, clientDeviceConfig.Port, uint16(1024))

	assert.NotNil(t, trippliteUps.SnmpServer.SnmpServerBase.DeviceConfig)
	serverDeviceConfig := trippliteUps.SnmpServer.SnmpServerBase.SnmpClient.DeviceConfig
	assert.Equal(t, serverDeviceConfig.Version, "V3")
	assert.Equal(t, serverDeviceConfig.Endpoint, "127.0.0.1")
	assert.Equal(t, serverDeviceConfig.ContextName, "public")
	assert.Equal(t, serverDeviceConfig.Timeout, thirtySeconds)
	assert.NotNil(t, serverDeviceConfig.SecurityParameters)
	assert.Equal(t, serverDeviceConfig.SecurityParameters.AuthenticationProtocol, core.AuthenticationProtocol(3))
	assert.Equal(t, serverDeviceConfig.SecurityParameters.PrivacyProtocol, core.PrivacyProtocol(3))
	assert.Equal(t, serverDeviceConfig.SecurityParameters.UserName, "simulator")
	assert.Equal(t, serverDeviceConfig.SecurityParameters.AuthenticationPassphrase, "auctoritas")
	assert.Equal(t, serverDeviceConfig.SecurityParameters.PrivacyPassphrase, "privatus")
	assert.Equal(t, serverDeviceConfig.Port, uint16(1024))

	// The verification here is done with emulator data from a different type of UPS.
	// In the future we can use data from other UPSes and get different results,
	// that just hasn't happened yet.

	// Verify device handlers by type.
	deviceHandlersByType := map[string]int{}
	for i := 0; i < len(trippliteUps.SnmpServer.DeviceConfigs); i++ {
		dhType := trippliteUps.SnmpServer.DeviceConfigs[i].Type
		count, ok := deviceHandlersByType[dhType]
		if ok {
			deviceHandlersByType[dhType] = count + 1
		} else {
			deviceHandlersByType[dhType] = 1
		}
	}

	assert.Equal(t, 10, len(deviceHandlersByType))
	assert.Equal(t, 4, deviceHandlersByType["current"])
	assert.Equal(t, 2, deviceHandlersByType["frequency"])
	assert.Equal(t, 1, deviceHandlersByType["identity"])
	assert.Equal(t, 1, deviceHandlersByType["minutes"])
	assert.Equal(t, 2, deviceHandlersByType["percentage"])
	assert.Equal(t, 3, deviceHandlersByType["power"])
	assert.Equal(t, 1, deviceHandlersByType["seconds"])
	assert.Equal(t, 5, deviceHandlersByType["status"])
	assert.Equal(t, 1, deviceHandlersByType["temperature"])
	assert.Equal(t, 4, deviceHandlersByType["voltage"])

}

// TestTrippliteUpsInitializationFailure tests that we get an error when we cannot initialize the UPS MIB.
// This test uses a privacy failure to fail initialization.
func TestTrippliteUpsInitializationFailure(t *testing.T) {
	t.Log("TestTrippliteUpsInitializationFailure start")
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
	data["model"] = "SU10000RT3UPM"
	data["version"] = "v3"

	_, err := NewTrippliteUps(data)
	assert.Error(t, err)
	assert.Equal(t, "incoming packet is not authentic, discarding", err.Error())
}
*/
