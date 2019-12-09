package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/pkg"
)

const (
	pluginName       = "snmp"
	pluginMaintainer = "vaporio"
	pluginDesc       = "A general-purpose SNMP plugin"
	pluginVcs        = "https://github.com/vapor-ware/synse-snmp-plugin"
)

func main() {
	// Set the plugin metadata
	sdk.SetPluginInfo(
		pluginName,
		pluginMaintainer,
		pluginDesc,
		pluginVcs,
	)

	plugin := pkg.MakePlugin()

	// Run the plugin
	if err := plugin.Run(); err != nil {
		log.Fatal(err)
	}
}
