package mibs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUpsAlarmsInfo tests that the device info string for the SNMP OIDs in the
// upsAlarmsTable is the same as the name in upsWellKnownAlarms (OIDs under
// .1.3.6.1.2.1.33.1.6.3)
func TestUpsAlarmsInfo(t *testing.T) {

	var expected = []string{
		"upsAlarmBatteryBad",
		"upsAlarmOnBattery",
		"upsAlarmLowBattery",
		"upsAlarmDepletedBattery",
		"upsAlarmTempBad",
		"upsAlarmInputBad",
		"upsAlarmOutputBad",
		"upsAlarmOutputOverload",
		"upsAlarmOnBypass",
		"upsAlarmBypassBad",
		"upsAlarmOutputOffAsRequested",
		"upsAlarmUpsOffAsRequested",
		"upsAlarmChargerFailed",
		"upsAlarmUpsOutputOff",
		"upsAlarmUpsSystemOff",
		"upsAlarmFanFailure",
		"upsAlarmFuseFailure",
		"upsAlarmGeneralFault",
		"upsAlarmDiagnosticTestFailed",
		"upsAlarmCommunicationsLost",
		"upsAlarmAwaitingPower",
		"upsAlarmShutdownPending",
		"upsAlarmShutdownImminent",
		"upsAlarmTestInProgress",
	}

	// Order matters.
	assert.True(t, reflect.DeepEqual(expected, upsAlarmsInfo))
}
