package mibs

import (
	"fmt"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
)

// UpsInputTable represents SNMP OID .1.3.6.1.2.1.33.1.3.3
type UpsInputTable struct {
	*core.SnmpTable // base class
}

// NewUpsInputTable constructs the UpsInputTable.
func NewUpsInputTable(snmpServerBase *core.SnmpServerBase) (
	table *UpsInputTable, err error) {

	// Initialize the base.
	snmpTable, err := core.NewSnmpTable(
		"UPS-MIB-UPS-Input-Table", // Table Name
		".1.3.6.1.2.1.33.1.3.3",   // WalkOid
		[]string{ // Column Names
			"upsInputLineIndex", // MIB says not accessable. Have seen it in walks.
			"upsInputFrequency", // .1 Hertz
			"upsInputVoltage",   // RMS Volts
			"upsInputCurrent",   // .1 RMS Amp
			"upsInputTruePower", // Down with False Power! (Watts)
		},
		snmpServerBase, // snmpServer
		"1",            // rowBase
		"",             // indexColumn
		"2",            // readableColumn
		false)          // flattened table
	if err != nil {
		return nil, err
	}

	table = &UpsInputTable{SnmpTable: snmpTable}
	table.DevEnumerator = UpsInputTableDeviceEnumerator{table}
	return table, nil
}

// UpsInputTableDeviceEnumerator overrides the default SnmpTable device
// enumerator for the input table.
type UpsInputTableDeviceEnumerator struct {
	Table *UpsInputTable // Pointer back to the table.
}

// DeviceEnumerator overrides the default SnmpTable device enumerator.
func (enumerator UpsInputTableDeviceEnumerator) DeviceEnumerator(
	data map[string]interface{}) (devices []*sdk.DeviceConfig, err error) {

	// Get the rack and board ids. Setup the location.
	rack, board, err := core.GetRackAndBoard(data)
	if err != nil {
		return nil, err
	}

	table := enumerator.Table
	mib := table.Mib.(*UpsMib)
	model := mib.UpsIdentityTable.UpsIdentity.Model

	snmpDeviceConfigMap, err := table.SnmpServerBase.DeviceConfig.ToMap()
	if err != nil {
		return nil, err
	}

	locationName := "snmp-location"
	cfg := &sdk.DeviceConfig{
		SchemeVersion: sdk.SchemeVersion{Version: "1.0"},
		Locations: []*sdk.LocationConfig{
			{
				Name:  locationName,
				Rack:  &sdk.LocationData{Name: rack},
				Board: &sdk.LocationData{Name: board},
			},
		},
		Devices: []*sdk.DeviceKind{},
	}

	// We will have "frequency", "voltage", "current", and "power" device kinds.
	// There is probably a better way of doing this, but this just gets things to
	// where they need to be for now.
	frequencyKind := &sdk.DeviceKind{
		Name: "frequency",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "frequency"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	voltageKind := &sdk.DeviceKind{
		Name: "voltage",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "voltage"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	currentKind := &sdk.DeviceKind{
		Name: "current",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "current"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	powerKind := &sdk.DeviceKind{
		Name: "watts.power",
		Metadata: map[string]string{
			"model": model,
		},
		Outputs: []*sdk.DeviceOutput{
			{Type: "watts.power"},
		},
		Instances: []*sdk.DeviceInstance{},
	}

	cfg.Devices = []*sdk.DeviceKind{
		frequencyKind,
		voltageKind,
		currentKind,
		powerKind,
	}

	for i := 0; i < len(table.Rows); i++ {
		// upsInputFrequency
		// deviceData gets shimmed into the DeviceConfig for each synse device.
		// It varies slightly for each device below.
		deviceData := map[string]interface{}{
			"info":       fmt.Sprintf("upsInputFrequency%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "2",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 2), // base_oid and integer column.
			"multiplier": ".1",                                  // Units are 0.1 Hertz
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device := &sdk.DeviceInstance{
			Location: locationName,
			Data:     deviceData,
		}
		frequencyKind.Instances = append(frequencyKind.Instances, device)

		// upsInputVoltage ----------------------------------------------------------
		deviceData = map[string]interface{}{
			"info":       fmt.Sprintf("upsInputVoltage%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "3",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 3), // base_oid and integer column.
			// No multiplier needed. Units are RMS Volts.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &sdk.DeviceInstance{
			Location: locationName,
			Data:     deviceData,
		}
		voltageKind.Instances = append(voltageKind.Instances, device)

		// upsInputCurrent ----------------------------------------------------------
		deviceData = map[string]interface{}{
			"info":       fmt.Sprintf("upsInputCurrent%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "4",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 4), // base_oid and integer column.
			"multiplier": ".1",                                  // Units are 0.1 RMS Amp
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &sdk.DeviceInstance{
			Location: locationName,
			Data:     deviceData,
		}
		currentKind.Instances = append(currentKind.Instances, device)

		// upsInputTruePower --------------------------------------------------------
		deviceData = map[string]interface{}{
			"info":       fmt.Sprintf("upsInputTruePower%d", i),
			"base_oid":   table.Rows[i].BaseOid,
			"table_name": table.Name,
			"row":        fmt.Sprintf("%d", i),
			"column":     "5",
			"oid":        fmt.Sprintf(table.Rows[i].BaseOid, 5), // base_oid and integer column.
			// Output is in Watts. No multiplier needed.
		}
		deviceData, err = core.MergeMapStringInterface(snmpDeviceConfigMap, deviceData)
		if err != nil {
			return nil, err
		}

		device = &sdk.DeviceInstance{
			Location: locationName,
			Data:     deviceData,
		}
		powerKind.Instances = append(powerKind.Instances, device)
	}

	devices = append(devices, cfg)
	return devices, err
}
