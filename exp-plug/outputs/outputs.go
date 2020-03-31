package outputs

import "github.com/vapor-ware/synse-sdk/sdk/output"

var (
	// Identity describes readings with identity outputs.
	Identity = output.Output{
		Name: "identity",
		Type: "identity",
	}
)
