package pkg

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/servers"
)

// mapOidsToInstances creates a map of SNMP OID to device instances and a list
// of OIDs so that Synse can determine the sort order for SNMP devices in a
// scan. In this case the OID is a string.
func mapOidsToInstances(deviceProtos []*config.DeviceProto) (oidMap map[string]*config.DeviceInstance, oidList []string, err error) {

	oidMap = map[string]*config.DeviceInstance{}
	// Iterate from the device config to each device instance.

	for i := 0; i < len(deviceProtos); i++ {
		proto := deviceProtos[i]

		for j := 0; j < len(proto.Instances); j++ {
			instance := proto.Instances[j]

			// Check for errors and add the oid as a string and a pointer to the
			// instance to the map value.
			oidData, ok := instance.Data["oid"]
			if !ok {
				return nil, []string{}, fmt.Errorf(
					"oid is not a key in instance data, instance.Data: %+v", instance.Data)
			}
			oidStr, ok := oidData.(string)
			if !ok {
				return nil, []string{}, fmt.Errorf("oid data is not a string, %T, %+v", oidData, oidData)
			}
			_, exists := oidMap[oidStr]
			if exists {
				return nil, []string{}, fmt.Errorf("oid %v already exists. Should not be duplicated", oidStr)
			}
			oidMap[oidStr] = instance
			oidList = append(oidList, oidStr)
		}
	}
	return oidMap, oidList, nil
}

// deviceIdentifier defines the SNMP-specific way of uniquely identifying a
// device through its device configuration.
// TODO: This will work for the initial cut. This may change later if/when
// we need to support the entity mib and entity sensor mib where joins may be
// required.
func deviceIdentifier(data map[string]interface{}) string {
	return fmt.Sprint(data["oid"])
}

// deviceEnumerator allows the sdk to enumerate devices.
func deviceEnumerator(data map[string]interface{}) (deviceConfigs []*config.DeviceProto, err error) {
	// Load the MIB from the configuration still.
	// Factory class for initializing servers via config is TODO:
	log.Info("[snmp] initializing UPS")
	pxgmsUps, err := servers.NewPxgmsUps(data)
	if err != nil {
		log.WithError(err).Error("Unable to initialize PxgmsUps")
		os.Exit(1)
	}
	log.Info("[snmp] UPS initialized")

	// First get a map of each OID to each device instance.
	oidMap, oidList, err := mapOidsToInstances(pxgmsUps.DeviceConfigs)
	if err != nil {
		log.WithError(err).Error("[snmp] failed mapping OIDs to instances")
		return nil, err
	}

	// Create an OidTrie and sort it.
	oidTrie, err := core.NewOidTrie(&oidList)
	if err != nil {
		log.WithError(err).Error("[snmp] failed to create OID Trie")
		return nil, err
	}
	sorted, err := oidTrie.Sort()
	if err != nil {
		log.WithError(err).Error("[snmp] failed to sort OID Trie")
		return nil, err
	}

	// Shim in the sort ordinal to the DeviceInstance Data.
	for ordinal := 0; ordinal < len(sorted); ordinal++ { // Zero based in list.
		oidMap[sorted[ordinal].ToString].SortIndex = int32(ordinal + 1) // One based sort ordinal.
	}

	// Dump SNMP device configurations.
	core.DumpDeviceConfigs(pxgmsUps.DeviceConfigs)
	return pxgmsUps.DeviceConfigs, nil
}
