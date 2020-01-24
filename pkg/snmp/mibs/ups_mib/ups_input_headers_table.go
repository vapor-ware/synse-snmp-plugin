package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsInputHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.3
type UpsInputHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsInputHeadersTable constructs the UpsInputHeadersTable.
func NewUpsInputHeadersTable(snmpServerBase *core.SnmpServerBase) (table *UpsInputHeadersTable, err error) {
	var tableName = "UPS-MIB-UPS-Input-Headers-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.3"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsInputLineBads",
			"upsInputNumLines",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true,           // flattened table
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"table": tableName,
		}).Error("[snmp] failed to create table")
		return nil, err
	}

	table = &UpsInputHeadersTable{SnmpTable: snmpTable}
	return table, nil
}
