package core

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
)

// SnmpMib is the base class for specific SNMP MIB implementations, a collection
// of SnmpTable
type SnmpMib struct {
	// The name of the MIB
	Name string
	// The tables that this MIB defines.
	Tables []*SnmpTable
}

// NewSnmpMib creates the SnmpMib structure.
func NewSnmpMib(name string, snmpTables []*SnmpTable) (*SnmpMib, error) {
	if name == "" {
		return nil, fmt.Errorf("NewSnmpMib. name is empty")
	}
	if len(snmpTables) == 0 {
		return nil, fmt.Errorf("NewSnmpMib. No snmpTables")
	}

	// Create the structure.
	mib := &SnmpMib{
		Name:   name,
		Tables: snmpTables,
	}

	// Initialize mib pointer for each table.
	for i := 0; i < len(snmpTables); i++ {
		snmpTables[i].Mib = mib
	}

	return mib, nil
}

// Dump all tables in the MIB to the log as CSV.
func (snmpMib *SnmpMib) Dump() {
	log.Debugf("Dumping SnmpMib %v. %d tables", snmpMib.Name, len(snmpMib.Tables))
	fmt.Printf("Dumping SnmpMib %+v\n", snmpMib.Name)

	for i := 0; i < len(snmpMib.Tables); i++ {
		snmpMib.Tables[i].Dump()
	}
	log.Debugf("End SnmpMib dump %v", snmpMib.Name)
}

// EnumerateDevices enumerates all synse devices supported by the mib.
func (snmpMib *SnmpMib) EnumerateDevices(data map[string]interface{}) (devices []*config.DeviceProto, err error) {
	for _, table := range snmpMib.Tables {
		deviceSet, err := table.DevEnumerator.DeviceEnumerator(data)
		if err != nil {
			return nil, err
		}
		devices = append(devices, deviceSet...)
	}
	return devices, nil
}

// Load all tables defined for the MIB.
func (snmpMib *SnmpMib) Load() error {
	for i := 0; i < len(snmpMib.Tables); i++ {
		err := snmpMib.Tables[i].Load()
		if err != nil {
			return err
		}
	}
	return nil
}

// Unload all tables defined for the MIB.
func (snmpMib *SnmpMib) Unload() {
	for i := 0; i < len(snmpMib.Tables); i++ {
		snmpMib.Tables[i].Unload()
	}
}
