package outputs

import "github.com/vapor-ware/synse-sdk/sdk"

var (
	Voltage = sdk.OutputType{
		Name:      "voltage",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "volts",
			Symbol: "V",
		},
	}

	Temperature = sdk.OutputType{
		Name:      "temperature",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "degrees celsius",
			Symbol: "C",
		},
	}

	Status = sdk.OutputType{
		Name: "status",
	}

	WattsPower = sdk.OutputType{
		Name:      "watts.power",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "watts",
			Symbol: "W",
		},
	}

	VAPower = sdk.OutputType{
		Name:      "va.power",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "volt-ampere",
			Symbol: "VA",
		},
	}

	Identity = sdk.OutputType{
		Name: "identity",
	}

	Frequency = sdk.OutputType{
		Name:      "frequency",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "hertz",
			Symbol: "Hz",
		},
	}

	Current = sdk.OutputType{
		Name:      "current",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "amps",
			Symbol: "A",
		},
	}
)
