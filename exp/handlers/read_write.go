package handlers

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

var ReadWrite = sdk.DeviceHandler{
	Name: "read-write",
	Read: func(device *sdk.Device) ([]*output.Reading, error) {
		return nil, nil
	},
	Write: func(device *sdk.Device, data *sdk.WriteData) error {
		return nil
	},
}
