package mibs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsInputTable represents SNMP OID .1.3.6.1.2.1.33.1.3.3
type UpsInputTable struct {
	*core.SnmpTable // base class
}

// NewUpsInputTable constructs the UpsInputTable.
func NewUpsInputTable(snmpServerBase *core.SnmpServerBase) (table *UpsInputTable, err error) {
	var tableName = "UPS-MIB-UPS-Input-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.3.3"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsInputLineIndex", // MIB says not accessable. Have seen it in walks.
			"upsInputFrequency", // .1 Hertz
			"upsInputVoltage",   // RMS Volts
			"upsInputCurrent",   // .1 RMS Amp
			"upsInputTruePower", // Down with False Power! (Watts)
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsInputTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsInputTableDeviceEnumerator{table}
	return table, nil
}

// UpsInputTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the input table.
type UpsInputTableDeviceEnumerator struct {
	Table *UpsInputTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsInputTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceProto, err error) {

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	// We will have "frequency", "voltage", "current", and "power" device kinds.
	// There is probably a better way of doing this, but this just gets things to
	// where they need to be for now.
	frequencyProto := &config.DeviceProto{
		Type: "frequency",
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
		frequencyProto,
		voltageProto,
		currentProto,
		powerProto,
	}

	for i := 0; i < len(table.Rows); i++ {
		// upsInputFrequency
		// deviceData gets shimmed into the DeviceConfig for each synse device.
		// It varies slightly for each device below.
		deviceData := map[string]interface{}{
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "2",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 2), // base_oid and integer column.
			"multiplier": float32(0.1),                          // Units are 0.1 Hertz
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device := &config.DeviceInstance{
			Info: fmt.Sprintf("upsInputFrequency%d", i),
			Data: deviceData,
		}
		frequencyProto.Instances = append(frequencyProto.Instances, device)

		// upsInputVoltage ----------------------------------------------------------
		deviceData = map[string]interface{}{
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "3",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 3), // base_oid and integer column.
			// No multiplier needed. Units are RMS Volts.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &config.DeviceInstance{
			Info: fmt.Sprintf("upsInputVoltage%d", i),
			Data: deviceData,
		}
		voltageProto.Instances = append(voltageProto.Instances, device)

		// upsInputCurrent ----------------------------------------------------------
		deviceData = map[string]interface{}{
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "4",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 4), // base_oid and integer column.
			"multiplier": float32(0.1),                          // Units are 0.1 RMS Amp
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &config.DeviceInstance{
			Info: fmt.Sprintf("upsInputCurrent%d", i),
			Data: deviceData,
		}
		currentProto.Instances = append(currentProto.Instances, device)

		// upsInputTruePower --------------------------------------------------------
		deviceData = map[string]interface{}{
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "5",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 5), // base_oid and integer column.
			// Output is in Watts. No multiplier needed.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &config.DeviceInstance{
			Info: fmt.Sprintf("upsInputTruePower%d", i),
			Data: deviceData,
		}
		powerProto.Instances = append(powerProto.Instances, device)
	}

	return devices, err
}
