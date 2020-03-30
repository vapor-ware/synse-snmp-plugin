package exp

import (
	"fmt"
	"github.com/vapor-ware/synse-sdk/sdk"
	"time"
)


// mibs is a package variable which is used to hold all of the MIBs that
// are registered with the SNMP base plugin.
// TODO (etd): need to figure out what this definition actually looks like.
var mibs = map[string]string{}


// TODO: figure out what a MIB actually looks like
func AddMib(name string, mib string) error {
	if _, exists := mibs[name]; exists {
		return fmt.Errorf("mib %s already exists", name)
	}
	mibs[name] = mib
	return nil
}



// PluginMetadata holds metadata for the plugin instance. It is used to
// provide identity to the plugin as well as some high level information
// about it and its source.
type PluginMetadata struct {
	Name string
	Maintainer string
	Description string
	VCS string
}

func customDeviceIdentifier(data map[string]interface{}) string {
	// TODO (etd): Do we need to include host information here too? That is to say,
	//	 can a single host (e.g. 1.2.3.4) expose data from multiple MIBs?
	oid, exists := data["oid"]
	if !exists {
		panic("unable to generate device identifier: OID does not exist in device data")
	}
	mibName, exists := data["mib"]
	if !exists {
		panic("unable to generate device identifier: MIB name does not exist in device data")
	}
	return fmt.Sprintf("%s:%s", mibName, oid)
}


type MIB struct {
	Name    string
	Devices []SNMPDevice
}

// This could be what gets converted into a proper SDK Device, so each subclassed plugin
// really only needs to create these, associate with a MIB, and register the MIB.=
type SNMPDevice struct {
	// Required
	OID          string
	Info         string
	Type         string
	Handler      string
	Output       string

	// Optional
	Tags         []*sdk.Tag
	Data         map[string]interface{} // does it need this? we could probably fill all of this in on out own
	Context      map[string]string
	Alias        string
	Transforms   []sdk.Transformer
	WriteTimeout time.Duration
}


/*
FOR REFERENCE:
This is the definition of a Device in the SDK

type Device struct {
	Type         string
	Info         string
	Tags         []*Tag
	Data         map[string]interface{}
	Context      map[string]string
	Handler      string
	SortIndex    int32
	Alias        string
	Transforms   []Transformer
	WriteTimeout time.Duration
	Output       string

	id           string
	idName       string
	handler      *DeviceHandler
}
 */

func customDynamicRegistration(data map[string]interface{}) ([]*sdk.Device, error) {
	// TODO:
	//  - load data from the `mibs` variable
	//  - look up the mib name from the current registration block and use that MIB
	//  - build devices for everything therein
	//  - return said devices


	return nil, nil
}

// NewSnmpBasePlugin creates a new SNMP base plugin.
//
// This base plugin can be used by other plugin implementations to inherit generic
// SNMP handling. Plugin implementations need only provide plugin metadata for the
// "subclassed" plugin and info mapping MIB devices to Synse devices.
//
//  TODO:  could pass in the device identifier? nah -- just require a device to be defined a certain way
//  TODO:  could pass in the dynamic function? possibly  - or could just have the base plugin check some
//		state for registered plugins and just dynamically register whatever is in there. Yeah.
func NewSnmpBasePlugin(metadata *PluginMetadata) (*sdk.Plugin, error) {
	sdk.SetPluginInfo(
		metadata.Name,
		metadata.Maintainer,
		metadata.Description,
		metadata.VCS,
	)

	plugin, err := sdk.NewPlugin(
		sdk.PluginConfigRequired(),
		sdk.DeviceConfigOptional(),
		sdk.CustomDeviceIdentifier(customDeviceIdentifier),
		sdk.CustomDynamicDeviceRegistration(customDynamicRegistration),
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
	err = plugin.RegisterDeviceHandlers(

	)
	if err != nil {
		return nil, err
	}

	return plugin, nil
}