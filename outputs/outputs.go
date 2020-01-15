package outputs

import "github.com/vapor-ware/synse-sdk/sdk"

var (
	// Voltage describes readings with voltage outputs.
	Voltage = sdk.OutputType{
		Name:      "voltage",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "volts",
			Symbol: "V",
		},
	}

	// Temperature describes readings with temperature (celsius) outputs.
	Temperature = sdk.OutputType{
		Name:      "temperature",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "degrees celsius",
			Symbol: "C",
		},
	}

	// StatusInt describes readings with status outputs where status is an
	// integer. Status is a very generic output and avoids a plethora of more
	// specific output types.
	StatusInt = sdk.OutputType{
		Name: "status-int",
	}

	// StatusUint describes readings with status outputs where status is an
	// integer. Status is a very generic output and avoids a plethora of more
	// specific output types.
	StatusUint = sdk.OutputType{
		Name: "status-uint",
	}

	// StatusString describes readings with status outputs where status is a
	// string, for example an enumeration.
	StatusString = sdk.OutputType{
		Name: "status-string",
	}

	// WattsPower describes readings with power (watts) outputs.
	WattsPower = sdk.OutputType{
		Name:      "watts.power",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "watts",
			Symbol: "W",
		},
	}

	// VAPower describes readings with power (volt-ampere) outputs.
	VAPower = sdk.OutputType{
		Name:      "va.power",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "volt-ampere",
			Symbol: "VA",
		},
	}

	// Identity describes readings with identity outputs.
	Identity = sdk.OutputType{
		Name: "identity",
	}

	// Frequency describes readings with frequency (Hz) outputs.
	Frequency = sdk.OutputType{
		Name:      "frequency",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "hertz",
			Symbol: "Hz",
		},
	}

	// Current describes readings with current (Amps) outputs.
	Current = sdk.OutputType{
		Name:      "current",
		Precision: 3,
		Unit: sdk.Unit{
			Name:   "amps",
			Symbol: "A",
		},
	}
)
