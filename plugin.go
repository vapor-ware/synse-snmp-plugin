package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/vapor-ware/synse-sdk/sdk"
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
	log.Info("SNMP Plugin initializing UPS.")
	pxgmsUps, err := servers.NewPxgmsUps(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to create NewPxgmUps: %v", err)
	}
	log.Infof("Initialized PxgmsUps: %+v\n", pxgmsUps)

	// Dump PxgmsUps device configurations.
	log.Info("SNMP Plugin Dumping device configs")
	core.DumpDeviceConfigs(pxgmsUps.DeviceConfigs)
	return pxgmsUps.DeviceConfigs, nil
}

func main() {
	log.Info("SNMP Plugin start")
	// Set the plugin metadata
	sdk.SetPluginMeta(
		pluginName,
		pluginMaintainer,
		pluginDesc,
		pluginVcs,
	)

	log.Info("SNMP Plugin calling NewPlugin")
	plugin := sdk.NewPlugin(
		sdk.CustomDeviceIdentifier(deviceIdentifier),
		sdk.CustomDynamicDeviceConfigRegistration(deviceEnumerator),
	)

	// Register the supported output types
	log.Info("SNMP Plugin registering output types")
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
		log.Fatal(err)
	}

	// Register Device Handlers for all supported devices we interact with over SNMP.
	log.Info("SNMP Plugin registering device handlers")
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
	log.Debugf("plugin: %+v", plugin)
	log.Debugf("plugin.Config: %+v", plugin.Config)
	log.Debugf("plugin.Config.AutoEnumerate: %+v", plugin.Config.AutoEnumerate)

	// Run the plugin.
	log.Info("SNMP Plugin running plugin")
	if err := plugin.Run(); err != nil {
		log.Fatalf("FATAL SNMP PLUGIN ERROR: %v", err)
	}
}
