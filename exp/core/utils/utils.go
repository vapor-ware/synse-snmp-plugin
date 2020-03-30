package utils

import (
	"github.com/pkg/errors"
	"github.com/soniah/gosnmp"
	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidAuthProtocol = errors.New("unsupported SNMP auth protocol specified")
	ErrInvalidPrivProtocol = errors.New("unsupported SNMP privacy protocol specified")
	ErrInvalidSNMPVersion = errors.New("invalid SNMP version specified")
)

// GetSNMPVersion gets the version of SNMP corresponding to the given string.
func GetSNMPVersion(s string) (gosnmp.SnmpVersion, error) {
	switch s {
	case "v1":
		return gosnmp.Version1, nil
	case "v2", "v2c":
		return  gosnmp.Version2c, nil
	case "v3":
		return gosnmp.Version3, nil
	default:
		log.WithFields(log.Fields{
			"version": s,
		}).Error("[snmp] invalid SNMP version specified")
		return -1, ErrInvalidSNMPVersion
	}
}

// GetPrivProtocol gets the privacy protocol constant corresponding to the
// given string.
func GetPrivProtocol(s string) (gosnmp.SnmpV3PrivProtocol, error) {
	switch s {
	case "AES":
		return gosnmp.AES, nil
	case "DES":
		return gosnmp.DES, nil
	case "none":
		return gosnmp.NoPriv,  nil
	default:
		log.WithFields(log.Fields{
			"privacy": s,
		}).Error("[snmp] invalid SNMP privacy protocol specified")
		return 0, ErrInvalidPrivProtocol
	}
}

// GetAuthProtocol gets the authentication protocol constant corresponding to
// the given string.
func GetAuthProtocol(s string) (gosnmp.SnmpV3AuthProtocol, error) {
	switch s {
	case "MD5":
		return gosnmp.MD5, nil
	case "SHA":
		return gosnmp.SHA, nil
	case "none":
		return gosnmp.NoAuth, nil
	default:
		log.WithFields(log.Fields{
			"authentication": s,
		}).Error("[snmp] invalid SNMP authentication protocol specified")
		return 0, ErrInvalidAuthProtocol
	}
}