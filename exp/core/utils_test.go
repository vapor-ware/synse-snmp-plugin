package core

import (
	"testing"

	"github.com/soniah/gosnmp"
	"github.com/stretchr/testify/assert"
)

func TestTagOrPanic(t *testing.T) {
	tag := TagOrPanic("vapor/test:tag")
	assert.Equal(t, "vapor", tag.Namespace)
	assert.Equal(t, "test", tag.Annotation)
	assert.Equal(t, "tag", tag.Label)
}

func TestTagOrPanic_Panics(t *testing.T) {
	assert.Panics(t, func() {
		_ = TagOrPanic("this is not a correct tag")
	})
}

func TestGetSNMPVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected gosnmp.SnmpVersion
	}{
		{
			name:     "version 1 (lowercase)",
			version:  "v1",
			expected: gosnmp.Version1,
		},
		{
			name:     "version 1 (uppercase)",
			version:  "V1",
			expected: gosnmp.Version1,
		},
		{
			name:     "version 2 (lowercase)",
			version:  "v2",
			expected: gosnmp.Version2c,
		},
		{
			name:     "version 2 (uppercase)",
			version:  "V2",
			expected: gosnmp.Version2c,
		},
		{
			name:     "version 2c (lowercase)",
			version:  "v2c",
			expected: gosnmp.Version2c,
		},
		{
			name:     "version 2c (uppercase)",
			version:  "V2c",
			expected: gosnmp.Version2c,
		},
		{
			name:     "version 3 (lowercase)",
			version:  "v3",
			expected: gosnmp.Version3,
		},
		{
			name:     "version 3 (uppercase)",
			version:  "V3",
			expected: gosnmp.Version3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := GetSNMPVersion(test.version)
			assert.Equal(t, test.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestGetSNMPVersion_Error(t *testing.T) {
	actual, err := GetSNMPVersion("v99")
	assert.Equal(t, gosnmp.SnmpVersion(0), actual)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidSNMPVersion, err)
}

func TestGetPrivProtocol(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
		expected gosnmp.SnmpV3PrivProtocol
	}{
		{
			name:     "AES (lowercase)",
			protocol: "aes",
			expected: gosnmp.AES,
		},
		{
			name:     "AES (uppercase)",
			protocol: "AES",
			expected: gosnmp.AES,
		},
		{
			name:     "DES (lowercase)",
			protocol: "des",
			expected: gosnmp.DES,
		},
		{
			name:     "DES (uppercase)",
			protocol: "DES",
			expected: gosnmp.DES,
		},
		{
			name:     "None (lowercase)",
			protocol: "none",
			expected: gosnmp.NoPriv,
		},
		{
			name:     "None (uppercase)",
			protocol: "NONE",
			expected: gosnmp.NoPriv,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := GetPrivProtocol(test.protocol)
			assert.Equal(t, test.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestGetPrivProtocol_Error(t *testing.T) {
	actual, err := GetPrivProtocol("unsupported value")
	assert.Equal(t, gosnmp.SnmpV3PrivProtocol(0), actual)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPrivProtocol, err)
}

func TestGetAuthProtocol(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
		expected gosnmp.SnmpV3AuthProtocol
	}{
		{
			name:     "MD5 (lowercase)",
			protocol: "md5",
			expected: gosnmp.MD5,
		},
		{
			name:     "MD5 (uppercase)",
			protocol: "MD5",
			expected: gosnmp.MD5,
		},
		{
			name:     "SHA (lowercase)",
			protocol: "sha",
			expected: gosnmp.SHA,
		},
		{
			name:     "SHA (uppercase)",
			protocol: "SHA",
			expected: gosnmp.SHA,
		},
		{
			name:     "None (lowercase)",
			protocol: "none",
			expected: gosnmp.NoAuth,
		},
		{
			name:     "None (uppercase)",
			protocol: "NONE",
			expected: gosnmp.NoAuth,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := GetAuthProtocol(test.protocol)
			assert.Equal(t, test.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestGetAuthProtocol_Error(t *testing.T) {
	actual, err := GetAuthProtocol("unsupported value")
	assert.Equal(t, gosnmp.SnmpV3AuthProtocol(0), actual)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAuthProtocol, err)
}

func TestGetSecurityFlags(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected gosnmp.SnmpV3MsgFlags
	}{
		{
			name:     "No auth, no priv (lowercase)",
			level:    "noauthnopriv",
			expected: gosnmp.NoAuthNoPriv,
		},
		{
			name:     "No auth, no priv (mixed case)",
			level:    "NoAuthNoPriv",
			expected: gosnmp.NoAuthNoPriv,
		},
		{
			name:     "Auth, no priv (lowercase)",
			level:    "authnopriv",
			expected: gosnmp.AuthNoPriv,
		},
		{
			name:     "Auth, no priv (mixed case)",
			level:    "AuthNoPriv",
			expected: gosnmp.AuthNoPriv,
		},
		{
			name:     "Auth, Priv (lowercase)",
			level:    "authpriv",
			expected: gosnmp.AuthPriv,
		},
		{
			name:     "Auth, Priv (mixed case)",
			level:    "AuthPriv",
			expected: gosnmp.AuthPriv,
		},
		{
			name:     "Reportable (lowercase)",
			level:    "reportable",
			expected: gosnmp.Reportable,
		},
		{
			name:     "Reportable (mixed case)",
			level:    "Reportable",
			expected: gosnmp.Reportable,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := GetSecurityFlags(test.level)
			assert.Equal(t, test.expected, actual)
			assert.NoError(t, err)
		})
	}
}

func TestGetSecurityFlags_Error(t *testing.T) {
	actual, err := GetSecurityFlags("unsupported value")
	assert.Equal(t, gosnmp.SnmpV3MsgFlags(0), actual)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidMessageFlag, err)
}

func TestBytesIfaceToASCII(t *testing.T) {
	data := []uint8{0x68, 0x65, 0x6c, 0x6c, 0x6f}

	str, err := BytesIfaceToASCII(data)
	assert.NoError(t, err)
	assert.Equal(t, "hello", str)
}

func TestBytesIfaceToASCII_BadIface(t *testing.T) {
	data := "not a byte array"

	str, err := BytesIfaceToASCII(data)
	assert.Error(t, err)
	assert.Equal(t, "", str)
}

func TestBytesIfaceToASCII_UnableToConvert(t *testing.T) {
	data := []uint8{0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x00}

	str, err := BytesIfaceToASCII(data)
	assert.Error(t, err)
	assert.Equal(t, "", str)
}
