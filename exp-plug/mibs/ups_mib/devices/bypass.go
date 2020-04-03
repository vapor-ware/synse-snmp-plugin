package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// BypassDevices contains the definitions of all the "upsBypass" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.5
// http://www.oidview.com/mibs/0/UPS-MIB.html
var BypassDevices = []*mibs.SnmpDevice{}

// UPS-MIB upsBypass device definitions.
var ()
