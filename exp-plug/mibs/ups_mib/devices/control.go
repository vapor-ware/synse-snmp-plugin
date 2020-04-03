package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// ControlDevices contains the definitions of all the "upsControl" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.8
// http://www.oidview.com/mibs/0/UPS-MIB.html
var ControlDevices = []*mibs.SnmpDevice{}

// UPS-MIB upsControl device definitions.
var ()
