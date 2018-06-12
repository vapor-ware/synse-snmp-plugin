package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-sdk/sdk/logger"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsIdentity contains identification information for a UPS.
// These come in as binary data from a single SnmpRow. This struct is a helper/
// translator to provide string data.
type UpsIdentity struct {
	Manufacturer         string
	Model                string
	UpsSoftwareVersion   string
	AgentSoftwareVersion string
	Name                 string
	AttachedDevices      string
}

// UpsIdentityTable represts SNMP OID .1.3.6.1.2.1.33.1.1
type UpsIdentityTable struct {
	*core.SnmpTable              // base class
	UpsIdentity     *UpsIdentity // Identity information.
}

// NewUpsIdentityTable constructs the UpsIdentityTable.
func NewUpsIdentityTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsIdentityTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Identity-Table", // Table Name
		".1.3.6.1.2.1.33.1.1",        // WalkOid
		[]string{ // Column Names
			"upsIdentManufacturer",
			"upsIdentModel",
			"upsIdentUPSSoftwareVersion",
			"upsIdentAgentSoftwareVersion",
			"upsIdentName",
			"upsIdentAttachedDevices",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsIdentityTable{SnmpTable: snmpTable}
	table.UpsIdentity = table.loadIdentity()
	table.DevEnumerator = UpsIdentityTableDeviceEnumerator{table}
	return table, nil
}

// loadIdentity loads the UpsIdentity data.
func (table *UpsIdentityTable) loadIdentity() *UpsIdentity { // nolint: gocyclo
	// Defaults are empty strings.
	manufacturer := ""
	model := ""
	upsSoftwareVersion := ""
	agentSoftwareVersion := ""
	name := ""
	attachedDevices := ""

	// Need these variable declarations before the gotos.
	var snmpRow core.SnmpRow
	var field string
	var ok bool

	if table == nil || len(table.Rows) < 1 {
		logger.Warn("No identity information.")
		goto end
	}

	snmpRow = table.Rows[0]
	if snmpRow.RowData == nil || len(snmpRow.RowData) < 6 {
		logger.Warn("No identity information.")
		goto end
	}

	// Get each field by column from the row.
	field, ok = snmpRow.RowData[0].Data.(string)
	if ok {
		manufacturer = field
	}

	field, ok = snmpRow.RowData[1].Data.(string)
	if ok {
		model = field
	}

	field, ok = snmpRow.RowData[2].Data.(string)
	if ok {
		upsSoftwareVersion = field
	}

	field, ok = snmpRow.RowData[3].Data.(string)
	if ok {
		agentSoftwareVersion = field
	}

	field, ok = snmpRow.RowData[4].Data.(string)
	if ok {
		name = field
	}

	field, ok = snmpRow.RowData[5].Data.(string)
	if ok {
		attachedDevices = field
	}

end:
	return &UpsIdentity{
		Manufacturer:         manufacturer,
		Model:                model,
		UpsSoftwareVersion:   upsSoftwareVersion,
		AgentSoftwareVersion: agentSoftwareVersion,
		Name:                 name,
		AttachedDevices:      attachedDevices,
	}
}

// UpsIdentityTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the input headers table.
type UpsIdentityTableDeviceEnumerator struct {
	Table *UpsIdentityTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsIdentityTableDeviceEnumerator) DeviceEnumerator(
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

	// Pull out the table, mib, device model, SNMP DeviceConfig.
	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	// This is always a single row table.

	// upsIdentManufacturer
	// deviceData gets shimmed into the DeviceConfig for each synse device.
	// It varies slightly for each device below.
	deviceData := map[string]string{
		"info":       "upsIdentManufacturer",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "1",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 1), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	// Create the synse device.
	device := config.DeviceConfig{
		Version:  "1",
		Type:     "identity",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device)

	// upsIdentModel
	deviceData = map[string]string{
		"info":       "upsIdentModel",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "2",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 2), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device2 := config.DeviceConfig{
		Version:  "1",
		Type:     "identity",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device2)

	// upsIdentUPSSoftwareVersion
	deviceData = map[string]string{
		"info":       "upsIdentSoftwareVersion",
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
		Type:     "identity",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device3)

	// upsIdentAgentSoftwareVersion ----------------------------------------------
	deviceData = map[string]string{
		"info":       "upsIdentAgentSoftwareVersion",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "4",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 4), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device4 := config.DeviceConfig{
		Version:  "1",
		Type:     "identity",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device4)

	// upsIdentName ---------------------------------------------------------------
	deviceData = map[string]string{
		"info":       "upsIdentName",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "5",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 5), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device5 := config.DeviceConfig{
		Version:  "1",
		Type:     "identity",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device5)

	// upsIdentAttachedDevices ----------------------------------------------------
	deviceData = map[string]string{
		"info":       "upsIdentAttachedDevices",
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "6",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 6), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringString(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device6 := config.DeviceConfig{
		Version:  "1",
		Type:     "identity",
		Model:    model,
		Location: location,
		Data:     deviceData,
	}
	devices = append(devices, &device6)

	return devices, err
}
