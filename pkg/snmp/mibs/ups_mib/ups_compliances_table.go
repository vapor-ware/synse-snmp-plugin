package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsCompliancesTable represents SNMP OID .1.3.6.1.2.1.33.3.1
type UpsCompliancesTable struct {
	*core.SnmpTable // base class
}

// NewUpsCompliancesTable constructs the UpsCompliancesTable.
func NewUpsCompliancesTable(snmpServerBase *core.SnmpServerBase) (table *UpsCompliancesTable, err error) {
	var tableName = "UPS-MIB-UPS-Compliances-Table"
	var walkOid = ".1.3.6.1.2.1.33.3.1"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsSubsetCompliance",
			"upsBasicCompliance",
			"upsFullCompliance",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		false,          // flattened table
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"table": tableName,
		}).Error("[snmp] failed to create table")
		return nil, err
	}

	table = &UpsCompliancesTable{SnmpTable: snmpTable}
	return table, nil
}
