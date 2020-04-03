package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// OutputDevices contains the definitions of all the "upsOutput" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.4
// http://www.oidview.com/mibs/0/UPS-MIB.html
var OutputDevices = []*mibs.SnmpDevice{}

// UPS-MIB upsOutput device definitions.
var ()
