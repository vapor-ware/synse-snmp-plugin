package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsBypassHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.5
type UpsBypassHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsBypassHeadersTable constructs the UpsBypassHeadersTable.
func NewUpsBypassHeadersTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsBypassHeadersTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Bypass-Headers-Table", // Table Name
		".1.3.6.1.2.1.33.1.5",              // WalkOid
		[]string{ // Column Names
			"upsBypassFrequency",
			"upsBypassNumLines",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsBypassHeadersTable{SnmpTable: snmpTable}
	return table, nil
}
