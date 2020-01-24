package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsSubsetGroupsTable represents SNMP OID .1.3.6.1.2.1.33.3.2.1
type UpsSubsetGroupsTable struct {
	*core.SnmpTable // base class
}

// NewUpsSubsetGroupsTable constructs the UpsSubsetGroupsTable.
func NewUpsSubsetGroupsTable(snmpServerBase *core.SnmpServerBase) (table *UpsSubsetGroupsTable, err error) {
	var tableName = "UPS-MIB-UPS-Subset-Groups-Table"
	var walkOid = ".1.3.6.1.2.1.33.3.2.1"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
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
