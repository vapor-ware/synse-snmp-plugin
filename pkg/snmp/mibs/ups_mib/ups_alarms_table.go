package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsAlarmsTable represents SNMP OID .1.3.6.1.2.1.33.1.6.2
// There are no rows in this table when no alarms are present.
// We have no real data for this row at this time (5/16/2018)
type UpsAlarmsTable struct {
	*core.SnmpTable // base class
}

// NewUpsAlarmsTable constructs the UpsAlarmsTable.
func NewUpsAlarmsTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsAlarmsTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Alarms-Table", // Table Name
		".1.3.6.1.2.1.33.1.6.2",    // WalkOid
		[]string{ // Column Names
			"upsAlarmId",
			"upsAlarmDescr",
			"upsAlarmTime",
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsAlarmsTable{SnmpTable: snmpTable}
	return table, nil
}
