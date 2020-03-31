package exp_plug

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/exp"
	"github.com/vapor-ware/synse-snmp-plugin/exp-plug/mibs/ups_mib"
	"github.com/vapor-ware/synse-snmp-plugin/exp-plug/outputs"
	"github.com/vapor-ware/synse-snmp-plugin/exp/mibs"
)

// MakePlugin creates a new instance of the plugin.
//
// It ensures all plugin-specific configuration is done and all necessary items
// are registered with the plugin.
func MakePlugin() *sdk.Plugin {

	// Build a base SNMP plugin instance.
	plugin, err := exp.NewSnmpBasePlugin(&exp.PluginMetadata{
		Name:        "experimental plugin",
		Maintainer:  "vapor-exp",
		Description: "plugin to test experimental plugin 'inheritance'",
		VCS:         "n/a",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Register custom output types.
	err = plugin.RegisterOutputs(
		&outputs.Identity,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register plugin-defined SNMP MIBs.
	err = mibs.Register(
		ups_mib.Mib,
	)
	if err != nil {
		log.Fatal(err)
	}

	return plugin
}
