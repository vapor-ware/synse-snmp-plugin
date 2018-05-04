package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsSubsetGroupsTable represts SNMP OID .1.3.6.1.2.1.33.3.2.1
type UpsSubsetGroupsTable struct {
	*core.SnmpTable // base class
}

// NewUpsSubsetGroupsTable constructs the UpsSubsetGroupsTable.
func NewUpsSubsetGroupsTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsSubsetGroupsTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Subset-Groups-Table", // Table Name
		".1.3.6.1.2.1.33.3.2.1",           // WalkOid
		[]string{ // Column Names
			"upsSubsetIdentGroup",
			"upsSubsetBatteryGroup",
			"upsSubsetInputGroup",
			"upsSubsetOutputGroup",
			"upsSubsetAlarmGroup",
			"upsSubsetControlGroup",
			"upsSubsetConfigGroup",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsSubsetGroupsTable{SnmpTable: snmpTable}
	return table, nil
}
