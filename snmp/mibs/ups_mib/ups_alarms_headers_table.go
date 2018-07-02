package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsAlarmsHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.6
// This is just the alarms present column. OID: .1.3.6.1.2.1.33.1.6.1.0
type UpsAlarmsHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsAlarmsHeadersTable constructs the UpsAlarmsHeadersTable.
func NewUpsAlarmsHeadersTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsAlarmsHeadersTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Alarms-Headers-Table", // Table Name
		".1.3.6.1.2.1.33.1.6",              // WalkOid
		[]string{ // Column Names
			"upsAlarmsPresent", // The present number of active alarm conditions.
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsAlarmsHeadersTable{SnmpTable: snmpTable}
	return table, nil
}
