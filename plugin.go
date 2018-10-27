package main

import (
	"fmt"

	logger "github.com/Sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/policies"
	"github.com/vapor-ware/synse-snmp-plugin/devices"
	"github.com/vapor-ware/synse-snmp-plugin/outputs"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/servers"
)

const (
	pluginName       = "snmp"
	pluginMaintainer = "vaporio"
	pluginDesc       = "A general-purpose SNMP plugin"
	pluginVcs        = "https://github.com/vapor-ware/synse-snmp-plugin"
)

// mapOidsToInstances creates a map of SNMP OID to device instances and a list
// of OIDs so that Synse can determine the sort order for SNMP devices in a
// scan. In this case the OID is a string.
func mapOidsToInstances(deviceConfigs []*sdk.DeviceConfig) (oidMap map[string]*sdk.DeviceInstance, oidList []string, err error) {

	oidMap = map[string]*sdk.DeviceInstance{}
	// Iterate from the device config to each device instance.
	for i := 0; i < len(deviceConfigs); i++ {
		devices := deviceConfigs[i].Devices
		for j := 0; j < len(devices); j++ {
			device := devices[j]
			for k := 0; k < len(device.Instances); k++ {
				instance := device.Instances[k]

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
func deviceEnumerator(data map[string]interface{}) (deviceConfigs []*sdk.DeviceConfig, err error) {
	// Load the MIB from the configuration still.
	// Factory class for initializing servers via config is TODO:
	logger.Info("SNMP Plugin initializing UPS.")
	pxgmsUps, err := servers.NewPxgmsUps(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to create NewPxgmUps: %v", err)
	}
	logger.Infof("Initialized PxgmsUps: %+v\n", pxgmsUps)

	// First get a map of each OID to each device instance.
	oidMap, oidList, err := mapOidsToInstances(pxgmsUps.DeviceConfigs)
	if err != nil {
		return nil, err
	}

	// Create an OidTrie and sort it.
	oidTrie, err := core.NewOidTrie(&oidList)
	if err != nil {
		return nil, err
	}
	sorted, err := oidTrie.Sort()
	if err != nil {
		return nil, err
	}

	// Shim in the sort ordinal to the DeviceInstance Data.
	for ordinal := 0; ordinal < len(sorted); ordinal++ { // Zero based in list.
		oidMap[sorted[ordinal].ToString].SortOrdinal = int32(ordinal + 1) // One based sort ordinal.
	}

	// Dump SNMP device configurations.
	core.DumpDeviceConfigs(pxgmsUps.DeviceConfigs)
	return pxgmsUps.DeviceConfigs, nil
}

func main() {
	logger.SetLevel(logger.DebugLevel)
	logger.Info("SNMP Plugin start")
	// Set the plugin metadata
	sdk.SetPluginMeta(
		pluginName,
		pluginMaintainer,
		pluginDesc,
		pluginVcs,
	)

	// Set the device config file policy to optional. That way, it it doesn't
	// exist, we are okay (since we should then get some kind of config dynamically).
	logger.Info("SNMP Plugin - setting policies")
	policies.Add(policies.DeviceConfigFileOptional)

	logger.Info("SNMP Plugin calling NewPlugin")
	plugin := sdk.NewPlugin(
		sdk.CustomDeviceIdentifier(deviceIdentifier),
		sdk.CustomDynamicDeviceConfigRegistration(deviceEnumerator),
	)

	// Register the supported output types
	logger.Info("SNMP Plugin registering output types")
	err := plugin.RegisterOutputTypes(
		&outputs.Current,
		&outputs.Frequency,
		&outputs.Identity,
		&outputs.VAPower,
		&outputs.WattsPower,
		&outputs.Status,
		&outputs.Temperature,
		&outputs.Voltage,
	)
	if err != nil {
		logger.Fatal(err)
	}

	// Register Device Handlers for all supported devices we interact with over SNMP.
	logger.Info("SNMP Plugin registering device handlers")
	plugin.RegisterDeviceHandlers(
		&devices.SnmpCurrent,
		&devices.SnmpFrequency,
		&devices.SnmpIdentity,
		&devices.SnmpPower,
		&devices.SnmpStatus,
		&devices.SnmpTemperature,
		&devices.SnmpVoltage,
	)

	// Run the plugin.
	logger.Info("SNMP Plugin running plugin")
	if err := plugin.Run(); err != nil {
		logger.Fatalf("FATAL SNMP PLUGIN ERROR: %v", err)
	}
}
