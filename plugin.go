package main

import (
	"fmt"
	"log"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-sdk/sdk/logger"
	"github.com/vapor-ware/synse-snmp-plugin/devices"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/servers"
)

// Build time variables for setting the version info of a Plugin.
var (
	BuildDate     string
	GitCommit     string
	GitTag        string
	GoVersion     string
	VersionString string
)

// DeviceIdentifier defines the SNMP-specific way of uniquely identifying a
// device through its device configuration.
// TODO: This will work for the initial cut. This may change later if/when
// we need to support the entity mib and entity sensor mib where joins may be
// required.
func DeviceIdentifier(data map[string]string) string {
	return data["oid"]
}

// DeviceEnumerator allows the sdk to enumerate devices.
func DeviceEnumerator(data map[string]interface{}) (deviceConfigs []*config.DeviceConfig, err error) {
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
	logger.Info("SNMP Plugin start")

	logger.Info("SNMP Plugin initializing handlers")
	handlers, err := sdk.NewHandlers(DeviceIdentifier, DeviceEnumerator)
	if err != nil {
		log.Fatalf("FATAL SNMP PLUGIN ERROR (NewHandlers): %v", err)
	}

	logger.Info("SNMP Plugin calling NewPlugin")
	plugin, err := sdk.NewPlugin(handlers, nil)
	if err != nil {
		log.Fatalf("FATAL SNMP PLUGIN ERROR (NewPlugin): %v", err)
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

	// Set build-time version info.
	logger.Info("SNMP Plugin setting version")
	plugin.SetVersion(sdk.VersionInfo{
		BuildDate:     BuildDate,
		GitCommit:     GitCommit,
		GitTag:        GitTag,
		GoVersion:     GoVersion,
		VersionString: VersionString,
	})

	// Trace things out that the sdk is using for device enumeration.
	logger.Debugf("plugin: %+v", plugin)
	logger.Debugf("plugin.Config: %+v", plugin.Config)
	logger.Debugf("plugin.Config.AutoEnumerate: %+v", plugin.Config.AutoEnumerate)

	// Run the plugin.
	logger.Info("SNMP Plugin running plugin")
	err = plugin.Run()
	if err != nil {
		log.Fatalf("FATAL SNMP PLUGIN ERROR: %v", err)
	}
}
