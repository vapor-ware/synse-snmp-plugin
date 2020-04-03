package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// ConfigDevices contains the definitions of all the "upsConfig" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.9
// http://www.oidview.com/mibs/0/UPS-MIB.html
var ConfigDevices = []*mibs.SnmpDevice{}

// UPS-MIB upsConfig device definitions.
var ()
