package core

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
)

// CheckPrivacyAndAuthFromData determines a MsgFlag based on parsed data
func (d *DeviceConfig) CheckPrivacyAndAuthFromData(data map[string]interface{}) error {
	auth, ok := data["authenticationProtocol"].(string)
	if !ok {
		return fmt.Errorf("cannot find authenticationProtocol")
	}
	priv, ok := data["privacyProtocol"].(string)
	if !ok {
		return fmt.Errorf("cannot find privacyProtocol")
	}

	switch {
	case auth == "None" && priv == "None":
		d.MsgFlag = gosnmp.NoAuthNoPriv
	case auth != "None" && priv == "None":
		d.MsgFlag = gosnmp.AuthNoPriv
	default:
		d.MsgFlag = gosnmp.AuthPriv
	}

	return nil
}
