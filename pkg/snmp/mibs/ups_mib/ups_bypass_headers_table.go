package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsBypassHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.5
type UpsBypassHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsBypassHeadersTable constructs the UpsBypassHeadersTable.
func NewUpsBypassHeadersTable(snmpServerBase *core.SnmpServerBase) (table *UpsBypassHeadersTable, err error) {
	var tableName = "UPS-MIB-UPS-Bypass-Headers-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.5"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsBypassFrequency",
			"upsBypassNumLines",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsBypassHeadersTable{SnmpTable: snmpTable}
	return table, nil
}
