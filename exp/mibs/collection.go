package mibs

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// pluginMibs holds the collection of MIBs which have been registered with the
// SNMP plugin base.
var pluginMibs = map[string]*MIB{}

// Errors for SNMP base plugin MIB operations.
var (
	ErrMibExists = errors.New("MIB already registered")
	ErrNilMib    = errors.New("MIB cannot be nil")
)

// Register MIBs defined by a plugin implementation with the SNMP base plugin.
//
// If a MIB with the same name is already registered, an error is returned.
func Register(mibs ...*MIB) error {
	for _, mib := range mibs {
		if mib == nil {
			log.Error("[snmp] cannot register nil MIB")
			return ErrNilMib
		}

		// Check if the MIB being registered conflicts with a MIB that has already
		// been registered.
		if _, exists := pluginMibs[mib.Name]; exists {
			log.WithFields(log.Fields{
				"mib": mib.Name,
			}).Error("[snmp] unable to register MIB; conflicts with existing MIB")
			return ErrMibExists
		}

		pluginMibs[mib.Name] = mib
	}
	return nil
}

// Get a registered MIB with the given name.
//
// If there is no MIB registered with the provided name, nil is returned.
func Get(name string) *MIB {
	return pluginMibs[name]
}

// GetAll returns all of the MIBs which have been registered with the SNMP base
// plugin. The order in which the MIBs are returned is not guaranteed. If no MIBs
// have been registered, an empty slice is returned.
func GetAll() []*MIB {
	mibs := make([]*MIB, len(pluginMibs))
	i := 0
	for _, mib := range pluginMibs {
		mibs[i] = mib
		i++
	}
	return mibs
}

// Clear removes all data from the global MIB collection.
//
// Generally, this should not be used by a plugin implementation, however
// it is useful for testing.
func Clear() {
	pluginMibs = map[string]*MIB{}
}
