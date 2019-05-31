package core

import (
	"fmt"
)

// SnmpServerBase is a base class for specific SnmpServer implementations.
type SnmpServerBase struct {
	SnmpClient   *SnmpClient
	DeviceConfig *DeviceConfig
}

// NewSnmpServerBase constructs common code for all SNMP Servers.
func NewSnmpServerBase(
	client *SnmpClient, deviceConfig *DeviceConfig) (*SnmpServerBase, error) {
	// Parameter checks.
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	if deviceConfig == nil {
		return nil, fmt.Errorf("deviceConfig is nil")
	}

	// Construct the struct.
	return &SnmpServerBase{
		SnmpClient:   client,
		DeviceConfig: deviceConfig,
	}, nil
}
