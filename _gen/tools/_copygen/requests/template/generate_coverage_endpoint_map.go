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

	content.WriteString("package main" + "\n")
	content.WriteString("var (\n")
	content.WriteString("// endpoints represents a dependency graph of endpoints (map[dependent][]dependencies).\n")
	content.WriteString("endpoints = map[string]bool{\n")

	var funcs strings.Builder
	for i := range gen.Functions {
		funcs.WriteString("\"" + gen.Functions[i].Name + "\": true,\n")
	}

	content.WriteString(funcs.String())
	content.WriteString("}\n") // map
	content.WriteString(")\n") // var

	return content.String(), nil
}
