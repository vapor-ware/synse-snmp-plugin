package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
)

// MIB is a logical grouping of SnmpDevices which a SNMP plugin implementation
// should define. A single plugin implementation may define multiple MIBs. These
// MIBs are registered with the SNMP base plugin.
type MIB struct {
	Name    string
	Devices []*SnmpDevice
}

// NewMIB creates a new MIB with the specified devices.
func NewMIB(name string, devices ...*SnmpDevice) *MIB {
	return &MIB{
		Name:    name,
		Devices: devices,
	}
}

// String returns a human-readable string, useful for identifying the
// MIB in logs.
func (mib *MIB) String() string {
	return fmt.Sprintf("[MIB %s]", mib.Name)
}

func (mib *MIB) LoadDevices() ([]*sdk.Device, error) {
	// TODO
	return nil, nil
}
