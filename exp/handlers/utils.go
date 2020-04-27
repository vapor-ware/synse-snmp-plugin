package handlers

import (
	"fmt"

	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
)

// getAgent is a convenience function to safely get the "agent" value out of a device's
// Data field and cast it to the appropriate type.
//
// Since the "agent" field is expected to exist in the device Data and it is expected to
// be a string, this function returns an error if it does not exist or cannot be cast to
// a string.
func getAgent(data map[string]interface{}) (string, error) {
	agentIface, exists := data["agent"]
	if !exists {
		return "", fmt.Errorf("expected field 'agent' in device data, but not found")
	}
	agent, ok := agentIface.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast 'agent' value (%T) to string", agentIface)
	}
	return agent, nil
}

// getOid is a convenience function to safely get the "oid" value out of a device's
// Data field and cast it to the appropriate type.
//
// Since the "oid" field is expected to exist in the device Data and it is expected to
// be a string, this function returns an error if it does not exist or cannot be cast to
// a string.
func getOid(data map[string]interface{}) (string, error) {
	oidIface, exists := data["oid"]
	if !exists {
		return "", fmt.Errorf("expected field 'oid' in device data, but not found")
	}
	oid, ok := oidIface.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast 'oid' value (%T) to string", oidIface)
	}
	return oid, nil
}

// getTargetConfig is a convenience function to safely get the "target_cfg" value out of a device's
// Data field and cast it to the appropriate type.
//
// Since the "target_cfg" field is expected to exist in the device Data and it is expected to
// be an SnmpTargetConfiguration, this function returns an error if it does not exist or cannot
// be cast to an SnmpTargetConfiguration.
func getTargetConfig(data map[string]interface{}) (*core.SnmpTargetConfiguration, error) {
	cfgIface, exists := data["target_cfg"]
	if !exists {
		return nil, fmt.Errorf("expected field 'target_cfg' in device data, but not found")
	}
	cfg, ok := cfgIface.(*core.SnmpTargetConfiguration)
	if !ok {
		return nil, fmt.Errorf("failed to cast 'target_cfg' value (%T) to SnmpTargetConfiguration", cfgIface)
	}
	return cfg, nil
}
