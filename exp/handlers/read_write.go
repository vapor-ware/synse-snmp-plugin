package handlers

import (
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
)

// ReadWrite
var ReadWrite = sdk.DeviceHandler{
	Name: "read-write",
	Read: func(device *sdk.Device) ([]*output.Reading, error) {
		return []*output.Reading{}, nil
	},
	Write: func(device *sdk.Device, data *sdk.WriteData) error {
		return nil
	},
}
