package core

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// SnmpRow encapsulates one row of an SNMP table.
// The row may only contain one column.
type SnmpRow struct {
	// This is different than the parent table.
	// For a table row this is an oid to walk to retrieve all columns.
	// For a single row this is the oid to get to retrieve the single column.
	BaseOid string
	// Pointer back to the table for this row.
	Table *SnmpTable
	// SNMP data for the walked row.
	RowData []*ReadResult
}

// NewSnmpRow creates the SnmpRow structure.
func NewSnmpRow(baseOid string, table *SnmpTable, rowData []*ReadResult) (*SnmpRow, error) {
	// Arg checks.
	if table == nil {
		return nil, fmt.Errorf("NewSnmpRow. table is nil")
	}

	if rowData == nil {
		return nil, fmt.Errorf("NewSnmpRow. rowData is nil")
	}

	// Verify table column list length is the same as the length of rowData for sanity check.
	if len(table.ColumnList) != len(rowData) {
		return nil, fmt.Errorf(
			"Given %d column names and %d data columns. Unable to map unless equal",
			len(table.ColumnList), len(rowData))
	}

	return &SnmpRow{
		BaseOid: baseOid,
		Table:   table,
		RowData: rowData,
	}, nil
}

// Dump the SnmpRow to the debug log.
func (snmpRow *SnmpRow) Dump() {
	log.Debugf("Dumping row for SNMP table %v", snmpRow.Table.Name)
	log.Debugf("baseOid: %v", snmpRow.BaseOid)
	for i := 0; i < len(snmpRow.Table.ColumnList); i++ {
		log.Debugf("row[%v] = %v", snmpRow.Table.ColumnList[i], snmpRow.RowData[i])
	}
}
