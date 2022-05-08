package core

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckPrivacyAndAuth(t *testing.T) {
	cases := []struct {
		name     string
		data     map[string]interface{}
		config   *DeviceConfig
		expected gosnmp.SnmpV3MsgFlags
	}{
		{
			name: "No auth, no priv",
			data: map[string]interface{}{
				"authenticationProtocol": "None",
				"privacyProtocol":        "None",
			},
			config:   &DeviceConfig{},
			expected: gosnmp.NoAuthNoPriv,
		},
		{
			name: "Auth, no priv",
			data: map[string]interface{}{
				"authenticationProtocol": "MD5",
				"privacyProtocol":        "None",
			},
			config:   &DeviceConfig{},
			expected: gosnmp.AuthNoPriv,
		},
		{
			name: "No auth, priv",
			data: map[string]interface{}{
				"authenticationProtocol": "None",
				"privacyProtocol":        "AES",
			},
			config:   &DeviceConfig{},
			expected: gosnmp.AuthPriv, // gosnmp does not support NoAuthPriv
		},
		{
			name: "Auth and priv",
			data: map[string]interface{}{
				"authenticationProtocol": "SHA",
				"privacyProtocol":        "AES",
			},
			config:   &DeviceConfig{},
			expected: gosnmp.AuthPriv,
		},
	}

	for _, tc := range cases {
		err := tc.config.CheckPrivacyAndAuthFromData(tc.data)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, tc.config.MsgFlag, fmt.Sprintf("case %v: expected %v got %v", tc.name, tc.expected, tc.config.MsgFlag))
	}
}
