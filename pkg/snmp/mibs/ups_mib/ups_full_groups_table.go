package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsFullGroupsTable represents SNMP OID .1.3.6.1.2.1.33.3.2.3
type UpsFullGroupsTable struct {
	*core.SnmpTable // base class
}

// NewUpsFullGroupsTable constructs the UpsFullGroupsTable.
func NewUpsFullGroupsTable(snmpServerBase *core.SnmpServerBase) (table *UpsFullGroupsTable, err error) {
	var tableName = "UPS-MIB-UPS-Full-Groups-Table"
	var walkOid = ".1.3.6.1.2.1.33.3.2.3"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
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
