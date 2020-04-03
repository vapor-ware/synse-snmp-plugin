package devices

import "github.com/vapor-ware/synse-snmp-plugin/exp/mibs"

// BatteryDevices contains the definitions of all the "upsBattery" objects
// in the UPS-MIB definition.
//
// See UPS-MIB 1.3.6.1.2.1.33.1.2
// http://www.oidview.com/mibs/0/UPS-MIB.html
var BatteryDevices = []*mibs.SnmpDevice{
	&UpsBatteryStatus,
	&UpsSecondsOnBattery,
	&UpsEstimatedMinutesRemaining,
	&UpsEstimatedChargeRemaining,
	&UpsBatteryVoltage,
	&UpsBatteryCurrent,
	&UpsBatteryTemperature,
}

// UPS-MIB upsBattery device definitions.
var (
	// TODO (etd): this should be an enumeration of:
	//   unknown, batteryNormal, batteryLow, batteryDepleted
	UpsBatteryStatus = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.2.1.0",
		Info:    "upsBatteryStatus",
		Type:    "status",
		Output:  "status",
		Handler: "read-only",
	}

	UpsSecondsOnBattery = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.2.2.0",
		Info:    "upsSecondsOnBattery",
		Type:    "status", // FIXME (etd): units are in seconds
		Output:  "status",
		Handler: "read-only",
	}

	UpsEstimatedMinutesRemaining = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.2.3.0",
		Info:    "upsEstimatedMinutesRemaining",
		Type:    "status", // FIXME (etd): units are in minutes
		Output:  "status",
		Handler: "read-only",
	}

	UpsEstimatedChargeRemaining = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.2.4.0",
		Info:    "upsEstimatedChargeRemaining",
		Type:    "percent",
		Output:  "percentage",
		Handler: "read-only",
	}

	UpsBatteryVoltage = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.2.5.0",
		Info:    "upsBatteryVoltage",
		Type:    "voltage", // FIXME (etd): units are "0.1 Volt DC"
		Output:  "voltage",
		Handler: "read-only",
	}

	UpsBatteryCurrent = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.2.6.0",
		Info:    "upsBatteryCurrent",
		Type:    "current", // FIXME (etd): units are "0.1 Amp DC"
		Output:  "electric-current",
		Handler: "read-only",
	}

	UpsBatteryTemperature = mibs.SnmpDevice{
		OID:     "1.3.6.1.2.1.33.1.2.7.0",
		Info:    "upsBatteryTemperature",
		Type:    "temperature",
		Output:  "temperature",
		Handler: "read-only",
	}
)
