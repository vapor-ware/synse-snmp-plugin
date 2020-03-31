package handlers

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

var ReadOnly = sdk.DeviceHandler{
	Name: "read-only",
	Read: func(device *sdk.Device) ([]*output.Reading, error) {
		return nil, nil
	},
}
