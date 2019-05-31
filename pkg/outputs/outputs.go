package outputs

import "github.com/vapor-ware/synse-sdk/sdk/output"

var (
	// WattsPower describes readings with power (watts) outputs.
	WattsPower = output.Output{
		Name:      "watts",
		Type:      "power",
		Precision: 3,
		Unit: &output.Unit{
			Name:   "watts",
			Symbol: "W",
		},
	}

	// VAPower describes readings with power (volt-ampere) outputs.
	VAPower = output.Output{
		Name:      "volt-ampere",
		Type:      "power",
		Precision: 3,
		Unit: &output.Unit{
			Name:   "volt-ampere",
			Symbol: "VA",
		},
	}

	// Identity describes readings with identity outputs.
	Identity = output.Output{
		Name: "identity",
	}
)
