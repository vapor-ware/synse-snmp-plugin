package core

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
)

// DeviceEnumeratorInterface is an interface that child classes can override in order to
// dynamically discover devices on a scan.
type DeviceEnumeratorInterface interface {
	DeviceEnumerator(data map[string]interface{}) ([]*config.DeviceProto, error)
}

// SnmpTable encapsulates a table of SNMP OIDs and data.
type SnmpTable struct {
	// The name of the table.
	Name string
	// The SNMP OID to walk the whole table.
	WalkOid string
	// Column names for the table.
	ColumnList []string
	// Pointer to the SnmpServer with the data for this row.
	SnmpServerBase *SnmpServerBase

	// If present each row is keyed at:
	// walk_oid + '.' + row_base + '.' + <column_index> + '.' + <row_index>
	RowBase string
	// This is appended to the WalkOid to get row indexes. The column may
	// not be readable however, in which case None is fine.
	IndexColumn string
	// Required when the index_column is marked not-accessible in the MIB.
	// Provides a way to tell how many rows are in the table.
	ReadableColumn string
	// A flattened table is really just a group of OIDs at the same level
	// and not truly an SNMP table. It's just simpler to read it in as a
	// single row table.
	FlattenedTable bool

	// The row data in the table.
	Rows []SnmpRow

	// Overrideable interface for device enumeration.
	DevEnumerator DeviceEnumeratorInterface

	// Pointer back to the SnmpMib derrived class. This is not in the constructor
	// because things would get difficult to initialize. Initialized in the
	// SnmpMib derrived class constructor.
	Mib interface{}
}

// DeviceEnumerator is the default SnmpTable Implementation of DeviceEnumerator that returns no devices.
func (snmpTable *SnmpTable) DeviceEnumerator(data map[string]interface{}) ([]*config.DeviceProto, error) {
	return []*config.DeviceProto{}, nil
}

// EnumerateDevices calls the specific implementation for enumerating devices.
func (snmpTable *SnmpTable) EnumerateDevices(data map[string]interface{}) ([]*config.DeviceProto, error) {
	return snmpTable.DevEnumerator.DeviceEnumerator(data)
}

// SnmpTableDefaultEnumerator is the default device enumerator.
type SnmpTableDefaultEnumerator struct{}

// DeviceEnumerator is the default enumerator which returns no devices.
func (enumerator SnmpTableDefaultEnumerator) DeviceEnumerator(data map[string]interface{}) ([]*config.DeviceProto, error) {
	return []*config.DeviceProto{}, nil
}

// NewSnmpTable creates the SnmpTable structure.
func NewSnmpTable(
	name string,
	walkOid string,
	columnList []string,
	snmpServerBase *SnmpServerBase,
	rowBase string,
	indexColumn string,
	readableColumn string,
	flattenedTable bool) (*SnmpTable, error) {

	// Arg Checks
	if name == "" {
		return nil, fmt.Errorf("NewSnmpTable. name is empty")
	}
	if !strings.HasPrefix(walkOid, ".") {
		return nil, fmt.Errorf(
			"NewSnmpTable. walkOid must start with a period, walkOid: %v", walkOid)
	}
	if columnList == nil {
		return nil, fmt.Errorf("NewSnmpTable. columnList is nil")
	}
	if snmpServerBase == nil {
		return nil, fmt.Errorf("NewSnmpTable. snmpServerBase is nil")
	}

	// Create the table. This should load the rows.
	snmpTable := &SnmpTable{
		Name:           name,
		WalkOid:        walkOid,
		ColumnList:     columnList,
		SnmpServerBase: snmpServerBase,
		RowBase:        rowBase,
		IndexColumn:    indexColumn,
		ReadableColumn: readableColumn,
		FlattenedTable: flattenedTable,
		DevEnumerator:  SnmpTableDefaultEnumerator{},
	}

	err := snmpTable.Load()
	if err != nil {
		return nil, err
	}
	return snmpTable, nil
}

// Dump to the log as CSV. Also to console.
func (snmpTable *SnmpTable) Dump() {
	// Header
	log.Debugf("Dumping %v table. %d rows. walk oid: %v",
		snmpTable.Name, len(snmpTable.Rows), snmpTable.WalkOid)
	fmt.Printf("Dumping %v table. %d rows. walk oid: %v\n",
		snmpTable.Name, len(snmpTable.Rows), snmpTable.WalkOid)
	// Column list
	log.Debugf("%v", strings.Join(snmpTable.ColumnList, ","))
	fmt.Printf("%v\n", strings.Join(snmpTable.ColumnList, ","))

	for i := 0; i < len(snmpTable.Rows); i++ { // for each row in the table
		var data []string
		row := snmpTable.Rows[i]
		for j := 0; j < len(row.RowData); j++ {
			cell := row.RowData[j]
			data = append(data, fmt.Sprintf("%v", cell.Data))
		}
		log.Debugf("%v\n", strings.Join(data, ","))
		fmt.Printf("%v\n", strings.Join(data, ","))
	}
}

// Get the row with the given base OID from the table or nil if not present.
// This is just a get from the cache. It is not a get from the SNMP server.
func (snmpTable *SnmpTable) Get(baseOid string) *SnmpRow {
	for i := 0; i < len(snmpTable.Rows); i++ {
		log.Debugf("baseOid: %v", snmpTable.Rows[i].BaseOid)
		if snmpTable.Rows[i].BaseOid == baseOid {
			log.Debugf("found row: %v", snmpTable.Rows[i])
			return &snmpTable.Rows[i]
		}
	}
	return nil
}

// Load the data from the SNMP Server.
// Walk the walk_oid on the SNMP server. Translate the data to SnmpRows.
func (snmpTable *SnmpTable) Load() error {
	// SNMP Walk the table.
	rawResults, err := snmpTable.SnmpServerBase.SnmpClient.Walk(snmpTable.WalkOid)
	if err != nil {
		return err
	}
	err = snmpTable.translate(rawResults)
	return err
}

// Unload cached row data once we're done with it.
func (snmpTable *SnmpTable) Unload() {
	snmpTable.Rows = nil
	log.Debugf("Unloaded SnmpTable %v", snmpTable.Name)
}

// Update the table by removing the row with the same base_oid as row,
// then adding row from the parameter list.
// row: The row to update.
// NOTE: This is an upsert.
func (snmpTable *SnmpTable) Update(row *SnmpRow) {
	// Delete the existing SNMP row in the variable table by base_oid.
	log.Debugf("before delete row count %d", len(snmpTable.Rows))
	for i := 0; i < len(snmpTable.Rows); i++ {
		if snmpTable.Rows[i].BaseOid == row.BaseOid {
			// Delete the row.
			snmpTable.Rows = append(snmpTable.Rows[:i], snmpTable.Rows[i+1:]...)
			break
		}
	}
	log.Debugf("after delete row count %d", len(snmpTable.Rows))

	// Add the new row.
	snmpTable.Rows = append(snmpTable.Rows, *row)
}

// UpdateCell updates the table data. Used on successful write.
// baseOid: The base oid of the row to update.
// index: The 1 based column index to update.
// data: The data for the update.
func (snmpTable *SnmpTable) UpdateCell(
	baseOid string, index int, data interface{}) (err error) {
	for i := 0; i < len(snmpTable.Rows); i++ {
		row := snmpTable.Rows[i]
		if row.BaseOid == baseOid {
			row.RowData[index-1].Data = data
			return nil
		}
	}
	return fmt.Errorf("Unable to update table %v, baseOid %v, index %d",
		snmpTable.Name, baseOid, index)
}

// getRowIndexes gets the index portion of the OID so that we can line up the rows.
func (snmpTable *SnmpTable) getRowIndexes(tableData []ReadResult) []string {
	var rowIndexes []string

	prefix := snmpTable.WalkOid + "."
	if snmpTable.RowBase != "" {
		prefix += snmpTable.RowBase + "."
	}

	if snmpTable.IndexColumn != "" {
		// We can read the index column to get the row indexes.
		prefix += snmpTable.IndexColumn + "."
		//  Get all values where the key starts with the prefix.
		for i := 0; i < len(tableData); i++ {
			if strings.HasPrefix(tableData[i].Oid, prefix) {
				rowIndexes = append(rowIndexes, tableData[i].Oid)
			}
		}
	} else {
		// We can't read the index column. Read the column passed in at
		// snmpTable.ReadableColumn.
		if snmpTable.ReadableColumn != "" {
			prefix += snmpTable.ReadableColumn + "."
		}
		var oids []string
		for i := 0; i < len(tableData); i++ {
			if strings.HasPrefix(tableData[i].Oid, prefix) {
				oids = append(oids, tableData[i].Oid)
			}
		}
		for i := 0; i < len(oids); i++ {
			rowIndex := oids[i][len(prefix):]
			rowIndexes = append(rowIndexes, rowIndex)
		}
	}
	return rowIndexes
}

// getData is a helper to get ReadResult from a slice.
func getData(oid string, results []ReadResult) *ReadResult {

	for i := 0; i < len(results); i++ {
		if results[i].Oid == oid {
			return &results[i]
		}
	}
	// This can happen on unreadable columns. Oid is non-nil, data are nil.
	return &ReadResult{
		Oid:  oid,
		Data: nil,
	}
}

// Translate into a structure of SnmpRow.
func (snmpTable *SnmpTable) translate(tableData []ReadResult) error {
	snmpTable.Rows = *new([]SnmpRow)
	rowIndexes := snmpTable.getRowIndexes(tableData)

	if len(rowIndexes) == 0 {
		return nil // No rows.
	}

	if !snmpTable.FlattenedTable {
		for i := 0; i < len(rowIndexes); i++ {
			columnIndex := 1
			var rowData []*ReadResult

			baseOid := snmpTable.WalkOid + "." + snmpTable.RowBase + ".%d." + rowIndexes[i]

			for j := 0; j < len(snmpTable.ColumnList); j++ {
				dataOid := fmt.Sprintf(baseOid, columnIndex)
				data := getData(dataOid, tableData)
				rowData = append(rowData, data)
				columnIndex++
			}
			row, err := NewSnmpRow(baseOid, snmpTable, rowData)
			if err != nil {
				return err
			}
			snmpTable.Rows = append(snmpTable.Rows, *row)
		}
	} else {
		baseOid := snmpTable.WalkOid + ".%d.0"
		columnIndex := 1
		var rowData []*ReadResult
		for i := 0; i < len(snmpTable.ColumnList); i++ {
			dataOid := fmt.Sprintf(baseOid, columnIndex)
			data := getData(dataOid, tableData)
			rowData = append(rowData, data)
			columnIndex++
		}
		row, err := NewSnmpRow(baseOid, snmpTable, rowData)
		if err != nil {
			return err
		}
		snmpTable.Rows = append(snmpTable.Rows, *row)
	}

	return nil
}
