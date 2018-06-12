package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsControlTable represts SNMP OID 1.3.6.1.2.1.33.1.8
type UpsControlTable struct {
	*core.SnmpTable // base class
}

// NewUpsControlTable constructs the UpsControlTable.
func NewUpsControlTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsControlTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Control-Table", // Table Name
		".1.3.6.1.2.1.33.1.8",       // WalkOid
		[]string{ // Column Names
			"upsShutdownType",
			"upsShutdownTypeAfterDelay", // Seconds
			"upsStartupAfterDelay",      // Seconds
			"upsRebootWithDuration",     // Seconds
			"upsAutoRestart",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsControlTable{SnmpTable: snmpTable}
	return table, nil
}
