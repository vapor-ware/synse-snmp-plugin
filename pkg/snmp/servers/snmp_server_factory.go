package servers

import (
	"fmt"
	"strings"
	log "github.com/sirupsen/logrus"
)

// CreateSnmpServer creates a SnmpServer from the configuration data model string.
func CreateSnmpServer(data map[string]interface{}) (server *SnmpServer, err error) {
	model, ok := data["model"].(string)
	if !ok {
		err = fmt.Errorf("No snmp server model")
		return
	}

	// The model string comes from the configuration.
	// The string itself is a prefix of what is returned from snmpget on OID .1.3.6.1.2.1.33.1.1.2.0
	if strings.HasPrefix(model, "PXGMS UPS") {
		var pxgmups *PxgmsUps
		pxgmups, err = NewPxgmsUps(data)
		if err == nil {
			server = pxgmups.SnmpServer
			return
		}
		return nil, err
	}

	if strings.HasPrefix(model, "Galaxy VM") {
		var galaxyups *GalaxyUps
		galaxyups, err = NewGalaxyUps(data)
		if err == nil {
			server = galaxyups.SnmpServer
			return
		}
		return nil, err
	}

	if model == "SU10000RT3UPM" {
		var trippliteups *TrippliteUps
		trippliteups, err = NewTrippliteUps(data)
		log.Debugf("Error is: [%s]", err)
		if err == nil {
			server = trippliteups.SnmpServer
			return
		}
		return nil, err
	}

	err = fmt.Errorf("Unkown snmp server model: %v", model)
	return
}
