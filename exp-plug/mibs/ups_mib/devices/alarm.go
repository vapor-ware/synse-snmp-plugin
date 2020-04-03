package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// AlarmDevices contains the definitions of all the "upsAlarm" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.6
// http://www.oidview.com/mibs/0/UPS-MIB.html
var AlarmDevices = []*mibs.SnmpDevice{}

// UPS-MIB upsAlarm device definitions.
var ()
