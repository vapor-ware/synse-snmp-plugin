package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsInputTable represts SNMP OID .1.3.6.1.2.1.33.1.3.3
type UpsInputTable struct {
	*core.SnmpTable // base class
}

// NewUpsInputTable constructs the UpsInputTable.
func NewUpsInputTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsInputTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Input-Table", // Table Name
		".1.3.6.1.2.1.33.1.3.3",   // WalkOid
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
	data map[string]interface{}) (devices []*sdk.DeviceConfig, err error) {

	// Get the rack and board ids. Setup the location.
	rack, board, err := core.GetRackAndBoard(data)
	if err != nil {
		return nil, err
	}
	location := sdk.Location{
		Rack:  rack,
		Board: board,
	}

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(table.Rows); i++ {
		// upsInputFrequency
		// deviceData gets shimmed into the DeviceConfig for each synse device.
		// It varies slightly for each device below.
		deviceData := map[string]string{
			"info":       fmt.Sprintf("upsInputFrequency%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "2",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 2), // base_oid and integer column.
			"multiplier": ".1",                                  // Units are 0.1 Hertz
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device := sdk.DeviceConfig{
			Version:  "1",
			Type:     "frequency",
			Model:    model,
			Location: location,
			Data:     deviceData,
		}
		devices = append(devices, &device)

		// upsInputVoltage ----------------------------------------------------------
		deviceData = map[string]string{
			"info":       fmt.Sprintf("upsInputVoltage%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "3",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 3), // base_oid and integer column.
			// No multiplier needed. Units are RMS Volts.
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device2 := sdk.DeviceConfig{
			Version:  "1",
			Type:     "voltage",
			Model:    model,
			Location: location,
			Data:     deviceData,
		}
		devices = append(devices, &device2)

		// upsInputCurrent ----------------------------------------------------------
		deviceData = map[string]string{
			"info":       fmt.Sprintf("upsInputCurrent%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "4",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 4), // base_oid and integer column.
			"multiplier": ".1",                                  // Units are 0.1 RMS Amp
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device3 := sdk.DeviceConfig{
			Version:  "1",
			Type:     "current",
			Model:    model,
			Location: location,
			Data:     deviceData,
		}
		devices = append(devices, &device3)

		// upsInputTruePower --------------------------------------------------------
		deviceData = map[string]string{
			"info":       fmt.Sprintf("upsInputTruePower%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "5",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 5), // base_oid and integer column.
			// Output is in Watts. No multiplier needed.
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device4 := sdk.DeviceConfig{
			Version:  "1",
			Type:     "power",
			Model:    model,
			Location: location,
			Data:     deviceData,
		}
		devices = append(devices, &device4)
	}
	return devices, err
}
