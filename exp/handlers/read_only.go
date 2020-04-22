package handlers

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
)

// ReadOnly
var ReadOnly = sdk.DeviceHandler{
	Name: "read-only",
	Read: func(device *sdk.Device) ([]*output.Reading, error) {
		if device == nil {
			return nil, errors.New("unable to read from nil device")
		}

		// Get data cached in device.Data
		agent, err := getAgent(device.Data)
		if err != nil {
			return nil, err
		}
		oid, err := getOid(device.Data)
		if err != nil {
			return nil, err
		}
		targetConfig, err := getTargetConfig(device.Data)
		if err != nil {
			return nil, err
		}

		// Create a new client with the target configuration.
		c, err := core.NewClient(targetConfig)
		if err != nil {
			return nil, err
		}
		defer c.Close()

		log.WithFields(log.Fields{
			"agent": agent,
			"oid":   oid,
		}).Debug("[snmp] reading OID")

		result, err := c.GetOid(oid)
		if err != nil {
			return nil, err
		}

		log.WithFields(log.Fields{
			"value": result.Value,
			"name":  result.Name,
			"type":  result.Type,
		}).Debug("[snmp] got reading value for OID")

		var value interface{}
		switch result.Type {
		case gosnmp.OctetString:
			ascii, err := TranslatePrintableASCII(result.Value)
			if err != nil {
				return nil, err
			}
			value = ascii
		default:
			value = result.Value
		}

		log.WithFields(log.Fields{
			"value": value,
		}).Debug("[snmp] final value")

		o := output.Get(device.Output)
		if o == nil {
			return nil, fmt.Errorf("unable to format reading: device output not defined")
		}

		return []*output.Reading{
			o.MakeReading(value),
		}, nil
	},
}

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

// TranslatePrintableASCII translates byte arrays from gosnmp to a printable
// string if possible. If this call fails, the caller should normally just keep
// the raw byte array. This call makes no attempt to support extended (8bit)
// ASCII. We need this function since there is no differentiation between
// strings and byte arrays in the SNMP protocol.
func TranslatePrintableASCII(x interface{}) (string, error) {
	bytes, ok := x.([]uint8)
	if !ok {
		return "", fmt.Errorf("failure converting type: %T, data: %v to byte array", x, x)
	}

	for i := 0; i < len(bytes); i++ {
		if !(bytes[i] < 0x80 && bytes[i] > 0x1f) {
			return "", fmt.Errorf("unable to convert %x byte %x at %d to printable Ascii", bytes, bytes[i], i)
		}
	}
	return string(bytes), nil
}
