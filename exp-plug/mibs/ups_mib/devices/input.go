package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// InputDevices contains the definitions of all the "upsInput" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.3
// http://www.oidview.com/mibs/0/UPS-MIB.html
var InputDevices = []*mibs.SnmpDevice{}

// UPS-MIB upsInput device definitions.
var ()
