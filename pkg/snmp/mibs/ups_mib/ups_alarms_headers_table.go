package mibs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsAlarmsHeadersTable represents SNMP OID .1.3.6.1.2.1.33.1.6
// This is just the alarms present column. OID: .1.3.6.1.2.1.33.1.6.1.0
type UpsAlarmsHeadersTable struct {
	*core.SnmpTable // base class
}

// NewUpsAlarmsHeadersTable constructs the UpsAlarmsHeadersTable.
func NewUpsAlarmsHeadersTable(snmpServerBase *core.SnmpServerBase) (table *UpsAlarmsHeadersTable, err error) {
	var tableName = "UPS-MIB-UPS-Alarms-Headers-Table"
	var walkOid = ".1.3.6.1.2.1.33.1.6"

	log.WithFields(log.Fields{
		"name": tableName,
		"oid":  walkOid,
	}).Debug("[snmp] creating new table")

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		tableName,
		walkOid,
		[]string{ // Column Names
			"upsAlarmsPresent", // The present number of active alarm conditions.
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

	table = &UpsAlarmsHeadersTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsAlarmsHeadersTableDeviceEnumerator{table}
	return table, nil
}

// UpsAlarmsHeadersTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the alarms headers table.
type UpsAlarmsHeadersTableDeviceEnumerator struct {
	Table *UpsAlarmsHeadersTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsAlarmsHeadersTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*config.DeviceProto, err error) {

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return
	}

	statusProto := &config.DeviceProto{
		Type: "status",
		Context: map[string]string{
			"model": model,
		},
		Instances: []*config.DeviceInstance{},
	}

	devices = []*config.DeviceProto{
		statusProto,
	}

	// This is always a single row table.

	// upsAlarmsPresent ---------------------------------------------------------
	deviceData := map[string]interface{}{
		"base_oid":   table.Rows[0].BaseOid,
		"table_name": table.Name,
		"row":        "1",
		"column":     "0",
		"oid":        fmt.Sprintf(table.Rows[0].BaseOid, 0), // base_oid and integer column.
	}
	deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
	if err != nil {
		return nil, err
	}

	device := &config.DeviceInstance{
		Info: "upsAlarmsPresent",
		Data: deviceData,
	}
	statusProto.Instances = append(statusProto.Instances, device)

	return
}
