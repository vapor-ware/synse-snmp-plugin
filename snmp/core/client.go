package core

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/soniah/gosnmp"
)

// AuthenticationProtocol enumeration for authentication algorithms.
type AuthenticationProtocol uint8

const (
	// NoAuthentication for SNMP V3.
	NoAuthentication AuthenticationProtocol = 1
	// MD5 Authentication for SNMP V3.
	MD5 AuthenticationProtocol = 2
	// SHA Authentication for SNMP V3.
	SHA AuthenticationProtocol = 3
)

// PrivacyProtocol enumeration for encryption algorithms.
type PrivacyProtocol uint8

const (
	// NoPrivacy Protocol for SNMP V3.
	NoPrivacy PrivacyProtocol = 1
	// DES Privacy Protocol for SNMP V3.
	DES PrivacyProtocol = 2
	// AES Privacy Protocoli for SNMP V3.
	AES PrivacyProtocol = 3
)

// SecurityParameters is a subset of SNMP V3 USM parameters.
type SecurityParameters struct {
	AuthenticationProtocol   AuthenticationProtocol
	PrivacyProtocol          PrivacyProtocol
	UserName                 string // SNMP user name.
	AuthenticationPassphrase string
	PrivacyPassphrase        string
}

// NewSecurityParameters constructs a SecurityParameters.
func NewSecurityParameters(
	userName string,
	authenticationProtocol AuthenticationProtocol,
	authenticationPassphrase string,
	privacyProtocol PrivacyProtocol,
	privacyPassphrase string) (*SecurityParameters, error) {

	// For now, require authorization and privacy.
	// Empty user/passwords are okay.
	if !(authenticationProtocol == MD5 || authenticationProtocol == SHA) {
		return nil, fmt.Errorf("Unsupported authentication protocol [%v]",
			authenticationProtocol)
	}

	if !(privacyProtocol == DES || privacyProtocol == AES) {
		return nil, fmt.Errorf("Unsupported privacy protocol [%v]",
			privacyProtocol)
	}

	return &SecurityParameters{
		UserName:                 userName,
		AuthenticationProtocol:   authenticationProtocol,
		AuthenticationPassphrase: authenticationPassphrase,
		PrivacyProtocol:          privacyProtocol,
		PrivacyPassphrase:        privacyPassphrase,
	}, nil
}

// DeviceConfig is a thin wrapper around the configuration for gosnmp using SNMP V3.
type DeviceConfig struct {
	Version            string              // SNMP protocol version. Currently only SNMP V3 is supported.
	Endpoint           string              // Endpoint of the SNMP server to connect to.
	ContextName        string              // Context name for SNMP V3 messages.
	Timeout            time.Duration       // Timeout for the SNMP query.
	SecurityParameters *SecurityParameters // SNMP V3 security parameters.
	Port               uint16              // UDP port to connect to.
}

// checkForEmptyString checks for an empty string variable and fails with an
// attempt of a reasonable error message on failure.
func checkForEmptyString(variable string, variableName string) (err error) {
	if variable == "" {
		return fmt.Errorf("%v is an empty string, but should not be", variableName)
	}
	return nil
}

// NewDeviceConfig creates an DeviceConfig.
func NewDeviceConfig(
	version string,
	endpoint string,
	port uint16,
	securityParameters *SecurityParameters,
	contextName string) (*DeviceConfig, error) {

	// Check parameters.
	versionUpper := strings.ToUpper(version)
	if versionUpper != "V3" {
		return nil, fmt.Errorf("Version [%v] unsupported", version)
	}

	if securityParameters == nil {
		return nil, fmt.Errorf("securityParameters is nil")
	}

	// Check strings for emptyness. Version is already checked.
	err := checkForEmptyString(endpoint, "endpoint")
	if err != nil {
		return nil, err
	}

	return &DeviceConfig{
		Version:            versionUpper,
		Endpoint:           endpoint,
		Port:               port,
		SecurityParameters: securityParameters,
		ContextName:        contextName,
		Timeout:            time.Duration(30) * time.Second,
	}, nil
}

// GetDeviceConfig takes the instance configuration for an SNMP device and
// parses it into a DeviceConfig struct, filling in default values for anything
// that is missing and has a default value defined.
// This is just a deserializer which creates a DeviceConfig from
// map[string]string.
func GetDeviceConfig(instanceData map[string]interface{}) (*DeviceConfig, error) {

	// Parse out each field. The constructor call will check the parameters.
	version := fmt.Sprint(instanceData["version"])
	endpoint := fmt.Sprint(instanceData["endpoint"])
	userName := fmt.Sprint(instanceData["userName"])
	privacyPassphrase := fmt.Sprint(instanceData["privacyPassphrase"])
	authenticationPassphrase := fmt.Sprint(instanceData["authenticationPassphrase"])
	contextName := fmt.Sprint(instanceData["contextName"])

	authProtocolString := fmt.Sprint(instanceData["authenticationProtocol"])
	privProtocolString := fmt.Sprint(instanceData["privacyProtocol"])

	port, ok := instanceData["port"].(uint16)
	if !ok {
		return nil, fmt.Errorf("port should be an int")
	}

	// Only MD5 and SHA are currently supported.
	var authenticationProtocol AuthenticationProtocol
	switch strings.ToUpper(authProtocolString) {
	case "MD5":
		authenticationProtocol = MD5
	case "SHA":
		authenticationProtocol = SHA
	default:
		return nil, fmt.Errorf("Unsupported authentication protocol [%v]", authProtocolString)
	}

	// Only DES and AES are currently supported.
	var privacyProtocol PrivacyProtocol
	switch strings.ToUpper(privProtocolString) {
	case "DES":
		privacyProtocol = DES
	case "AES":
		privacyProtocol = AES
	default:
		return nil, fmt.Errorf("Unsupported privacy protocol [%v]", privProtocolString)
	}

	// Create security parameters
	securityParameters, err := NewSecurityParameters(
		userName,
		authenticationProtocol,
		authenticationPassphrase,
		privacyProtocol,
		privacyPassphrase,
	)
	if err != nil {
		return nil, err
	}

	// Create the config.
	return NewDeviceConfig(
		version,
		endpoint,
		port,
		securityParameters,
		contextName)
}

// ToMap serializes DeviceConfig to map[string]string.
func (deviceConfig *DeviceConfig) ToMap() (m map[string]string, err error) {

	if deviceConfig.SecurityParameters == nil {
		return nil, fmt.Errorf("No security parameters")
	}

	m = make(map[string]string)
	m["version"] = deviceConfig.Version
	m["endpoint"] = deviceConfig.Endpoint
	m["port"] = fmt.Sprintf("%d", deviceConfig.Port)
	m["contextName"] = deviceConfig.ContextName

	securityParameters := deviceConfig.SecurityParameters
	m["userName"] = securityParameters.UserName
	if securityParameters.AuthenticationProtocol == MD5 {
		m["authenticationProtocol"] = "MD5"
	} else if securityParameters.AuthenticationProtocol == SHA {
		m["authenticationProtocol"] = "SHA"
	}
	m["authenticationPassphrase"] = securityParameters.AuthenticationPassphrase
	if securityParameters.PrivacyProtocol == DES {
		m["privacyProtocol"] = "DES"
	} else if securityParameters.PrivacyProtocol == AES {
		m["privacyProtocol"] = "AES"
	}
	m["privacyPassphrase"] = securityParameters.PrivacyPassphrase
	return m, nil
}

// SnmpClient is a thin wrapper around gosnmp.
type SnmpClient struct {
	DeviceConfig *DeviceConfig
}

// NewSnmpClient constructs SnmpClient.
func NewSnmpClient(deviceConfig *DeviceConfig) (*SnmpClient, error) {
	if deviceConfig == nil {
		return nil, fmt.Errorf("deviceConfig is nil")
	}

	return &SnmpClient{
		DeviceConfig: deviceConfig,
	}, nil
}

// ReadResult is the result structure for any SNMP read.
type ReadResult struct {
	Oid  string      // The SNMP OID read.
	Data interface{} // The data for the OID. See gosnmp decodeValue() https://github.com/soniah/gosnmp/blob/master/helper.go#L67
}

// Get performs an SNMP get on the given OID.
func (client *SnmpClient) Get(oid string) (result ReadResult, err error) {

	goSnmp, err := client.createGoSNMP()
	if err != nil {
		return result, err
	}

	oids := []string{oid}
	snmpPacket, err := goSnmp.Get(oids)
	err2 := goSnmp.Conn.Close() // Do not leak connection.

	// Return first error.
	if err != nil {
		return result, err
	}
	if err2 != nil {
		return result, err2
	}

	data := snmpPacket.Variables[0]

	// If it looks like an ASCII string, try to translate it.
	if data.Type == gosnmp.OctetString {
		ascii, err := TranslatePrintableASCII(data.Value)
		if err == nil {
			data.Value = ascii
		}
		// err above is deliberately ignored here. SNMP does not differentiate
		// between ASCII strings and byte array.
	}

	return ReadResult{
		Oid:  data.Name,
		Data: data.Value,
	}, nil
}

// Walk performs an SNMP bulk walk on the given OID.
func (client *SnmpClient) Walk(rootOid string) (results []ReadResult, err error) {

	goSnmp, err := client.createGoSNMP()
	if err != nil {
		return nil, err
	}

	resultSet, err := goSnmp.BulkWalkAll(rootOid)
	err2 := goSnmp.Conn.Close() // Do not leak connection.

	// Return first error.
	if err != nil {
		return nil, err
	}
	if err2 != nil {
		return nil, err2
	}

	// Package results.
	for _, snmpPdu := range resultSet {

		// If it looks like an ASCII string, try to translate it.
		if snmpPdu.Type == gosnmp.OctetString {
			ascii, err := TranslatePrintableASCII(snmpPdu.Value)
			if err == nil {
				snmpPdu.Value = ascii
			}
			// err above is deliberately ignored here. SNMP does not differentiate
			// between ASCII strings and byte array.
		}

		results = append(results, ReadResult{
			Oid:  snmpPdu.Name,
			Data: snmpPdu.Value,
		})
	}
	return results, nil
}

// createGoSNMP is a helper to create gosnmp.GoSNMP from SnmpClient.
// On success, the connection is open.
func (client *SnmpClient) createGoSNMP() (*gosnmp.GoSNMP, error) {

	// Argument checks
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	// Map DeviceConfig parameters to gosnmp parameters.
	securityParameters := client.DeviceConfig.SecurityParameters
	var authProtocol gosnmp.SnmpV3AuthProtocol
	var privProtocol gosnmp.SnmpV3PrivProtocol

	if securityParameters.AuthenticationProtocol == MD5 {
		authProtocol = gosnmp.MD5
	} else if securityParameters.AuthenticationProtocol == SHA {
		authProtocol = gosnmp.SHA
	} else {
		return nil, fmt.Errorf("Unsupported authentication protocol [%v]", securityParameters.AuthenticationProtocol)
	}

	if securityParameters.PrivacyProtocol == DES {
		privProtocol = gosnmp.DES
	} else if securityParameters.PrivacyProtocol == AES {
		privProtocol = gosnmp.AES
	} else {
		return nil, fmt.Errorf("Unsupported privacy protocol [%v]", securityParameters.PrivacyProtocol)
	}

	goSnmp := &gosnmp.GoSNMP{
		Target:        client.DeviceConfig.Endpoint,
		Port:          client.DeviceConfig.Port,
		Version:       gosnmp.Version3,
		Timeout:       client.DeviceConfig.Timeout,
		SecurityModel: gosnmp.UserSecurityModel,
		MsgFlags:      gosnmp.AuthPriv,
		SecurityParameters: &gosnmp.UsmSecurityParameters{
			UserName:                 client.DeviceConfig.SecurityParameters.UserName,
			AuthenticationProtocol:   authProtocol,
			AuthenticationPassphrase: client.DeviceConfig.SecurityParameters.AuthenticationPassphrase,
			PrivacyProtocol:          privProtocol,
			PrivacyPassphrase:        client.DeviceConfig.SecurityParameters.PrivacyPassphrase,
		},
		ContextName: client.DeviceConfig.ContextName,
	}

	// Connect
	err := goSnmp.Connect()
	if err != nil {
		log.Error("gosnmp failed to connect")
		return nil, fmt.Errorf("Failed to connect gosnmp: %+v", err)
	}
	return goSnmp, err
}
