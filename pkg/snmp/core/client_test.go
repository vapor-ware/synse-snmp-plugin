package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestClientPxgmsUps is the initial positive test against the emulator.
// Uses SNMPv3 with SHA/AES.
func TestClientPxgmsUps(t *testing.T) {
	// Create SecurityParameters for the config that should connect to the emulator.
	securityParameters, err := NewSecurityParameters(
		"simulator",  // User Name
		SHA,          // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		AES,          // Privacy Protocol
		"privatus",   // Privacy Passphrase
	)
	assert.NoError(t, err)

	// Create a config.
	config, err := NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public",   //  Context name
		[]string{}, // tags (none)
	)
	assert.NoError(t, err)

	// Create a client.
	client, err := NewSnmpClient(config)
	assert.NoError(t, err)

	// Walk SNMP OID "1.3.6.1" and print results.
	results, err := client.Walk("1.3.6.1")
	assert.NoError(t, err)

	// Log output.
	for i, result := range results {
		t.Logf("%d: OID: %v, Data: %v", i, result.Oid, result.Data)
	}
	// Assert result set size.
	assert.Equal(t, 709, len(results))
}

// getExpectedConfigShaAes gets an expected valid device configuration.
func getExpectedConfigShaAes() *DeviceConfig {
	return &DeviceConfig{
		Version:     "V3",
		Endpoint:    "127.0.0.1",
		Port:        1024,
		ContextName: "public",
		Timeout:     time.Duration(30) * time.Second,
		SecurityParameters: &SecurityParameters{
			UserName:                 "simulator",
			AuthenticationProtocol:   SHA,
			AuthenticationPassphrase: "auctorias",
			PrivacyProtocol:          AES,
			PrivacyPassphrase:        "privatus",
		},
	}
}

// getExpectedConfigShaAes gets an expected valid device configuration.
func getExpectedConfigMd5Des() *DeviceConfig {
	return &DeviceConfig{
		Version:     "V3",
		Endpoint:    "127.0.0.1",
		Port:        1024,
		ContextName: "public",
		Timeout:     time.Duration(30) * time.Second,
		SecurityParameters: &SecurityParameters{
			UserName:                 "simulator",
			AuthenticationProtocol:   MD5,
			AuthenticationPassphrase: "auctorias",
			PrivacyProtocol:          DES,
			PrivacyPassphrase:        "privatus",
		},
	}
}

// verifyConfig checks expected DeviceConfig fields versus the actual.
func verifyConfig(t *testing.T, expected *DeviceConfig, actual *DeviceConfig) {
	assert.Equal(t, expected.Version, actual.Version)
	assert.Equal(t, expected.Endpoint, actual.Endpoint)
	assert.Equal(t, expected.Port, actual.Port)
	assert.Equal(t, expected.ContextName, actual.ContextName)
	assert.Equal(t, expected.SecurityParameters.UserName, actual.SecurityParameters.UserName)
	assert.Equal(t, expected.SecurityParameters.AuthenticationProtocol, actual.SecurityParameters.AuthenticationProtocol)
	assert.Equal(t, expected.SecurityParameters.AuthenticationPassphrase, actual.SecurityParameters.AuthenticationPassphrase)
	assert.Equal(t, expected.SecurityParameters.PrivacyProtocol, actual.SecurityParameters.PrivacyProtocol)
	assert.Equal(t, expected.SecurityParameters.PrivacyPassphrase, actual.SecurityParameters.PrivacyPassphrase)
}

// Configuration tests.

// Test a valid configuration.
func TestValidConfigMapShaAes(t *testing.T) {
	yamlConfig := map[string]interface{}{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     1024,
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	// actual is a SNMP DeviceConfig created by the constructor.
	actual, err := GetDeviceConfig(yamlConfig)
	assert.NoError(t, err)

	// Test each field against expected.
	expected := getExpectedConfigShaAes()
	verifyConfig(t, expected, actual)
}

func TestValidConfigMapMd5Des(t *testing.T) {
	yamlConfig := map[string]interface{}{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     1024,
		"userName":                 "simulator",
		"authenticationProtocol":   "MD5",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "DES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	// actual is a SNMP DeviceConfig created by the constructor.
	actual, err := GetDeviceConfig(yamlConfig)
	assert.NoError(t, err)

	// Test each field against expected.
	expected := getExpectedConfigMd5Des()
	verifyConfig(t, expected, actual)
}

// Test invalid configurations, one for each required field missing or invalid.

func TestConfigMapInvalidVersion(t *testing.T) {
	yamlConfig := map[string]interface{}{
		"version":                  "v2c", // v2c is not currently supported.
		"endpoint":                 "127.0.0.1",
		"port":                     1024,
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	assert.Error(t, err)
	assert.Equal(t, "version [v2c] unsupported", err.Error())
}

func TestConfigMapForgotVersion(t *testing.T) {
	yamlConfig := map[string]interface{}{
		//"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     1024,
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	assert.Error(t, err)
	assert.Equal(t, "version should be a string", err.Error())
}

func TestConfigMapForgotEndpoint(t *testing.T) {
	yamlConfig := map[string]interface{}{
		"version": "v3",
		//"endpoint":                 "127.0.0.1",
		"port":                     1024,
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	assert.Error(t, err)
	assert.Equal(t, "endpoint should be a string", err.Error())
}

func TestConfigMapNonNumericPort(t *testing.T) {
	yamlConfig := map[string]interface{}{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     "Z",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	assert.Error(t, err)
	assert.Equal(t, "port should be an int or uint16", err.Error())
}

func TestConfigMapForgotPort(t *testing.T) {
	yamlConfig := map[string]interface{}{
		"version":  "v3",
		"endpoint": "127.0.0.1",
		//"port":                     "Z",
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
	}
	_, err := GetDeviceConfig(yamlConfig)
	assert.Error(t, err)
	assert.Equal(t, "port required, but not specified", err.Error())
}

func TestConfigMapForgotContextName(t *testing.T) {
	// This one should be okay. Empty string is valid for context name and is in
	// fact the SNMP default.
	yamlConfig := map[string]interface{}{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     1024,
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
	}
	_, err := GetDeviceConfig(yamlConfig)
	assert.NoError(t, err)
}

// Test a valid configuration with an additional field.
// The additional field should be ignored.
func TestValidConfigMapShaAesExtraField(t *testing.T) {
	yamlConfig := map[string]interface{}{
		"version":                  "v3",
		"endpoint":                 "127.0.0.1",
		"port":                     1024,
		"userName":                 "simulator",
		"authenticationProtocol":   "SHA",
		"authenticationPassphrase": "auctorias",
		"privacyProtocol":          "AES",
		"privacyPassphrase":        "privatus",
		"contextName":              "public",
		"extraBaggage":             "anything",
	}
	// actual is a SNMP DeviceConfig created by the constructor.
	actual, err := GetDeviceConfig(yamlConfig)
	assert.NoError(t, err)

	// Test each field against expected.
	expected := getExpectedConfigShaAes()
	verifyConfig(t, expected, actual)
}

// TestDeviceConfigSerialization tests serialization to and from a map[string]string.
func TestDeviceConfigSerialization(t *testing.T) {
	// Create SecurityParameters for the config that should connect to the emulator.
	securityParameters, err := NewSecurityParameters(
		"simulator",  // User Name
		SHA,          // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		AES,          // Privacy Protocol
		"privatus",   // Privacy Passphrase
	)
	assert.NoError(t, err)

	// Create a config.
	config, err := NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1024,        // Port
		securityParameters,
		"public",   //  Context name
		[]string{}, // tags (none)
	)
	assert.NoError(t, err)

	// Serialize
	serialized, err := config.ToMap()
	assert.NoError(t, err)

	// Deserialize
	deserialized, err := GetDeviceConfig(serialized)
	assert.NoError(t, err)

	// Compare. config is the expected (original), deserialized is actual.
	verifyConfig(t, config, deserialized)
}

// TestClientTrippliteUps is the initial positive test against the emulator.
// Uses SNMPv3 with MD5/DES.
func TestClientTrippliteUps(t *testing.T) {
	// Create SecurityParameters for the config that should connect to the emulator.
	securityParameters, err := NewSecurityParameters(
		"simulator",  // User Name
		MD5,          // Authentication Protocol
		"auctoritas", // Authentication Passphrase
		DES,          // Privacy Protocol
		"privatus",   // Privacy Passphrase
	)
	assert.NoError(t, err)

	// Create a config.
	config, err := NewDeviceConfig(
		"v3",        // SNMP v3
		"127.0.0.1", // Endpoint
		1025,        // Port
		securityParameters,
		"public",   //  Context name
		[]string{}, // tags (none)
	)
	assert.NoError(t, err)

	// Create a client.
	client, err := NewSnmpClient(config)
	assert.NoError(t, err)

	// Walk SNMP OID "1.3.6.1" and print results.
	results, err := client.Walk("1.3.6.1")
	assert.NoError(t, err)

	// Log output.
	for i, result := range results {
		t.Logf("%d: OID: %v, Data: %v", i, result.Oid, result.Data)
	}
	// Assert result set size.
	// This will be different from the PxgmsUps because the data are different.
	assert.Equal(t, 347, len(results))
}
