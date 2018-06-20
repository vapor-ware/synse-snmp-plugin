package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsBatteryTable represents SNMP OID .1.3.6.1.2.1.33.1.2
type UpsBatteryTable struct {
	*core.SnmpTable // base class
}

// NewUpsBatteryTable constructs the UpsBatteryTable.
func NewUpsBatteryTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsBatteryTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Battery-Table", // Table Name
		".1.3.6.1.2.1.33.1.2",       // WalkOid
		[]string{ // Column Names
			"upsBatteryStatus",
			"upsSecondsOnBattery", // Zero if not on battery power.
			"upsEstimatedMinutesRemaining",
			"upsEstimatedChargeRemaining", // Percentage
			"upsBatteryVoltage",           // Units .1 VDC.
			"upsBatteryCurrent",           // Units .1 Amp DC.
			"upsBacontteryTemperature",    // Units degrees C.
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsBatteryTable{SnmpTable: snmpTable}
	// Override the default Device Enumerator
	table.DevEnumerator = UpsBatteryTableDeviceEnumerator{table}
	return table, nil
}

// UpsBatteryTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the battery table.
type UpsBatteryTableDeviceEnumerator struct {
	Table *UpsBatteryTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsBatteryTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*sdk.DeviceConfig, err error) {

	// Get the rack and board ids. Setup the location.
	rack, board, err := core.GetRackAndBoard(data)
	if err != nil {
		return nil, err
	}

	// Pull out the table, mib, device model, SNMP DeviceConfig
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

	temperatureKind := &sdk.DeviceKind{
		Name: "temperature",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "temperature"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	cfg.Devices = []*sdk.DeviceKind{
		statusKind,
		voltageKind,
		currentKind,
		temperatureKind,
	}

	// This is always a single row table.

	// upsBatteryStatus
	// deviceData gets shimmed into the DeviceConfig for each synse device.
	// It varies slightly for each device below.
	deviceData := map[string]interface{}{
		"info":       "upsBatteryStatus",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "1",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 1), // base_oid and integer column.
		// This is an enumeration. We need to translate the integer we read to a string.
		// Enumeration data. For now we have map[string]string to work with so the
		// key is fmt.Sprintf("enumeration%d", reading).
		"enumeration1": "unknown",
		"enumeration2": "batteryNormal",
		"enumeration3": "batteryLow",
		"enumeration4": "batteryDepleted",
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device := &sdk.DeviceInstance{
		Location: snmpLocation,
		Data:     deviceData,
	}
	statusKind.Instances = append(statusKind.Instances, device)

	// upsSecondsOnBattery --------------------------------------------------------
	deviceData = map[string]interface{}{
		"info":       "upsSecondsOnBattery",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "2",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 2), // base_oid and integer column.
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

	// upsEstimatedMinutesRemaining -----------------------------------------------
	deviceData = map[string]interface{}{
		"info":       "upsEstimatedMinutesRemaining",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "3",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 3), // base_oid and integer column.
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

	// upsEstimatedChargeRemaining ------------------------------------------------
	deviceData = map[string]interface{}{
		"info":       "upsEstimatedChargeRemaining",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "4",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 4), // base_oid and integer column.
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

	// upsBatteryVoltage ----------------------------------------------------------
	deviceData = map[string]interface{}{
		"info":       "upsBatteryVoltage",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "5",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 5), // base_oid and integer column.
		"multiplier": ".1",                                  // Units are 0.1 Volt DC.
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device = &sdk.DeviceInstance{
		Location: snmpLocation,
		Data:     deviceData,
	}
	voltageKind.Instances = append(voltageKind.Instances, device)

	// upsBatteryCurrent ---------------------------------------------------------
	deviceData = map[string]interface{}{
		"info":       "upsBatteryCurrent",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "6",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 6), // base_oid and integer column.
		"multiplier": ".1",                                  // Units are 0.1 Amp DC.
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

	// upsBatteryTemperature  -----------------------------------------------------
	deviceData = map[string]interface{}{
		"info":       "upsBatteryTemperature",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "7",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 7), // base_oid and integer column.
		// No multiplier needed. Units are degrees C.
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device = &sdk.DeviceInstance{
		Location: snmpLocation,
		Data:     deviceData,
	}
	temperatureKind.Instances = append(temperatureKind.Instances, device)

	devices = append(devices, cfg)
	return devices, err
}
