package exp

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
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
	host, exists := data["host"]
	if !exists {
		panic("unable to generate device ID: 'host' not found in device data")
	}
	port, exists := data["port"]
	if !exists {
		panic("unable to generate device ID: 'port' not found in device data")
	}

	return fmt.Sprintf("%s:%v-%s:%s", host, port, mibName, oid)
}

// SnmpDeviceRegistrar is the dynamic registration function used by the SDK to
// build devices at runtime.
//
// This function is defined for thee base SNMP plugin and is subsequently used
// by all plugins which use the base.
func SnmpDeviceRegistrar(data map[string]interface{}) ([]*sdk.Device, error) {
	// TODO (etd): validate config

	// TODO (etd): implement
	//  - load data from the `mibs` variable
	//  - look up the mib name from the current registration block and use that MIB
	//  - build devices for everything therein
	//  - return said devices
	return nil, nil
}
