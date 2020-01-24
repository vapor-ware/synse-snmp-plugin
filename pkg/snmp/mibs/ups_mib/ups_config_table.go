package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsConfigTable represents SNMP OID 1.3.6.1.2.1.33.1.9
type UpsConfigTable struct {
	*core.SnmpTable // base class
}

// NewUpsConfigTable constructs the UpsConfigTable.
func NewUpsConfigTable(snmpServerBase *core.SnmpServerBase) (table *UpsConfigTable, err error) {
	var tableName = "UPS-MIB-UPS-Config-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.9"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsConfigInputVoltage",             // RMS Volts
			"upsConfigInputFreq",                // 0.1 Hertz
			"upsConfigOutputVoltage",            // RMS Volts
			"upsConfigOutputFreq",               // 0.1 Hertz
			"upsConfigOutputVA",                 // Volt Amps
			"upsConfigOutputPower",              // Watts
			"upsConfigLowBattTime",              // Minutes
			"upsConfigAudibleStatus",            // Disabled, Enabled, Muted
			"upsConfigLowVoltageTransferPoint",  // RMS Volts
			"upsConfigHighVoltageTransferPoint", // RMS Volts
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

	table = &UpsConfigTable{SnmpTable: snmpTable}
	return table, nil
}
