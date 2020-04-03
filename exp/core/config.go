package core

import (
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/utils"
)

type SnmpTargetConfiguration struct {
	MIB       string
	Version   string
	Agent     string
	Community string
	Timeout   time.Duration
	Retries   int
	Security  *SnmpV3Security
}

type SnmpV3Security struct {
	Level          string
	Context        string
	Username       string
	Authentication *SnmpV3SecurityAuthentication
	Privacy        *SnmpV3SecurityPrivacy
}

type SnmpV3SecurityAuthentication struct {
	Protocol   string
	Passphrase string
}

type SnmpV3SecurityPrivacy struct {
	Protocol   string
	Passphrase string
}

func LoadTargetConfiguration(raw map[string]interface{}) (*SnmpTargetConfiguration, error) {
	var cfg SnmpTargetConfiguration

	if err := mapstructure.Decode(raw, &cfg); err != nil {
		log.WithFields(log.Fields{
			"data": utils.RedactPasswords(raw),
		}).Error("[snmp] failed decoding SNMP target configuration into struct")
		return nil, err
	}

	// Set defaults
	if cfg.Retries == 0 {
		cfg.Retries = 1
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 3 * time.Second
	}

	return &cfg, nil
}
