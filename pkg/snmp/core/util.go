package core

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-sdk/sdk/utils"
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
func DumpDeviceConfigs(deviceConfigs []*config.DeviceProto) (err error) {

	var redactedContext interface{}
	var redactedData interface{}

	if deviceConfigs == nil {
		log.Infof("[snmp] no device prototype configs to dump")
		return
	}

	log.WithField("count", len(deviceConfigs)).Info("[snmp] found device prototype configs")

	for i := 0; i < len(deviceConfigs); i++ {
		proto := deviceConfigs[i]

		redactedContext, err = utils.RedactPasswords(proto.Context)
		if err != nil {
			return
		}
		redactedData, err = utils.RedactPasswords(proto.Data)
		if err != nil {
			return
		}

		log.WithFields(log.Fields{
			"idx":           i,
			"instances":     len(proto.Instances),
			"transforms":    len(proto.Transforms),
			"write timeout": proto.WriteTimeout,
			"handler":       proto.Handler,
			"tags":          proto.Tags,
			"type":          proto.Type,
			"context":       redactedContext,
			"data":          redactedData,
		}).Info("[snmp] dumping device prototype config")

		for j := 0; j < len(proto.Instances); j++ {
			instance := proto.Instances[j]
			redactedContext, err = utils.RedactPasswords(instance.Context)
			if err != nil {
				return
			}
			redactedData, err = utils.RedactPasswords(instance.Data)
			if err != nil {
				return
			}
			log.WithFields(log.Fields{
				"idx":           j,
				"prototype idx": i,
				"write timeout": instance.WriteTimeout,
				"transforms":    len(instance.Transforms),
				"type":          instance.Type,
				"tags":          instance.Tags,
				"handler":       instance.Handler,
				"info":          instance.Info,
				"alias":         instance.Alias,
				"output":        instance.Output,
				"context":       redactedContext,
				"data":          redactedData,
			}).Info("[snmp] dumping device instance config")
		}
	}
	return
}

// MergeMapStringInterface returns a new map with the contents of both maps passed
// in. Errors out on duplicate keys.
func MergeMapStringInterface(a map[string]interface{}, b map[string]interface{}) (map[string]interface{}, error) {
	merged := CopyMapStringInterface(a)
	for k, v := range b {
		_, inMap := merged[k]
		if inMap {
			return nil, fmt.Errorf("key %v already in merged map: %v", k, merged)
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
		return "", fmt.Errorf("failure converting type: %T, data: %v to byte array", x, x)
	}

	for i := 0; i < len(bytes); i++ {
		if !(bytes[i] < 0x80 && bytes[i] > 0x1f) {
			return "", fmt.Errorf("unable to convert %x byte %x at %d to printable Ascii", bytes, bytes[i], i)
		}
	}
	return string(bytes), nil
}
