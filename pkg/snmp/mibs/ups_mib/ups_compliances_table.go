package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsCompliancesTable represents SNMP OID .1.3.6.1.2.1.33.3.1
type UpsCompliancesTable struct {
	*core.SnmpTable // base class
}

// NewUpsCompliancesTable constructs the UpsCompliancesTable.
func NewUpsCompliancesTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsCompliancesTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Compliances-Table", // Table Name
		".1.3.6.1.2.1.33.3.1",           // WalkOid
		[]string{ // Column Names
			"upsSubsetCompliance",
			"upsBasicCompliance",
			"upsFullCompliance",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsCompliancesTable{SnmpTable: snmpTable}
	return table, nil
}
