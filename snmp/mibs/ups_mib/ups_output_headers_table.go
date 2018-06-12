package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsOutputHeadersTable represts SNMP OID .1.3.6.1.2.1.33.1.4
type UpsOutputHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsOutputHeadersTable constructs the UpsOutputHeadersTable.
func NewUpsOutputHeadersTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsOutputHeadersTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Output-Headers-Table", // Table Name
		".1.3.6.1.2.1.33.1.4",              // WalkOid
		[]string{ // Column Names
			"upsOutputSource",
			"upsOutputFrequency",
			"upsOutputNumLines",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsOutputHeadersTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsOutputHeadersTableDeviceEnumerator{table}
	return table, nil
}

// UpsOutputHeadersTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the output headers table.
type UpsOutputHeadersTableDeviceEnumerator struct {
	Table *UpsOutputHeadersTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsOutputHeadersTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceConfig, err error) {

	// Get the rack and board ids. Setup the location.
	rack, board, err := core.GetRackAndBoard(data)
	if err != nil {
		return nil, err
	}
	location := config.Location{
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

	// This is always a single row table.

	// upsOutputSource
	// deviceData gets shimmed into the DeviceConfig for each synse device.
	// It varies slightly for each device below.
	deviceData := map[string]string{
		"info":       "upsOutputSource",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "1",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 1), // base_oid and integer column.
		// This is an enumeration. We need to translate the integer we read to a string.
		"enumeration": "true", // Defines that this is an enumeration.
		// Enumeration data. For now we have map[string]string to work with so the
		// key is fmt.Sprintf("enumeration%d", reading).
		"enumeration1": "other",
		"enumeration2": "none",
		"enumeration3": "normal",
		"enumeration4": "bypass",
		"enumeration5": "battery",
		"enumeration6": "booster",
		"enumeration7": "reducer",
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device := config.DeviceConfig{
		Version:  "1",
		Type:     "status",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device)

	// upsOutputFrequency --------------------------------------------------------
	deviceData = map[string]string{
		"info":       "upsOutputFrequency",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "2",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 2), // base_oid and integer column.
		"multiplier": ".1",                                  // Units are 0.1 Hertz
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device2 := config.DeviceConfig{
		Version:  "1",
		Type:     "frequency",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device2)

	// upsOutputNumLines ---------------------------------------------------------
	deviceData = map[string]string{
		"info":       "upsOutputNumLines",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "3",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 3), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device3 := config.DeviceConfig{
		Version:  "1",
		Type:     "status",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device3)

	return devices, err
}
