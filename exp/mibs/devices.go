package mibs

import (
	"fmt"
	"time"

	"github.com/iancoleman/strcase"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
)

// SnmpDevice models a device and serves as the bridge between SNMP configuration
// and Synse device configuration.
//
// Each SNMP plugin implementation which uses the plugin base should define these
// devices in the code base. A collection of SnmpDevices makes up a MIB. These
// devices are used to generate Synse devices on plugin init.
type SnmpDevice struct {
	// Required fields
	OID     string
	Info    string
	Type    string
	Handler string
	Output  string

	// Optional fields
	Tags         []*sdk.Tag
	Data         map[string]interface{}
	Context      map[string]string
	Alias        string
	Transforms   []sdk.Transformer
	WriteTimeout time.Duration
}

// String returns a human-readable string, useful for identifying the
// device in logs.
func (device *SnmpDevice) String() string {
	return fmt.Sprintf("[SnmpDevice %s: %s]", device.OID, device.Info)
}

// ToDevice converts the plugin-specific SnmpDevice to a Synse SDK Device.
func (device *SnmpDevice) ToDevice() (*sdk.Device, error) {
	log.WithFields(log.Fields{
		"oid":  device.OID,
		"info": device.Info,
	}).Debug("[snmp] creating synse device from MIB device")
	// Construct the device data.
	data := map[string]interface{}{}
	for k, v := range device.Data {
		data[k] = v
	}
	// Note: this will be augmented with MIB and Agent info later
	// (via MIB.LoadDevices)
	data["oid"] = device.OID

	// Construct the device context.
	context := map[string]string{}
	for k, v := range device.Context {
		context[k] = v
	}
	context["oid"] = device.OID

	// Ensure that the device info can be made into a tag. Eliminate any spaces
	// which may be present in the string.
	normalizedInfo := strcase.ToLowerCamel(device.Info)

	// Create a set of standard tags for the device.
	tags := []*sdk.Tag{
		core.TagOrPanic("protocol/snmp"),
		core.TagOrPanic(fmt.Sprintf("snmp/oid:%s", device.OID)),
		core.TagOrPanic(fmt.Sprintf("snmp/name:%s", normalizedInfo)),
	}
	tags = append(tags, device.Tags...)

	// Check that the output exists with the plugin (either as an SDK
	// built-in or as a custom plugin defined output).
	if o := output.Get(device.Output); o == nil {
		return nil, fmt.Errorf("unable to create synse device: output '%s' not registered", device.Output)
	}

	return &sdk.Device{
		Type:         device.Type,
		Info:         device.Info,
		Handler:      device.Handler,
		Alias:        device.Alias,
		Output:       device.Output,
		WriteTimeout: device.WriteTimeout,
		Transforms:   device.Transforms,
		Tags:         tags,
		Data:         data,
		Context:      context,
	}, nil
}
