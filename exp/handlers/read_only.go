package handlers

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-plugin/exp/core"
)

// ReadOnly
var ReadOnly = sdk.DeviceHandler{
	Name: "read-only",
	Read: func(device *sdk.Device) ([]*output.Reading, error) {
		if device == nil {
			return nil, errors.New("unable to read from nil device")
		}

		agent := device.Data["agent"].(string)
		c := core.GetClient(agent)
		if c == nil {
			return nil, fmt.Errorf("no client cached for device agent %s", agent)
		}

		oid := device.Data["oid"].(string)

		log.WithFields(log.Fields{
			"agent": agent,
			"oid": oid,
		}).Debug("[snmp] reading OID")

		result, err := c.GetOid(oid)
		if err != nil {
			return nil, err
		}

		log.WithFields(log.Fields{
			"value": result.Value,
			"name": result.Name,
			"type": result.Type,
		}).Debug("[snmp] got reading value for OID")

		var value interface{}
		switch result.Type {
		case gosnmp.OctetString:
			ascii, err := TranslatePrintableASCII(result.Value)
			if err != nil {
				return nil, err
			}
			value = ascii
		default:
			value = result.Value
		}

		log.WithFields(log.Fields{
			"value": value,
		}).Debug("[snmp] final value")

		o := output.Get(device.Output)
		if o == nil {
			return nil, fmt.Errorf("unable to format reading: device output not defined")
		}

		return []*output.Reading{
			o.MakeReading(value),
		}, nil
	},
}


// TranslatePrintableASCII translates byte arrays from gosnmp to a printable
// string if possible. If this call fails, the caller should normally just keep
// the raw byte array. This call makes no attempt to support extended (8bit)
// ASCII. We need this function since there is no differentiation between
// strings and byte arrays in the SNMP protocol.
func TranslatePrintableASCII(x interface{}) (string, error) {
	bytes, ok := x.([]uint8)
	if !ok {
		return "", fmt.Errorf("failure converting type: %T, data: %v to byte array", x, x)
	}

	for i := 0; i < len(bytes); i++ {
		if !(bytes[i] < 0x80 && bytes[i] > 0x1f) {
			return "", fmt.Errorf("unable to convert %x byte %x at %d to printable Ascii", bytes, bytes[i], i)
		}
	}
	return string(bytes), nil
}
