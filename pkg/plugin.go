package pkg

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/devices"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/outputs"
)

// PluginMetadata holds metadata for the plugin instance. It is used to
// provide identity to the plugin as well as some high level information
// about it and its source.
type PluginMetadata struct {
	Name string
	Maintainer string
	Description string
	VCS string
}

// NewSnmpPluginBase creates a new SNMP base plugin.
//
// This base plugin can be used by other plugin implementations to inherit generic
// SNMP handling. Plugin implementations need only provide plugin metadata for the
// "subclassed" plugin and info mapping MIB devices to Synse devices.
func NewSnmpPluginBase(metadata *PluginMetadata) (*sdk.Plugin, error) {
	sdk.SetPluginInfo(
		metadata.Name,
		metadata.Maintainer,
		metadata.Description,
		metadata.VCS,
	)

	plugin, err := sdk.NewPlugin(
		// TODO (etd): custom handler functions here???
		sdk.PluginConfigRequired(),
		sdk.DeviceConfigOptional(),
	)
	if err != nil {
		return nil, err
	}

	// Since this is a generic base, no custom output types are registered
	// to the plugin instance here. Plugins which use this as the base ar
	// free to add their own custom outputs once they have this generic base
	// plugin, e.g.
	//
	//   plugin, _ := NewSnmpPluginBase(...)
	//   err = plugin.RegisterOutputs(
	//       &customOutput,
	//   )


	// TODO (etd): Figure out how to register device handlers. Should there
	//   just be device handlers for the generic SNMP methods? Should each
	//   plugin subclass define their own device handlers? Should device handlers
	//   be generated from something that the plugin subclass gives this
	//   base constructor?

	return plugin, nil
}

// MakePlugin creates a new instance of the Synse SNMP Plugin.
func MakePlugin() *sdk.Plugin {
	plugin, err := sdk.NewPlugin(
		sdk.CustomDeviceIdentifier(deviceIdentifier),
		sdk.CustomDynamicDeviceConfigRegistration(deviceEnumerator),
		sdk.DeviceConfigOptional(),
		sdk.PluginConfigRequired(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register custom output types.
	err = plugin.RegisterOutputs(
		&outputs.Identity,
		&outputs.VAPower,
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
