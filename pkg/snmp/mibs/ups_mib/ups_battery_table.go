package mibs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsBatteryTable represents SNMP OID .1.3.6.1.2.1.33.1.2
type UpsBatteryTable struct {
	*core.SnmpTable // base class
}

// NewUpsBatteryTable constructs the UpsBatteryTable.
func NewUpsBatteryTable(snmpServerBase *core.SnmpServerBase) (table *UpsBatteryTable, err error) {
	var tableName = "UPS-MIB-UPS-Battery-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.2"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsBatteryStatus",
			"upsSecondsOnBattery", // Zero if not on battery power.
			"upsEstimatedMinutesRemaining",
			"upsEstimatedChargeRemaining", // Percentage
			"upsBatteryVoltage",           // Units .1 VDC.
			"upsBatteryCurrent",           // Units .1 Amp DC.
			"upsBatteryTemperature",       // Units degrees C.
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
	data map[string]interface{}) (devices []*config.DeviceProto, err error) {

	// Pull out the table, mib, device model, SNMP DeviceConfig
	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	// We will have "status", "voltage", "current", "temperature", "percentage", "minutes", and "seconds" device kinds.
	// There is probably a better way of doing this, but this just gets things to
	// where they need to be for now.
	// The code below is hooking up the DeviceHandler via the Type string to the Output Type string.
	statusProto := &config.DeviceProto{
		Type: "status",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
		Tags:      snmpDeviceConfigMap["deviceTags"].([]string),
	}

	voltageProto := &config.DeviceProto{
		Type: "voltage",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
		Tags:      snmpDeviceConfigMap["deviceTags"].([]string),
	}

	currentProto := &config.DeviceProto{
		Type: "current",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
		Tags:      snmpDeviceConfigMap["deviceTags"].([]string),
	}

	temperatureProto := &config.DeviceProto{
		Type: "temperature",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
		Tags:      snmpDeviceConfigMap["deviceTags"].([]string),
	}

	percentageProto := &config.DeviceProto{
		Type: "percentage",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
		Tags:      snmpDeviceConfigMap["deviceTags"].([]string),
	}

	minutesProto := &config.DeviceProto{
		Type: "minutes",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
		Tags:      snmpDeviceConfigMap["deviceTags"].([]string),
	}

	secondsProto := &config.DeviceProto{
		Type: "seconds",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
		Tags:      snmpDeviceConfigMap["deviceTags"].([]string),
	}

	devices = []*config.DeviceProto{
		statusProto,
		voltageProto,
		currentProto,
		temperatureProto,
		percentageProto,
		minutesProto,
		secondsProto,
	}

	// This is always a single row table.

	// upsBatteryStatus ---------------------------------------------------
	// deviceData gets shimmed into the DeviceConfig for each synse device.
	// It varies slightly for each device below.
	deviceData := map[string]interface{}{
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "1",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 1), // base_oid and integer column.
		// This is an enumeration. We need to translate the integer we read to a string.
		"enumeration": "true", // Defines that this is an enumeration.
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

	device := &config.DeviceInstance{
		Info: "upsBatteryStatus",
		Data: deviceData,
	}
	statusProto.Instances = append(statusProto.Instances, device)

	// upsSecondsOnBattery --------------------------------------------------------
	deviceData = map[string]interface{}{
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

	device = &config.DeviceInstance{
		Info: "upsSecondsOnBattery",
		Data: deviceData,
	}
	secondsProto.Instances = append(secondsProto.Instances, device)

	// upsEstimatedMinutesRemaining -----------------------------------------------
	deviceData = map[string]interface{}{
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

	device = &config.DeviceInstance{
		Info: "upsEstimatedMinutesRemaining",
		Data: deviceData,
	}
	minutesProto.Instances = append(minutesProto.Instances, device)

	// upsEstimatedChargeRemaining ------------------------------------------------
	deviceData = map[string]interface{}{
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

	device = &config.DeviceInstance{
		Info: "upsEstimatedChargeRemaining",
		Data: deviceData,
	}
	percentageProto.Instances = append(percentageProto.Instances, device)

	// upsBatteryVoltage ----------------------------------------------------------
	deviceData = map[string]interface{}{
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "5",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 5), // base_oid and integer column.
		"multiplier": float32(0.1),                          // Units are 0.1 Volt DC.
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device = &config.DeviceInstance{
		Info: "upsBatteryVoltage",
		Data: deviceData,
	}
	voltageProto.Instances = append(voltageProto.Instances, device)

	// upsBatteryCurrent ---------------------------------------------------------
	deviceData = map[string]interface{}{
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "6",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 6), // base_oid and integer column.
		"multiplier": float32(0.1),                          // Units are 0.1 Amp DC.
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device = &config.DeviceInstance{
		Info: "upsBatteryCurrent",
		Data: deviceData,
	}
	currentProto.Instances = append(currentProto.Instances, device)

	// upsBatteryTemperature  -----------------------------------------------------
	deviceData = map[string]interface{}{
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

	device = &config.DeviceInstance{
		Info: "upsBatteryTemperature",
		Data: deviceData,
	}
	temperatureProto.Instances = append(temperatureProto.Instances, device)

	return devices, err
}
