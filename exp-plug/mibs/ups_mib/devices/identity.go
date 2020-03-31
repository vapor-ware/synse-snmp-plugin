package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// IdentityDevices contains the definitions of all the "upsIdent" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.1
// http://www.oidview.com/mibs/0/UPS-MIB.html
var IdentityDevices = []*mibs.SnmpDevice{
	&UpsIdentManufacturer,
	&UpsIdentModel,
	&UpsIdentUPSSoftwareVersion,
	&UpsIdentAgentSoftwareVersion,
	&UpsIdentName,
	&UpsIdentAttachedDevices,
}

// UPS-MIB upsIdent device definitions.
var (
	UpsIdentManufacturer = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.1.1",
		Info:    "upsIdentManufacturer",
		Type:    "identity",
		Output:  "identity",
		Handler: "read-only",
	}

	UpsIdentModel = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.1.2",
		Info:    "upsIdentModel",
		Type:    "identity",
		Output:  "identity",
		Handler: "read-only",
	}

	UpsIdentUPSSoftwareVersion = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.1.3",
		Info:    "upsIdentUPSSoftwareVersion",
		Type:    "identity",
		Output:  "identity",
		Handler: "read-only",
	}

	UpsIdentAgentSoftwareVersion = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.1.4",
		Info:    "upsIdentAgentSoftwareVersion",
		Type:    "identity",
		Output:  "identity",
		Handler: "read-only",
	}

	UpsIdentName = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.1.5",
		Info:    "upsIdentName",
		Type:    "identity",
		Output:  "identity",
		Handler: "read-only",
	}

	UpsIdentAttachedDevices = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.1.6",
		Info:    "upsIdentAttachedDevices",
		Type:    "identity",
		Output:  "identity",
		Handler: "read-only",
	}
)
