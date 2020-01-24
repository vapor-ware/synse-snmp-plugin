package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsWellKnownTestsTable represents SNMP OID .1.3.6.1.2.1.33.1.7.7
type UpsWellKnownTestsTable struct {
	*core.SnmpTable // base class
}

// NewUpsWellKnownTestsTable constructs the UpsWellKnownTestsTable.
func NewUpsWellKnownTestsTable(snmpServerBase *core.SnmpServerBase) (table *UpsWellKnownTestsTable, err error) {
	var tableName = "UPS-MIB-UPS-Well-Known-Tests-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.7.7"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsTestNoTestsInitiated",
			"upsTestAbortTestInProgress",
			"upsTestGeneralSystemsTest",
			"upsTestQuickBatteryTest",
			"upsTestDeepBatteryCalibration",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"1",            // readableColumn
		false,          // flattened table
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"table": tableName,
		}).Error("[snmp] failed to create table")
		return nil, err
	}

	table = &UpsWellKnownTestsTable{SnmpTable: snmpTable}
	return table, nil
}
