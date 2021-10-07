package servers

import (
	"fmt"
	"strings"
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
		server = pxgmups.SnmpServer
		return
	}

	if strings.HasPrefix(model, "Galaxy VM") {
		var galaxyups *GalaxyUps
		galaxyups, err = NewGalaxyUps(data)
		server = galaxyups.SnmpServer
		return
	}

	if model == "SU10000RT3UPM" {
		var galaxyups *GalaxyUps
		galaxyups, err = NewGalaxyUps(data)
		server = galaxyups.SnmpServer
		return
	}

	err = fmt.Errorf("Unkown snmp server model: %v", model)
	return
}
