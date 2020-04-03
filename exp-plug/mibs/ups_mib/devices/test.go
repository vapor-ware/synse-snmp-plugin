package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// TestDevices contains the definitions of all the "upsTest" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.7
// http://www.oidview.com/mibs/0/UPS-MIB.html
var TestDevices = []*mibs.SnmpDevice{}

// UPS-MIB upsTest device definitions.
var ()
