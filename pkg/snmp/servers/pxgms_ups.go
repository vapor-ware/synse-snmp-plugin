package servers

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
	mibs "github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/mibs/ups_mib"
)

// PxgmsUps represents the PXGMS UPS + EATON 93PM SNMP Server.
type PxgmsUps struct {
	*core.SnmpServerBase                       // base class.
	UpsMib               *mibs.UpsMib          // Supported Mibs.
	DeviceConfigs        []*config.DeviceProto // Enumerated device configs.
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

	// FIXME (etd): Sorta a hack just to get things moving, but adding in a check against
	// the model here. There could probably be something at a higher level that checks this
	// and initializes the right stuff based on the specified model.
	// TODO: File a ticket on this. Checking against the model is ruining SNMP. Any MIB can
	// now only support one and only one model.
	// We intend to be able to share SNMP MIBs across models and this won't work at all.
	// (mhink): The SNMP server level factory is NYI due to other oblications. This will allow sharing MIBs.
	// Ticket is here: https://github.com/vapor-ware/synse-snmp-plugin/issues/10
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
		SnmpServerBase: snmpServerBase,
		UpsMib:         upsMib,
		DeviceConfigs:  snmpDevices,
	}, nil
}
