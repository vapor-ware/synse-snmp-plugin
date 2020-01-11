package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsAlarmsHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.6
// This is just the alarms present column. OID: .1.3.6.1.2.1.33.1.6.1.0
type UpsAlarmsHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsAlarmsHeadersTable constructs the UpsAlarmsHeadersTable.
func NewUpsAlarmsHeadersTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsAlarmsHeadersTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Alarms-Headers-Table", // Table Name
		".1.3.6.1.2.1.33.1.6",              // WalkOid
		[]string{ // Column Names
			"upsAlarmsPresent", // The present number of active alarm conditions.
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsAlarmsHeadersTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsAlarmsHeadersTableDeviceEnumerator{table}
	return table, nil
}

// UpsAlarmsHeadersTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the alarms headers table.
type UpsAlarmsHeadersTableDeviceEnumerator struct {
	Table *UpsAlarmsHeadersTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsAlarmsHeadersTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*sdk.DeviceConfig, err error) {

	// Get the rack and board ids. Setup the location.
	rack, board, err := core.GetRackAndBoard(data)
	if err != nil {
		return
	}

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return
	}

	cfg := &sdk.DeviceConfig{
		SchemeVersion: sdk.SchemeVersion{Version: "1.0"},
		Locations: []*sdk.LocationConfig{
			{
				Name:  snmpLocation,
				Rack:  &sdk.LocationData{Name: rack},
				Board: &sdk.LocationData{Name: board},
			},
		},
		Devices: []*sdk.DeviceKind{},
	}

	// We have only "status-int" device kinds.
	statusIntKind := &sdk.DeviceKind{
		Name: "status-int",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "status-int"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	// This gets the devices in the enumerated output, meaning they show up in a scan.
	cfg.Devices = []*sdk.DeviceKind{
		statusIntKind,
	}

	// This is always a single row table.

	// upsAlarmsPresent ---------------------------------------------------------
	deviceData := map[string]interface{}{
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "1",
		"column":     "0",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 0), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device := &sdk.DeviceInstance{
		Info:     "upsAlarmsPresent",
		Location: snmpLocation,
		Data:     deviceData,
	}
	statusIntKind.Instances = append(statusIntKind.Instances, device)

	devices = append(devices, cfg)
	return
}
