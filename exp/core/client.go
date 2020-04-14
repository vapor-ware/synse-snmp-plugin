package core

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
)

// Errors relating to SNMP plugin client creation and usage.
var (
	ErrNonV3SecurityParams = errors.New("cannot define security parameters for SNMP versions other than v3")
)

var clientCache = map[string]*Client{}

// Client is a wrapper around a GoSNMP struct which adds some utility
// functions around it. Notably, it enables lazy connecting to the client,
// so the SNMP agent does not need to be reachable at plugin startup.
type Client struct {
	*gosnmp.GoSNMP

	isConnected bool
}

// GetOid gets the value for a specified OID.
func (c *Client) GetOid(oid string) (*gosnmp.SnmpPDU, error) {
	if !c.isConnected {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	result, err := c.Get([]string{oid})
	if err != nil {
		return nil, err
	}

	// Since we are currently only reading one OID, the result value will be
	// the first and only returned variable in the response.
	data := result.Variables[0]

	return &data, nil
}

// GetClient returns the Client for the specified target.
//
// If there is no client for the specified target cached, nil is returned.
func GetClient(target string) *Client {
	return clientCache[target]
}

// CacheClient caches an SNMP client. The key a client is cached against is generated
// from the client configuration values.
//
// If a client is already cached with a given key, it will be overwritten with the new
// client instance.
func CacheClient(client *Client) {
	key := fmt.Sprintf("%s://%s:%d", client.Transport, client.Target, client.Port)

	if _, exists := clientCache[key]; exists {
		log.WithFields(log.Fields{
			"target": key,
		}).Warn("[snmp] overwriting previously cached client")
	}
	clientCache[key] = client
}

// NewClient creates a new instance of an SNMP Client for the given SNMP target
// configuration. The SNMP target configuration is defined in the dynamic configuration
// block for the plugin.
func NewClient(cfg *SnmpTargetConfiguration) (*Client, error) {
	// Verify that the configured version is valid.
	version, err := GetSNMPVersion(cfg.Version)
	if err != nil {
		return nil, err
	}

	// If configured to use SNMP v3, verify the security parameters.
	var securityModel gosnmp.SnmpV3SecurityModel
	var msgFlags gosnmp.SnmpV3MsgFlags
	var securityParams *gosnmp.UsmSecurityParameters
	var contextName string

	if version == gosnmp.Version3 {
		// If there are no security parameters defined, assume NoAuthNoPriv.
		if cfg.Security == nil {
			log.Info("[snmp] no security parameters defined, falling back to NoAuthNoPriv")

		} else {
			contextName = cfg.Security.Context

			msgFlags, err = GetSecurityFlags(cfg.Security.Level)
			if err != nil {
				return nil, err
			}
			if msgFlags != gosnmp.NoAuthNoPriv {
				securityModel = gosnmp.UserSecurityModel
			}

			var (
				authPass  = ""
				authProto = gosnmp.NoAuth
				privPass  = ""
				privProto = gosnmp.NoPriv
			)
			if cfg.Security.Authentication != nil {
				authPass = cfg.Security.Authentication.Passphrase
				authProto, err = GetAuthProtocol(cfg.Security.Authentication.Protocol)
				if err != nil {
					return nil, err
				}
			}
			if cfg.Security.Privacy != nil {
				privPass = cfg.Security.Privacy.Passphrase
				privProto, err = GetPrivProtocol(cfg.Security.Privacy.Protocol)
				if err != nil {
					return nil, err
				}
			}
			securityParams = &gosnmp.UsmSecurityParameters{
				UserName:                 cfg.Security.Username,
				AuthenticationPassphrase: authPass,
				AuthenticationProtocol:   authProto,
				PrivacyPassphrase:        privPass,
				PrivacyProtocol:          privProto,
			}
		}

	} else {
		// If not using SNMP v3, no security parameters should be defined. If they
		// are, it should be considered a misconfiguration.
		if cfg.Security != nil {
			log.WithFields(log.Fields{
				"version": cfg.Version,
			}).Error("[snmp] security parameters unsupported for configured SNMP version")
			return nil, ErrNonV3SecurityParams
		}
	}

	agent := cfg.Agent
	if !strings.Contains(agent, "://") {
		agent = "udp://" + agent
	}

	u, err := url.Parse(agent)
	if err != nil {
		return nil, err
	}

	var transport string
	switch u.Scheme {
	case "tcp":
		transport = "tcp"
	case "", "udp":
		transport = "udp"
	default:
		return nil, fmt.Errorf("unsupported transport scheme: %s", u.Scheme)
	}

	portStr := u.Port()
	if portStr == "" {
		portStr = "161"
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"host": u.Hostname(),
		"port": port,
		"transport": transport,
	}).Debug("[snmp] parsed client agent config")

	c := &Client{
		GoSNMP: &gosnmp.GoSNMP{
			Version:            version,
			Target:             u.Hostname(),
			Port:               uint16(port),
			Transport:          transport,
			Timeout:            cfg.Timeout,
			Retries:            cfg.Retries,
			Community:          cfg.Community,
			MsgFlags:           msgFlags,
			SecurityModel:      securityModel,
			SecurityParameters: securityParams,
			ContextName:        contextName,
			MaxOids:            gosnmp.MaxOids,
		},
	}

	log.WithFields(log.Fields{
		"cfg": *c.GoSNMP,
	}).Debug("[snmp] created new client")
	return c, nil
}
