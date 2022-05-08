package servers

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
	mibs "github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/mibs/ups_mib"
)

// PxgmsUps represents the PXGMS UPS + EATON 93PM SNMP Server.
type PxgmsUps struct {
	*SnmpServer // base class.
}

// NewPxgmsUps creates the PxgmsUps structure.
// Sample data that works with the emulator:
// contextName:public
// endpoint:127.0.0.1
// userName:simulator
// privacyProtocol:AES
// privacyPassphrase:privatus
// port:1024
// authenticationProtocol:SHA
// authenticationPassphrase:auctoritas
// model:PXGMS UPS + EATON 93PM
// version:v3
func NewPxgmsUps(data map[string]interface{}) (ups *PxgmsUps, err error) { // nolint: gocyclo

	// Parameter check against data["model"]
	model := data["model"].(string)
	log.Debugf("model is: [%s]", model)
	if !strings.HasPrefix(model, "PXGMS UPS") {
		return nil, fmt.Errorf("only PXGMS UPS models are currently supported")
	}

	// Create the SNMP DeviceConfig,
	snmpDeviceConfig, err := core.GetDeviceConfig(data)
	if err != nil {
		log.WithError(err).Error("[snmp] failed to load device config")
		return nil, err
	}
	log.WithField("config", snmpDeviceConfig).Info("[snmp] loaded device config")

	if err := snmpDeviceConfig.CheckPrivacyAndAuthFromData(data); err != nil {
		return nil, err
	}

	// Create SNMP client.
	snmpClient, err := core.NewSnmpClient(snmpDeviceConfig)
	if err != nil {
		log.WithError(err).Error("[snmp] failed to create new SNMP client")
		return nil, err
	}
	log.Debug("[snmp] created new SNMP client")

	// Create SnmpServerBase.
	snmpServerBase, err := core.NewSnmpServerBase(snmpClient, snmpDeviceConfig)
	if err != nil {
		log.WithError(err).Error("[snmp] failed to create SNMP server base")
		return nil, err
	}
	log.Debug("[snmp] created SNMP server base")

	// Create the UpsMib.
	upsMib, err := mibs.NewUpsMib(snmpServerBase)
	if err != nil {
		log.WithError(err).Error("failed to create the UPS MIB")
		return nil, err
	}
	log.Debug("[snmp] created new UPS MIB")

	// Enumerate the mib.
	snmpDevices, err := upsMib.EnumerateDevices(map[string]interface{}{"rack": "site", "board": "ups"})
	if err != nil {
		log.WithError(err).Error("[snmp] failed to enumerate the UPS MIB")
		return nil, err
	}
	log.WithField("devices", len(snmpDevices)).Debug("[snmp] enumerated the UPS MIB")

	// Output enumerated devices.
	for _, dev := range snmpDevices {
		log.WithField("device", dev).Debug("[snmp] enumerated device")
	}

	// Set up the object.
	return &PxgmsUps{
		&SnmpServer{
			SnmpServerBase: snmpServerBase,
			UpsMib:         upsMib,
			DeviceConfigs:  snmpDevices,
		},
	}, nil
}
