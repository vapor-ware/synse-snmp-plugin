package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsWellKnownAlarmsTable represents SNMP OID .1.3.6.1.2.1.33.1.6.3
type UpsWellKnownAlarmsTable struct {
	*core.SnmpTable // base class
}

// NewUpsWellKnownAlarmsTable constructs the UpsWellKnownAlarmsTable.
func NewUpsWellKnownAlarmsTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsWellKnownAlarmsTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Well-Known-Alarms-Table", // Table Name
		".1.3.6.1.2.1.33.1.6.3",               // WalkOid
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
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsWellKnownAlarmsTable{SnmpTable: snmpTable}
	return table, nil
}
