package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsControlTable represents SNMP OID 1.3.6.1.2.1.33.1.8
type UpsControlTable struct {
	*core.SnmpTable // base class
}

// NewUpsControlTable constructs the UpsControlTable.
func NewUpsControlTable(snmpServerBase *core.SnmpServerBase) (table *UpsControlTable, err error) {
	var tableName = "UPS-MIB-UPS-Control-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.8"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsShutdownType",
			"upsShutdownTypeAfterDelay", // Seconds
			"upsStartupAfterDelay",      // Seconds
			"upsRebootWithDuration",     // Seconds
			"upsAutoRestart",
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

	table = &UpsControlTable{SnmpTable: snmpTable}
	return table, nil
}
