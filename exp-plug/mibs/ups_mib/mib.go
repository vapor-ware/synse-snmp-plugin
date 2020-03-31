package ups_mib

import (
	"github.com/vapor-ware/synse-snmp-plugin/exp-plug/mibs/ups_mib/devices"
	"github.com/vapor-ware/synse-snmp-plugin/exp/mibs"
)

// Mib is the MIB definition for the UPS MIB.
//
// See also: http://www.oidview.com/mibs/0/UPS-MIB.html
var Mib = mibs.NewMIB(
	"UPS-MIB",
	devices.IdentityDevices...,
)
