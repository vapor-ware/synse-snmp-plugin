package core

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
)

// This file contains utility functions. In the future we could put them in
// a "library" repo.

// CopyMapStringInterface returns a copy of the map passed in without error.
func CopyMapStringInterface(m map[string]interface{}) map[string]interface{} {
	target := make(map[string]interface{})
	for k, v := range m {
		target[k] = v
	}
	return target
}

// DumpDeviceConfigs to the log.
func DumpDeviceConfigs(deviceConfigs []*config.DeviceProto) {
	if deviceConfigs == nil {
		logger.Infof("No Device prototype Configs to dump\n")
		return
	}
	logger.Infof("Found %d device prototype configs.\n", len(deviceConfigs))
	for i := 0; i < len(deviceConfigs); i++ {
		logger.Infof("deviceProto[%d]: %T: %+v\n", i, deviceConfigs[i], deviceConfigs[i])
		//logger.Infof("deviceConfig[%d].Devices: %T: %+v", i, deviceConfigs[i].Devices, deviceConfigs[i].Devices)
		//devices := deviceConfigs[i].Devices
		instances := deviceConfigs[i].Instances
		for j := 0; j < len(instances); j++ {
			logger.Infof("deviceProto[%d].Instances[%d]: %T: %+v", i, j,
				instances[j], instances[j])

			logger.Infof("deviceProto[%d].Instances[%d].Output: %+v", i, j, instances[j].Output)
		}
	}
}

// MergeMapStringInterface returns a new map with the contents of both maps passed
// in. Errors out on duplicate keys.
func MergeMapStringInterface(a map[string]interface{}, b map[string]interface{}) (map[string]interface{}, error) {
	merged := CopyMapStringInterface(a)
	for k, v := range b {
		_, inMap := merged[k]
		if inMap {
			return nil, fmt.Errorf("Key %v already in merged map: %v", k, merged)
		}
		merged[k] = v
	}
	return merged, nil
}

// TranslatePrintableASCII translates byte arrays from gosnmp to a printable
// string if possible. If this call fails, the caller should normally just keep
// the raw byte array. This call makes no attempt to support extended (8bit)
// ASCII. We need this function since there is no differentiation between
// strings and byte arrays in the SNMP protocol.
func TranslatePrintableASCII(x interface{}) (string, error) {
	bytes, ok := x.([]uint8)
	if !ok {
		return "", fmt.Errorf("Failure converting type: %T, data: %v to byte array", x, x)
	}

	for i := 0; i < len(bytes); i++ {
		if !(bytes[i] < 0x80 && bytes[i] > 0x1f) {
			return "", fmt.Errorf("Unable to convert %x byte %x at %d to printable Ascii", bytes, bytes[i], i)
		}
	}
	return string(bytes), nil
}
