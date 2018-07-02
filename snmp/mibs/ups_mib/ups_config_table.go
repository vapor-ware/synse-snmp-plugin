package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsConfigTable represents SNMP OID 1.3.6.1.2.1.33.1.9
type UpsConfigTable struct {
	*core.SnmpTable // base class
}

// NewUpsConfigTable constructs the UpsConfigTable.
func NewUpsConfigTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsConfigTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Config-Table", // Table Name
		".1.3.6.1.2.1.33.1.9",      // WalkOid
		[]string{ // Column Names
			"upsConfigInputVoltage",             // RMS Volts
			"upsConfigInputFreq",                // 0.1 Hertz
			"upsConfigOutputVoltage",            // RMS Volts
			"upsConfigOutputFreq",               // 0.1 Hertz
			"upsConfigOutputVA",                 // Volt Amps
			"upsConfigOutputPower",              // Watts
			"upsConfigLowBattTime",              // Minutes
			"upsConfigAudibleStatus",            // Disabled, Enabled, Muted
			"upsConfigLowVoltageTransferPoint",  // RMS Volts
			"upsConfigHighVoltageTransferPoint", // RMS Volts
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"",             // readableColumn
		true)           // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsConfigTable{SnmpTable: snmpTable}
	return table, nil
}
