package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsBypassTable represts SNMP OID .1.3.6.1.2.1.33.1.5.3
type UpsBypassTable struct {
	*core.SnmpTable // base class
}

// NewUpsBypassTable constructs the UpsBypassTable.
func NewUpsBypassTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsBypassTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Bypass-Table", // Table Name
		".1.3.6.1.2.1.33.1.2",      // WalkOid
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
		false)          // flattened table
	if err != nil {
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

	// Pull out the table, mib, device model, SNMP DeviceConfig.
	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(table.Rows); i++ {
		// upsBypassVoltage
		// deviceData gets shimmed into the DeviceConfig for each synse device.
		// It varies slightly for each device below.
		deviceData := map[string]string{
			"info":       fmt.Sprintf("upsBypassVoltage%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "2",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 2), // base_oid and integer column.
			// No multiplier needed. Units are RMS Volts.
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device := sdk.DeviceConfig{
			Version:  "1",
			Type:     "voltage",
			Model:    model,
			Location: location,
			Data:     deviceData,
		}
		devices = append(devices, &device)

		// upsBypassCurrent ---------------------------------------------------------
		deviceData = map[string]string{
			"info":       fmt.Sprintf("upsBypassCurrent%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "3",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 3), // base_oid and integer column.
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device2 := sdk.DeviceConfig{
			Version:  "1",
			Type:     "current",
			Model:    model,
			Location: location,
			Data:     deviceData,
		}
		devices = append(devices, &device2)

		// upsBypassPower --------------------------------------------------------------
		deviceData = map[string]string{
			"info":       fmt.Sprintf("upsBypassPower%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "4",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 4), // base_oid and integer column.
			// Output is in Watts. No multiplier needed.
		}
		deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device3 := sdk.DeviceConfig{
			Version:  "1",
			Type:     "power",
			Model:    model,
			Location: location,
			Data:     deviceData,
		}
		devices = append(devices, &device3)
	}
	return devices, err
}
