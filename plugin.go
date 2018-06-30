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
	logger.Info("SNMP Plugin initializing UPS.")
	pxgmsUps, err := servers.NewPxgmsUps(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to create NewPxgmUps: %v", err)
	}
	logger.Infof("Initialized PxgmsUps: %+v\n", pxgmsUps)

	// Dump PxgmsUps device configurations.
	logger.Info("SNMP Plugin Dumping device configs")
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

	// Trace things out that the sdk is using for device enumeration.
	logger.Debugf("plugin: %+v", plugin)
	logger.Debugf("sdk.Config.Plugin: %+v", sdk.Config.Plugin)

	// Run the plugin.
	logger.Info("SNMP Plugin running plugin")
	if err := plugin.Run(); err != nil {
		logger.Fatalf("FATAL SNMP PLUGIN ERROR: %v", err)
	}
}
