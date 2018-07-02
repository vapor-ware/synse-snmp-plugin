package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsOutputTable represents SNMP OID .1.3.6.1.2.1.33.1.4.4
type UpsOutputTable struct {
	*core.SnmpTable // base class
}

// NewUpsOutputTable constructs the UpsOutputTable.
func NewUpsOutputTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsOutputTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Output-Table", // Table Name
		".1.3.6.1.2.1.33.1.4.4",    // WalkOid
		[]string{ // Column Names
			"upsOutputLineIndex", // MIB says not accessable. Have seen it in walks.
			"upsOutputVoltage",   // RMS Volts
			"upsOutputCurrent",   // .1 RMS Amp
			"upsOutputPower",     // Watts
			"upsOutputPercentLoad",
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsOutputTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsOutputTableDeviceEnumerator{table}
	return table, nil
}

// UpsOutputTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the output table.
type UpsOutputTableDeviceEnumerator struct {
	Table *UpsOutputTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsOutputTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*sdk.DeviceConfig, err error) {

	// Get the rack and board ids. Setup the location.
	rack, board, err := core.GetRackAndBoard(data)
	if err != nil {
		return nil, err
	}

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
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

	// We will have "status", "voltage", "current", and "temperature" device kinds.
	// There is probably a better way of doing this, but this just gets things to
	// where they need to be for now.
	statusKind := &sdk.DeviceKind{
		Name: "status",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "status"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	voltageKind := &sdk.DeviceKind{
		Name: "voltage",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "voltage"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	currentKind := &sdk.DeviceKind{
		Name: "current",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "current"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	powerKind := &sdk.DeviceKind{
		Name: "power",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "watts.power"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	cfg.Devices = []*sdk.DeviceKind{
		statusKind,
		voltageKind,
		currentKind,
		powerKind,
	}

	for i := 0; i < len(table.Rows); i++ {
		// upsOutputVoltage
		// deviceData gets shimmed into the DeviceConfig for each synse device.
		// It varies slightly for each device below.
		deviceData := map[string]interface{}{
			"info":       fmt.Sprintf("upsOutputVoltage%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "2",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 2), // base_oid and integer column.
			// No multiplier needed. Units are RMS Volts.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device := &sdk.DeviceInstance{
			Location: snmpLocation,
			Data:     deviceData,
		}
		voltageKind.Instances = append(voltageKind.Instances, device)

		// upsOutputCurrent ----------------------------------------------------------
		deviceData = map[string]interface{}{
			"info":       fmt.Sprintf("upsOutputCurrent%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "3",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 3), // base_oid and integer column.
			"multiplier": float32(0.1),                                  // Units are 0.1 RMS Amp
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &sdk.DeviceInstance{
			Location: snmpLocation,
			Data:     deviceData,
		}
		currentKind.Instances = append(currentKind.Instances, device)

		// upsOutputPower -------------------------------------------------------------
		deviceData = map[string]interface{}{
			"info":       fmt.Sprintf("upsOutputPower%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "4",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 4), // base_oid and integer column.
			// Output is in Watts. No multiplier needed.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &sdk.DeviceInstance{
			Location: snmpLocation,
			Data:     deviceData,
		}
		powerKind.Instances = append(powerKind.Instances, device)

		// upsOutputPercentLoad -------------------------------------------------------
		deviceData = map[string]interface{}{
			"info":       fmt.Sprintf("upsOutputPercentLoad%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "5",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 5), // base_oid and integer column.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &sdk.DeviceInstance{
			Location: snmpLocation,
			Data:     deviceData,
		}
		statusKind.Instances = append(statusKind.Instances, device)
	}

	devices = append(devices, cfg)
	return devices, err
}
