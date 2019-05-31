package mibs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
)

// UpsMib is the class for all SNMP operations on UPS-MIB, rfc 1628.
type UpsMib struct {
	*core.SnmpMib // base class

	// Tables defined in this MIB
	UpsIdentityTable        *UpsIdentityTable
	UpsBatteryTable         *UpsBatteryTable
	UpsInputHeadersTable    *UpsInputHeadersTable
	UpsInputTable           *UpsInputTable
	UpsOutputHeadersTable   *UpsOutputHeadersTable
	UpsOutputTable          *UpsOutputTable
	UpsBypassHeadersTable   *UpsBypassHeadersTable
	UpsBypassTable          *UpsBypassTable
	UpsAlarmsHeadersTable   *UpsAlarmsHeadersTable
	UpsAlarmsTable          *UpsAlarmsTable
	UpsWellKnownAlarmsTable *UpsWellKnownAlarmsTable
	UpsTestHeadersTable     *UpsTestHeadersTable
	UpsWellKnownTestsTable  *UpsWellKnownTestsTable
	UpsControlTable         *UpsControlTable
	UpsConfigTable          *UpsConfigTable
	UpsCompliancesTable     *UpsCompliancesTable
	UpsSubsetGroupsTable    *UpsSubsetGroupsTable
	UpsBasicGroupsTable     *UpsBasicGroupsTable
	UpsFullGroupsTable      *UpsFullGroupsTable
}

// NewUpsMib constructs the UpsMib.
func NewUpsMib(server *core.SnmpServerBase) (upsMib *UpsMib, err error) { // nolint: gocyclo
	log.Debugf("Initializing UpsMib")

	// Arg checks.
	if server == nil {
		return nil, fmt.Errorf("NewUpsMib, server is nil")
	}

	// Initialize Tables.
	upsIdentityTable, err := NewUpsIdentityTable(server)
	if err != nil {
		return nil, err
	}

	upsBatteryTable, err := NewUpsBatteryTable(server)
	if err != nil {
		return nil, err
	}

	upsInputHeadersTable, err := NewUpsInputHeadersTable(server)
	if err != nil {
		return nil, err
	}

	upsInputTable, err := NewUpsInputTable(server)
	if err != nil {
		return nil, err
	}

	upsOutputHeadersTable, err := NewUpsOutputHeadersTable(server)
	if err != nil {
		return nil, err
	}

	upsOutputTable, err := NewUpsOutputTable(server)
	if err != nil {
		return nil, err
	}

	upsBypassHeadersTable, err := NewUpsBypassHeadersTable(server)
	if err != nil {
		return nil, err
	}

	upsBypassTable, err := NewUpsBypassTable(server)
	if err != nil {
		return nil, err
	}

	upsAlarmsHeadersTable, err := NewUpsAlarmsHeadersTable(server)
	if err != nil {
		return nil, err
	}

	upsAlarmsTable, err := NewUpsAlarmsTable(server)
	if err != nil {
		return nil, err
	}

	upsWellKnownAlarmsTable, err := NewUpsWellKnownAlarmsTable(server)
	if err != nil {
		return nil, err
	}

	upsTestHeadersTable, err := NewUpsTestHeadersTable(server)
	if err != nil {
		return nil, err
	}

	upsWellKnownTestsTable, err := NewUpsWellKnownTestsTable(server)
	if err != nil {
		return nil, err
	}

	upsControlTable, err := NewUpsControlTable(server)
	if err != nil {
		return nil, err
	}

	upsConfigTable, err := NewUpsConfigTable(server)
	if err != nil {
		return nil, err
	}

	upsCompliancesTable, err := NewUpsCompliancesTable(server)
	if err != nil {
		return nil, err
	}

	upsSubsetGroupsTable, err := NewUpsSubsetGroupsTable(server)
	if err != nil {
		return nil, err
	}

	upsBasicGroupsTable, err := NewUpsBasicGroupsTable(server)
	if err != nil {
		return nil, err
	}

	upsFullGroupsTable, err := NewUpsFullGroupsTable(server)
	if err != nil {
		return nil, err
	}

	// Initialize the base class.
	snmpMib, err := core.NewSnmpMib(
		"UPS-MIB",
		[]*core.SnmpTable{
			upsIdentityTable.SnmpTable,
			upsBatteryTable.SnmpTable,
			upsInputHeadersTable.SnmpTable,
			upsInputTable.SnmpTable,
			upsOutputHeadersTable.SnmpTable,
			upsOutputTable.SnmpTable,
			upsBypassHeadersTable.SnmpTable,
			upsBypassTable.SnmpTable,
			upsAlarmsHeadersTable.SnmpTable,
			upsAlarmsTable.SnmpTable,
			upsWellKnownAlarmsTable.SnmpTable,
			upsTestHeadersTable.SnmpTable,
			upsWellKnownTestsTable.SnmpTable,
			upsControlTable.SnmpTable,
			upsConfigTable.SnmpTable,
			upsCompliancesTable.SnmpTable,
			upsSubsetGroupsTable.SnmpTable,
			upsBasicGroupsTable.SnmpTable,
			upsFullGroupsTable.SnmpTable,
		})
	if err != nil {
		return nil, err
	}

	// Initialize class.
	upsMib = &UpsMib{SnmpMib: snmpMib} // base mib class
	// Tables
	upsMib.UpsIdentityTable = upsIdentityTable
	upsMib.UpsBatteryTable = upsBatteryTable
	upsMib.UpsInputHeadersTable = upsInputHeadersTable
	upsMib.UpsInputTable = upsInputTable
	upsMib.UpsOutputHeadersTable = upsOutputHeadersTable
	upsMib.UpsOutputTable = upsOutputTable
	upsMib.UpsBypassHeadersTable = upsBypassHeadersTable
	upsMib.UpsBypassTable = upsBypassTable
	upsMib.UpsAlarmsHeadersTable = upsAlarmsHeadersTable
	upsMib.UpsAlarmsTable = upsAlarmsTable
	upsMib.UpsWellKnownAlarmsTable = upsWellKnownAlarmsTable
	upsMib.UpsTestHeadersTable = upsTestHeadersTable
	upsMib.UpsWellKnownTestsTable = upsWellKnownTestsTable
	upsMib.UpsControlTable = upsControlTable
	upsMib.UpsConfigTable = upsConfigTable
	upsMib.UpsCompliancesTable = upsCompliancesTable
	upsMib.UpsSubsetGroupsTable = upsSubsetGroupsTable
	upsMib.UpsBasicGroupsTable = upsBasicGroupsTable
	upsMib.UpsFullGroupsTable = upsFullGroupsTable

	// Update mib pointer for each table.
	upsMib.UpsIdentityTable.Mib = upsMib
	upsMib.UpsBatteryTable.Mib = upsMib
	upsMib.UpsInputHeadersTable.Mib = upsMib
	upsMib.UpsInputTable.Mib = upsMib
	upsMib.UpsOutputHeadersTable.Mib = upsMib
	upsMib.UpsOutputTable.Mib = upsMib
	upsMib.UpsBypassHeadersTable.Mib = upsMib
	upsMib.UpsBypassTable.Mib = upsMib
	upsMib.UpsAlarmsHeadersTable.Mib = upsMib
	upsMib.UpsAlarmsTable.Mib = upsMib
	upsMib.UpsWellKnownAlarmsTable.Mib = upsMib
	upsMib.UpsTestHeadersTable.Mib = upsMib
	upsMib.UpsWellKnownTestsTable.Mib = upsMib
	upsMib.UpsControlTable.Mib = upsMib
	upsMib.UpsConfigTable.Mib = upsMib
	upsMib.UpsCompliancesTable.Mib = upsMib
	upsMib.UpsSubsetGroupsTable.Mib = upsMib
	upsMib.UpsBasicGroupsTable.Mib = upsMib
	upsMib.UpsFullGroupsTable.Mib = upsMib

	log.Debugf("Initialized UpsMib")
	return upsMib, nil
}
