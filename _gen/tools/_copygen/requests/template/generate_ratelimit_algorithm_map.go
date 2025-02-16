// DO NOT CHANGE PACKAGE

// Package template provides a template used by copygen to generate custom code.
package template

import (
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

// Generate generates code.
// GENERATOR FUNCTION.
// EDITABLE.
// DO NOT REMOVE.
func Generate(gen *models.Generator) (string, error) {
	var content strings.Builder

	content.WriteString(string(gen.Keep) + "\n")
	content.WriteString("var (\n")
	content.WriteString("// RateLimitHashFuncs represents a map of routes to respective rate limit algorithms.\n")
	content.WriteString("// \n")
	content.WriteString("// Used to determine the hashing function for routes during runtime (map[routeID]algorithm).\n")
	content.WriteString("RateLimitHashFuncs = map[uint8]hash{\n")

	var funcs strings.Builder
	for i := range gen.Functions {
		funcs.WriteString("RouteIDs[\"" + gen.Functions[i].Name + "\"]: HashPerRoute,\n")
	}

	content.WriteString(funcs.String())
	content.WriteString("}\n") // map
	content.WriteString(")\n") // var

	return content.String(), nil
}
