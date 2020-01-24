package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsTestHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.7
type UpsTestHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsTestHeadersTable constructs the UpsTestHeadersTable.
func NewUpsTestHeadersTable(snmpServerBase *core.SnmpServerBase) (table *UpsTestHeadersTable, err error) {
	var tableName = "UPS-MIB-UPS-Test-Headers-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.7"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsTestId",
			"upsTestSpinLock",
			"upsTestResultsSummary",
			"upsTestResultsDetail",
			"upsTestStartTime",
			"upsTestElapsedTime",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsTestHeadersTable{SnmpTable: snmpTable}
	return table, nil
}
