package mibs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsOutputHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.4
type UpsOutputHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsOutputHeadersTable constructs the UpsOutputHeadersTable.
func NewUpsOutputHeadersTable(snmpServerBase *core.SnmpServerBase) (table *UpsOutputHeadersTable, err error) {
	var tableName = "UPS-MIB-UPS-Output-Headers-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.4"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsOutputSource",
			"upsOutputFrequency",
			"upsOutputNumLines",
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
	data map[string]interface{}) (devices []*config.DeviceProto, err error) {

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	// We will have "status-int", "status-string", and "frequency" device kinds.
	// There is probably a better way of doing this, but this just gets things to
	// where they need to be for now.
	statusIntProto := &config.DeviceProto{
		Type: "status-int",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
	}

	statusStringProto := &config.DeviceProto{
		Type: "status-string",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
	}

	frequencyProto := &config.DeviceProto{
		Type: "frequency",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
	}

	devices = []*config.DeviceProto{
		statusIntProto,
		statusStringProto,
		frequencyProto,
	}

	// This is always a single row table.

	// upsOutputSource
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
		"enumeration1": "other",
		"enumeration2": "none",
		"enumeration3": "normal",
		"enumeration4": "bypass",
		"enumeration5": "battery",
		"enumeration6": "booster",
		"enumeration7": "reducer",
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device := &config.DeviceInstance{
		Info: "upsOutputSource",
		Data: deviceData,
	}
	statusStringProto.Instances = append(statusStringProto.Instances, device)

	// upsOutputFrequency --------------------------------------------------------
	deviceData = map[string]interface{}{
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "0",
		"column":     "2",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 2), // base_oid and integer column.
		"multiplier": float32(0.1),                          // Units are 0.1 Hertz
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device = &config.DeviceInstance{
		Info: "upsOutputFrequency",
		Data: deviceData,
	}
	frequencyProto.Instances = append(frequencyProto.Instances, device)

	// upsOutputNumLines ---------------------------------------------------------
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
		Info: "upsOutputNumLines",
		Data: deviceData,
	}
	statusIntProto.Instances = append(statusIntProto.Instances, device)

	return devices, err
}
