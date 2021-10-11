package mibs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// upsAlarmsInfo is what the device info will be under SNMP OID
// .1.3.6.1.2.1.33.1.6.3 view it in a MIB browser to see the names.
var upsAlarmsInfo = []string{
	"upsAlarmBatteryBad",
	"upsAlarmOnBattery",
	"upsAlarmLowBattery",
	"upsAlarmDepletedBattery",
	"upsAlarmTempBad",
	"upsAlarmInputBad",
	"upsAlarmOutputBad",
	"upsAlarmOutputOverload",
	"upsAlarmOnBypass",
	"upsAlarmBypassBad",
	"upsAlarmOutputOffAsRequested",
	"upsAlarmUpsOffAsRequested",
	"upsAlarmChargerFailed",
	"upsAlarmUpsOutputOff",
	"upsAlarmUpsSystemOff",
	"upsAlarmFanFailure",
	"upsAlarmFuseFailure",
	"upsAlarmGeneralFault",
	"upsAlarmDiagnosticTestFailed",
	"upsAlarmCommunicationsLost",
	"upsAlarmAwaitingPower",
	"upsAlarmShutdownPending",
	"upsAlarmShutdownImminent",
	"upsAlarmTestInProgress",
}

// UpsAlarmsTable represents SNMP OID .1.3.6.1.2.1.33.1.6.2
// There are no rows in this table when no alarms are present.
// We have no real data for this row at this time (5/16/2018)
type UpsAlarmsTable struct {
	*core.SnmpTable // base class
}

// NewUpsAlarmsTable constructs the UpsAlarmsTable.
func NewUpsAlarmsTable(snmpServerBase *core.SnmpServerBase) (table *UpsAlarmsTable, err error) {
	var tableName = "UPS-MIB-UPS-Alarms-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.6.2"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsAlarmId",
			"upsAlarmDescr",
			"upsAlarmTime",
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false,          // flattened table
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"table": tableName,
		}).Error("[snmp] failed to create table")
		return nil, err
	}

	table = &UpsAlarmsTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsAlarmsTableDeviceEnumerator{table}
	return table, nil
}

// UpsAlarmsTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the alarms headers table.
type UpsAlarmsTableDeviceEnumerator struct {
	Table *UpsAlarmsTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsAlarmsTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceProto, err error) {

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return
	}

	statusProto := &config.DeviceProto{
		Type: "status",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
		Tags:      snmpDeviceConfigMap["deviceTags"].([]string),
	}

	devices = []*config.DeviceProto{
		statusProto,
	}

	for i := 0; i < len(table.Rows); i++ {

		// upsAlarmDescr ---------------------------------------------------------
		deviceData := map[string]interface{}{
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "2",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 2), // base_oid and integer column.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device := &config.DeviceInstance{
			Info: upsAlarmsInfo[i],
			Data: deviceData,
		}
		statusProto.Instances = append(statusProto.Instances, device)

		// upsAlarmTime ---------------------------------------------------------
		deviceData = map[string]interface{}{
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "3",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 3), // base_oid and integer column.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &config.DeviceInstance{
			Info: fmt.Sprintf("%sTime", upsAlarmsInfo[i]),
			Data: deviceData,
		}
		statusProto.Instances = append(statusProto.Instances, device)

	} // end for

	return
}
