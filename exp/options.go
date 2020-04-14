package exp

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
	"github.com/vapor-ware/synse-snmp-plugin/exp/mibs"
)

// SnmpDeviceIdentifier is the device identifier function used by the SDK
// in order to generate unique device IDs for the registered Synse devices.
//
// This function is defined for the base SNMP plugin and is subsequently used
// by all plugins which use the base.
//
// The expectation is that each device should be uniquely identifiable using a
// combination of the SNMP device's OID and MIB name. As such, those fields are
// expected in the device Data. If they are not present, the plugin will panic
// and terminate.
//
// Additionally, since there may be multiple SNMP-enabled servers configured
// which use the same MIB, the id generation needs to take the configured host/port
// into account.
//
// Generally, it is not the responsibility of the plugin writer to ensure that
// info exists per device because the SNMP base plugin provides utility functions
// which automatically fill this information in when building Synse devices.
func SnmpDeviceIdentifier(data map[string]interface{}) string {
	oid, exists := data["oid"]
	if !exists {
		panic("unable to generate device ID: 'oid' not found in device data")
	}
	mibName, exists := data["mib"]
	if !exists {
		panic("unable to generate device ID: 'mib' not found in device data")
	}
	agent, exists := data["agent"]
	if !exists {
		panic("unable to generate device ID: 'agent' not found in device data")
	}

	return fmt.Sprintf("%v-%s:%s", agent, mibName, oid)
}

// SnmpDeviceRegistrar is the dynamic registration function used by the SDK to
// build devices at runtime.
//
// This function is defined for thee base SNMP plugin and is subsequently used
// by all plugins which use the base.
func SnmpDeviceRegistrar(data map[string]interface{}) ([]*sdk.Device, error) {
	// Load the data into a configurations struct.
	config, err := core.LoadTargetConfiguration(data)
	if err != nil {
		return nil, err
	}

	// Verify that the configured agent as a MIB name specified. This is
	// required, as it determines which devices will be loaded for the agent.
	if config.MIB == "" {
		return nil, fmt.Errorf("invalid configuration: no MIB specified for agent %s", config.Agent)
	}

	// Create a new client for the configured SNMP agent. This client is cached
	// so devices can re-use the client.
	client, err := core.NewClient(config)
	if err != nil {
		return nil, err
	}
	core.CacheClient(client)

	// Get the specified MIB and load its devices for the agent.
	mib := mibs.Get(config.MIB)
	d, err := mib.LoadDevices(client)
	if err != nil {
		log.WithError(err).Error("[snmp failed to load devices from MIB")
		return nil, err
	}
	return d, nil
}
