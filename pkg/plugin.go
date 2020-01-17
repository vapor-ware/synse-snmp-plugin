package pkg

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/devices"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/outputs"
)

// MakePlugin creates a new instance of the Synse SNMP Plugin.
func MakePlugin() *sdk.Plugin {
	plugin, err := sdk.NewPlugin(
		sdk.CustomDeviceIdentifier(deviceIdentifier),
		sdk.CustomDynamicDeviceConfigRegistration(deviceEnumerator),
		sdk.DeviceConfigOptional(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register custom output types.
	err = plugin.RegisterOutputs(
		&outputs.Identity,
		&outputs.VAPower,
		&outputs.WattsPower,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register device handlers
	err = plugin.RegisterDeviceHandlers(
		devices.SNMPDeviceHandlers...,
	)
	if err != nil {
		log.Fatal(err)
	}

	return plugin
}
