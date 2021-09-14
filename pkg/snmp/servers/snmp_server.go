package servers

import (
	"github.com/vapor-ware/synse-sdk/sdk/config"
	"github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/core"
	mibs "github.com/vapor-ware/synse-snmp-plugin/pkg/snmp/mibs/ups_mib"
)

// SnmpServer is a base class for all SnmpServers.
// This is meant to represnt any device serving SNMP.
// Currently only the parts of the UPS-MIB are supported.
// TODO: Replace UpsMib with []*core.SnmpMib so each snmp server supports whatever mibs it supports.
type SnmpServer struct {
	*core.SnmpServerBase                       // base class.
	UpsMib               *mibs.UpsMib          // UpsMib is the only Mib currently supported. In the future this will be a slice or map of supported Mibs.
	DeviceConfigs        []*config.DeviceProto // Enumerated device configs.
}
