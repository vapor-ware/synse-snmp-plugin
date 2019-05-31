package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsBasicGroupsTable represents SNMP OID .1.3.6.1.2.1.33.1.6.1.0
// We currently have no data for this table (5/16/2018)
type UpsBasicGroupsTable struct {
	*core.SnmpTable // base class
}

// NewUpsBasicGroupsTable constructs the UpsBasicGroupsTable.
func NewUpsBasicGroupsTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsBasicGroupsTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Basic-Groups-Table", // Table Name
		".1.3.6.1.2.1.33.3.2.2",          // WalkOid
		[]string{ // Column Names
			"upsBasicIdentGroup",
			"upsBasicBatteryGroup",
			"upsBasicInpuGroup",
			"upsBasicOutputGroup",
			"upsBasicBypassGroup",
			"upsBasicAlarmGroup",
			"upsBasicTestGroup",
			"upsBasicControlGroup",
			"upsBasicConfigGroup",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsBasicGroupsTable{SnmpTable: snmpTable}
	return table, nil
}
