package mibs

import (
	"fmt"
	"time"

	"github.com/vapor-ware/synse-sdk/sdk"
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

func (device *SnmpDevice) ToDevice() (*sdk.Device, error) {
	// TODO
	return &sdk.Device{}, nil
}
