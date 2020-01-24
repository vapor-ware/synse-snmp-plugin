package mibs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsBypassTable represents SNMP OID .1.3.6.1.2.1.33.1.5.3
type UpsBypassTable struct {
	*core.SnmpTable // base class
}

// NewUpsBypassTable constructs the UpsBypassTable.
func NewUpsBypassTable(snmpServerBase *core.SnmpServerBase) (table *UpsBypassTable, err error) {
	var tableName = "UPS-MIB-UPS-Bypass-Table"
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
			"upsBypassLineIndex",
			"upsBypassVoltage",
			"upsBypassCurrent",
			"upsBypassPower",
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

	table = &UpsBypassTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsBypassTableDeviceEnumerator{table}
	return table, nil
}

// UpsBypassTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the bypass table.
type UpsBypassTableDeviceEnumerator struct {
	Table *UpsBypassTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsBypassTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceProto, err error) {

	// Pull out the table, mib, device model, SNMP DeviceConfig.
	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	// We will have "voltage", "current", and "power" device kinds.
	// There is probably a better way of doing this, but this just gets things to
	// where they need to be for now.
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
		voltageProto,
		currentProto,
		powerProto,
	}

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(table.Rows); i++ {
		// upsBypassVoltage ---------------------------------------------------
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
			Info: fmt.Sprintf("upsBypassVoltage%d", i),
			Data: deviceData,
		}
		voltageProto.Instances = append(voltageProto.Instances, device)

		// upsBypassCurrent ---------------------------------------------------------
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
			Info: fmt.Sprintf("upsBypassCurrent%d", i),
			Data: deviceData,
		}
		currentProto.Instances = append(currentProto.Instances, device)

		// upsBypassPower --------------------------------------------------------------
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
			Info: fmt.Sprintf("upsBypassPower%d", i),
			Data: deviceData,
		}
		powerProto.Instances = append(powerProto.Instances, device)
	}

	return devices, err
}
