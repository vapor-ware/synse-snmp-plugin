package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsFullGroupsTable represents SNMP OID .1.3.6.1.2.1.33.3.2.3
type UpsFullGroupsTable struct {
	*core.SnmpTable // base class
}

// NewUpsFullGroupsTable constructs the UpsFullGroupsTable.
func NewUpsFullGroupsTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsFullGroupsTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Full-Groups-Table", // Table Name
		".1.3.6.1.2.1.33.3.2.3",         // WalkOid
		[]string{ // Column Names
			"upsFullIdentGroup",
			"upsFullBatteryGroup",
			"upsFullInputGroup",
			"upsFullOutputGroup",
			"upsFullBypassGroup",
			"upsFullAlarmGroup",
			"upsFullTestGroup",
			"upsFullControlGroup",
			"upsFullConfigGroup",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsFullGroupsTable{SnmpTable: snmpTable}
	return table, nil
}
