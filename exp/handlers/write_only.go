package handlers

import "github.com/vapor-ware/synse-sdk/sdk"

var WriteOnly = sdk.DeviceHandler{
	Name: "write-only",
	Write: func(device *sdk.Device, data *sdk.WriteData) error {
		return nil
	},
}
