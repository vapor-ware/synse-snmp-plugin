package exp

import (
	"github.com/pkg/errors"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/exp/handlers"
)

// Errors for base plugin setup/creation.
var (
	ErrNoName       = errors.New("plugin metadata does not specify the required 'Name' field")
	ErrNoMaintainer = errors.New("plugin metadata does not specify the required 'Maintainer' field")
)

// PluginMetadata holds metadata for the plugin instance. It is used to
// provide identity to the plugin as well as some high level information
// about it and its source.
type PluginMetadata struct {
	Name        string
	Maintainer  string
	Description string
	VCS         string
}

// NewSnmpBasePlugin creates a new SNMP base plugin.
//
// This base plugin can be used by other plugin implementations to inherit generic
// SNMP handling. Plugin implementations need only provide plugin metadata for the
// "subclassed" plugin and info mapping MIB devices to Synse devices.
func NewSnmpBasePlugin(metadata *PluginMetadata) (*sdk.Plugin, error) {
	if metadata.Name == "" {
		return nil, ErrNoName
	}
	if metadata.Maintainer == "" {
		return nil, ErrNoMaintainer
	}

	sdk.SetPluginInfo(
		metadata.Name,
		metadata.Maintainer,
		metadata.Description,
		metadata.VCS,
	)

	plugin, err := sdk.NewPlugin(
		sdk.PluginConfigRequired(),
		sdk.DynamicConfigRequired(),
		sdk.DeviceConfigOptional(),
		sdk.CustomDeviceIdentifier(SnmpDeviceIdentifier),
		sdk.CustomDynamicDeviceRegistration(SnmpDeviceRegistrar),
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

	err = plugin.RegisterDeviceHandlers(
		&handlers.ReadOnly,
		&handlers.ReadWrite,
		&handlers.WriteOnly,
	)
	if err != nil {
		return nil, err
	}

	return plugin, nil
}
