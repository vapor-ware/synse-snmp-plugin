package mibs

import (
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsWellKnownAlarmsTable represents SNMP OID .1.3.6.1.2.1.33.1.6.3
type UpsWellKnownAlarmsTable struct {
	*core.SnmpTable // base class
}

// NewUpsWellKnownAlarmsTable constructs the UpsWellKnownAlarmsTable.
func NewUpsWellKnownAlarmsTable(snmpServerBase *core.SnmpServerBase) (table *UpsWellKnownAlarmsTable, err error) {
	var tableName = "UPS-MIB-UPS-Well-Known-Alarms-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.6.3"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsAlarmBatteryBad",
			"upsAlarmOnBattery",
			"upsAlarmLowBattery",
			"upsAlarmDepletedBattery",
			"upsAlarmTempBad",
			"upsAlarmInputBad",
			"upsAlarmOutputBad",
			"upsAlarmOutputOverload",
			"upsAlarmOnBypass",
			"upsAlarmOnBypassBad",
			"upsAlarmOutputOffAsRequested",
			"upsAlarmUpsOffAsRequested",
			"upsAlarmChargerFailed",
			"upsAlarmUpsOutputOff",
			"upsAlarmUpsSystemOff",
			"upsAlarmFanFailure",
			"upsAlarmFuseFailure",
			"upsAlarmGeneralFault",
			"upsAlarmUpsDiagnosticTestFailed",
			"upsAlarmCommunicationsLost",
			"upsAlarmAwaitingPower",
			"upsAlarmShutdownPending",
			"upsAlarmShutdownImminent",
			"upsAlarmTestInProgress",
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

	table = &UpsWellKnownAlarmsTable{SnmpTable: snmpTable}
	return table, nil
}
