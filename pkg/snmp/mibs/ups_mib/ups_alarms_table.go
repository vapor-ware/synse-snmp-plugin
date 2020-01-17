package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsAlarmsTable represents SNMP OID .1.3.6.1.2.1.33.1.6.2
// There are no rows in this table when no alarms are present.
// We have no real data for this row at this time (5/16/2018)
type UpsAlarmsTable struct {
	*core.SnmpTable // base class
}

// NewUpsAlarmsTable constructs the UpsAlarmsTable.
func NewUpsAlarmsTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsAlarmsTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Alarms-Table", // Table Name
		".1.3.6.1.2.1.33.1.6.2",    // WalkOid
		[]string{ // Column Names
			"upsAlarmId",
			"upsAlarmDescr",
			"upsAlarmTime",
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false)          // flattened table
	if err != nil {
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
	}

	devices = []*config.DeviceProto{
		statusProto,
	}

	//// We have "status-uint" and "status-string" device kinds.
	//statusUintKind := &sdk.DeviceKind{
	//	Name: "status-uint",
	//	Metadata: map[string]string{
	//		"model": model,
	//	},
	//	Outputs: []*sdk.DeviceOutput{
	//		{Type: "status-uint"},
	//	},
	//	Instances: []*sdk.DeviceInstance{},
	//}
	//
	//statusStringKind := &sdk.DeviceKind{
	//	Name: "status-string",
	//	Metadata: map[string]string{
	//		"model": model,
	//	},
	//	Outputs: []*sdk.DeviceOutput{
	//		{Type: "status-string"},
	//	},
	//	Instances: []*sdk.DeviceInstance{},
	//}
	//
	//// This gets the devices in the enumerated output, meaning they show up in a scan.
	//cfg.Devices = []*sdk.DeviceKind{
	//	statusUintKind,
	//	statusStringKind,
	//}

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
			Info: fmt.Sprintf("upsAlarm%d", i),
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
			Info: fmt.Sprintf("upsAlarmTime%d", i),
			Data: deviceData,
		}
		statusProto.Instances = append(statusProto.Instances, device)

	} // end for

	return
}
