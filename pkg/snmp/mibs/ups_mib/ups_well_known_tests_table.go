package mibs

import (
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsWellKnownTestsTable represents SNMP OID .1.3.6.1.2.1.33.1.7.7
type UpsWellKnownTestsTable struct {
	*core.SnmpTable // base class
}

// NewUpsWellKnownTestsTable constructs the UpsWellKnownTestsTable.
func NewUpsWellKnownTestsTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsWellKnownTestsTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Well-Known-Tests-Table", // Table Name
		".1.3.6.1.2.1.33.1.7.7",              // WalkOid
		[]string{ // Column Names
			"upsTestNoTestsInitiated",
			"upsTestAbortTestInProgress",
			"upsTestGeneralSystemsTest",
			"upsTestQuickBatteryTest",
			"upsTestDeepBatteryCalibration",
		},
		snmpServerBase, // snmpServer
		"",             // rowBase
		"",             // indexColumn
		"1",            // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsWellKnownTestsTable{SnmpTable: snmpTable}
	return table, nil
}
