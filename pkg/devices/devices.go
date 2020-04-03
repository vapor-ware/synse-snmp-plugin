package devices

// This file contains device utility functions and common device functions.

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// SNMPDeviceHandlers holds a reference to all of the device handlers used by
// the SNMP plugin.
var SNMPDeviceHandlers = []*sdk.DeviceHandler{
	&SnmpCurrent,
	&SnmpFrequency,
	&SnmpIdentity,
	&SnmpPower,
	&SnmpStatus,
	&SnmpTemperature,
	&SnmpVoltage,
}

// Get the raw reading from the SNMP server with error checks.
// Factors out common code.
func getRawReading(device *sdk.Device) (result core.ReadResult, err error) {
	// Arg checks.
	if device == nil {
		return result, fmt.Errorf("device is nil")
	}

	// Get the SNMP device config from the strings in data.
	data := device.Data
	snmpConfig, err := core.GetDeviceConfig(data)
	if err != nil {
		return result, err
	}

	// Create SnmpClient.
	snmpClient, err := core.NewSnmpClient(snmpConfig)
	if err != nil {
		return result, err
	}

	// Read the SNMP OID in the device config.
	return snmpClient.Get(fmt.Sprint(data["oid"]))
}