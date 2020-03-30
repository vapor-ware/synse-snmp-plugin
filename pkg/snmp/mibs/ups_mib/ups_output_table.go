package mibs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsOutputTable represents SNMP OID .1.3.6.1.2.1.33.1.4.4
type UpsOutputTable struct {
	*core.SnmpTable // base class
}

// NewUpsOutputTable constructs the UpsOutputTable.
func NewUpsOutputTable(snmpServerBase *core.SnmpServerBase) (table *UpsOutputTable, err error) {
	var tableName = "UPS-MIB-UPS-Output-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.4.4"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsOutputLineIndex", // MIB says not accessible. Have seen it in walks.
			"upsOutputVoltage",   // RMS Volts
			"upsOutputCurrent",   // .1 RMS Amp
			"upsOutputPower",     // Watts
			"upsOutputPercentLoad",
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
	data map[string]interface{}) (devices []*config.DeviceProto, err error) {

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	// We will have "status", "voltage", "current", and "temperature" device kinds.
	// There is probably a better way of doing this, but this just gets things to
	// where they need to be for now.
	statusProto := &config.DeviceProto{
		Type: "status",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
	}

	voltageProto := &config.DeviceProto{
		Type: "voltage",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
	}

	currentProto := &config.DeviceProto{
		Type: "current",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
	}

	powerProto := &config.DeviceProto{
		Type: "power",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
	}

	devices = []*config.DeviceProto{
		statusProto,
		voltageProto,
		currentProto,
		powerProto,
	}

	for i := 0; i < len(table.Rows); i++ {
		// upsOutputVoltage
		// deviceData gets shimmed into the DeviceConfig for each synse device.
		// It varies slightly for each device below.
		deviceData := map[string]interface{}{
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

		device := &config.DeviceInstance{
			Info: fmt.Sprintf("upsOutputVoltage%d", i),
			Data: deviceData,
		}
		voltageProto.Instances = append(voltageProto.Instances, device)

		// upsOutputCurrent ----------------------------------------------------------
		deviceData = map[string]interface{}{
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "3",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 3), // base_oid and integer column.
			"multiplier": float32(0.1),                          // Units are 0.1 RMS Amp
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &config.DeviceInstance{
			Info: fmt.Sprintf("upsOutputCurrent%d", i),
			Data: deviceData,
		}
		currentProto.Instances = append(currentProto.Instances, device)

		// upsOutputPower -------------------------------------------------------------
		deviceData = map[string]interface{}{
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

		device = &config.DeviceInstance{
			Info: fmt.Sprintf("upsOutputPower%d", i),
			Data: deviceData,
		}
		powerProto.Instances = append(powerProto.Instances, device)

		// upsOutputPercentLoad -------------------------------------------------------
		deviceData = map[string]interface{}{
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

		device = &config.DeviceInstance{
			Info: fmt.Sprintf("upsOutputPercentLoad%d", i),
			Data: deviceData,
		}
		statusProto.Instances = append(statusProto.Instances, device)
	}

	return devices, err
}
