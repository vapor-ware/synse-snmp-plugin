package servers

import (
	"fmt"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/core"
	"github.com/vapor-ware/synse-snmp-plugin/snmp/mibs/ups_mib"
)

// PxgmsUps represents the PXGMS UPS + EATON 93PM SNMP Server.
type PxgmsUps struct {
	*core.SnmpServerBase                     // base class.
	UpsMib               *mibs.UpsMib        // Supported Mibs.
	DeviceConfigs        []*sdk.DeviceConfig // Enumerated device configs.
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
	logger.Debugf("model is: [%s]", model)
	if !strings.HasPrefix(model, "PXGMS UPS") {
		return nil, fmt.Errorf("only PXGMS UPS models are currently supported")
	}

	// Create the SNMP DeviceConfig,
	snmpDeviceConfig, err := core.GetDeviceConfig(data)
	if err != nil {
		return nil, err
	}
	fmt.Printf("snmpDeviceConfig: %+v\n", snmpDeviceConfig)

	// Create SNMP client.
	snmpClient, err := core.NewSnmpClient(snmpDeviceConfig)
	if err != nil {
		return nil, err
	}
	fmt.Printf("snmpClient: %+v\n", snmpClient)

	// Create SnmpServerBase.
	snmpServerBase, err := core.NewSnmpServerBase(snmpClient, snmpDeviceConfig)
	if err != nil {
		return nil, err
	}
	fmt.Printf("snmpServerBase: %+v\n", snmpServerBase)

	// Create the UpsMib.
	upsMib, err := mibs.NewUpsMib(snmpServerBase)
	if err != nil {
		return nil, err
	}
	fmt.Printf("upsMib: %+v\n", upsMib)

	// Enumerate the mib.
	snmpDevices, err := upsMib.EnumerateDevices(
		map[string]interface{}{"rack": "site", "board": "ups"})
	if err != nil {
		return nil, err
	}

	// Output enumerated devices.
	for i := 0; i < len(snmpDevices); i++ {
		logger.Debugf("snmpDevice[%d]: %+v\n", i, snmpDevices[i])
	}

	// Set up the object.
	return &PxgmsUps{
		SnmpServerBase: snmpServerBase,
		UpsMib:         upsMib,
		DeviceConfigs:  snmpDevices,
	}, nil
}
