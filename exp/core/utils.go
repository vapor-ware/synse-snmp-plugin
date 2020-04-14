package core

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// Errors for utility functions, largely around parsing SNMP configurations.
var (
	ErrInvalidAuthProtocol = errors.New("unsupported SNMP auth protocol specified")
	ErrInvalidPrivProtocol = errors.New("unsupported SNMP privacy protocol specified")
	ErrInvalidSNMPVersion  = errors.New("invalid SNMP version specified")
	ErrInvalidMessageFlag  = errors.New("invalid security message flag specified")
)

// TagOrPanic is a utility function which creates a new SDK Tag or it panics.
//
// This is useful when you know the format of the tag is correct ahead of time
// and you which skip the error check.
func TagOrPanic(tag string) *sdk.Tag {
	t, err := sdk.NewTag(tag)
	if err != nil {
		panic(err)
	}
	return t
}

// GetSNMPVersion gets the version of SNMP corresponding to the given string.
func GetSNMPVersion(s string) (gosnmp.SnmpVersion, error) {
	switch strings.ToLower(s) {
	case "v1":
		return gosnmp.Version1, nil
	case "v2", "v2c":
		return gosnmp.Version2c, nil
	case "v3":
		return gosnmp.Version3, nil
	default:
		log.WithFields(log.Fields{
			"version": s,
		}).Error("[snmp] invalid SNMP version specified")
		return 0, ErrInvalidSNMPVersion
	}
}

// GetPrivProtocol gets the SNMP v3 privacy protocol constant corresponding to the
// given string.
func GetPrivProtocol(s string) (gosnmp.SnmpV3PrivProtocol, error) {
	switch strings.ToLower(s) {
	case "aes":
		return gosnmp.AES, nil
	case "des":
		return gosnmp.DES, nil
	case "none":
		return gosnmp.NoPriv, nil
	default:
		log.WithFields(log.Fields{
			"privacy": s,
		}).Error("[snmp] invalid SNMP privacy protocol specified")
		return 0, ErrInvalidPrivProtocol
	}
}

// GetAuthProtocol gets the SNMP v3 authentication protocol constant corresponding to
// the given string.
func GetAuthProtocol(s string) (gosnmp.SnmpV3AuthProtocol, error) {
	switch strings.ToLower(s) {
	case "md5":
		return gosnmp.MD5, nil
	case "sha":
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

// GetSecurityFlags gets the SNMP v3 security message flags constant corresponding to
// the given string.
func GetSecurityFlags(s string) (gosnmp.SnmpV3MsgFlags, error) {
	switch strings.ToLower(s) {
	case "noauthnopriv":
		return gosnmp.NoAuthNoPriv, nil
	case "authnopriv":
		return gosnmp.AuthNoPriv, nil
	case "authpriv":
		return gosnmp.AuthPriv, nil
	case "reportable":
		return gosnmp.Reportable, nil
	default:
		log.WithFields(log.Fields{
			"flag": s,
		}).Error("[snmp] invalid SNMP message flag specified")
		return 0, ErrInvalidMessageFlag
	}
}
