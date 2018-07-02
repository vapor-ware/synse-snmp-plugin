package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsInputHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.3
type UpsInputHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsInputHeadersTable constructs the UpsInputHeadersTable.
func NewUpsInputHeadersTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsInputHeadersTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Input-Headers-Table", // Table Name
		".1.3.6.1.2.1.33.1.3",             // WalkOid
		[]string{ // Column Names
			"upsInputLineBads",
			"upsInputNumLines",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsInputHeadersTable{SnmpTable: snmpTable}
	return table, nil
}
