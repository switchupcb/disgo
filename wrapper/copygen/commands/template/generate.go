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
	for i := range gen.Functions {
		content.WriteString(Function(&gen.Functions[i]) + "\n")
	}

	return content.String(), nil
}

// Function provides generated code for a function.
func Function(function *models.Function) string {
	var fn strings.Builder
	fn.WriteString(generateComment(function) + "\n")
	fn.WriteString(generateSignature(function) + "\n")
	fn.WriteString(generateBody(function))

	return fn.String()
}

////////////////////////////////////////////////////////////////////////////////
// Signature
////////////////////////////////////////////////////////////////////////////////

// generateComment generates a function comment.
func generateComment(function *models.Function) string {
	return "// Command sends an Opcode " + function.Options.Custom["opcode"][0] + " " + function.Name + " command to the Discord Gateway."
}

// generateSignature generates a function's signature.
func generateSignature(function *models.Function) string {
	return "func (c " + function.From[0].Field.FullDefinition() + ") Command(session *Session) (" + generateResultParameters(function) + ") {"
}

// generateResultParameters generates the result parameters of a function.
func generateResultParameters(function *models.Function) string {
	var parameters strings.Builder
	for i, toType := range function.To {
		if i+1 == len(function.To) {
			parameters.WriteString(" " + toType.Name())
			break
		}

		parameters.WriteString(toType.Name() + ", ")
	}

	return parameters.String()
}

// generateParameters generates the parameters of a function.
func generateParameters(function *models.Function) string {
	var parameters strings.Builder
	for _, toType := range function.To {
		parameters.WriteString(toType.Field.VariableName + " " + toType.Name() + ", ")
	}

	for i, fromType := range function.From {
		if i+1 == len(function.From) {
			parameters.WriteString(fromType.Field.VariableName + " " + fromType.Name())
			break
		}

		parameters.WriteString(fromType.Field.VariableName + " " + fromType.Name() + ", ")
	}

	return parameters.String()
}

// generateBody generates the body of a function.
func generateBody(function *models.Function) string {
	var body strings.Builder
	name := function.From[0].Field.FullDefinitionWithoutPointer()

	// rename GatewayPresenceUpdate to PresenceUpdate.
	if name == "GatewayPresenceUpdate" {
		name = "PresenceUpdate"
	}

	body.WriteString("if err := writeEvent(session, FlagGatewayOpcode" + name + ", " +
		"FlagGatewayCommandName" + name + ", c); err != nil {\n")
	body.WriteString("return err\n")
	body.WriteString("}\n\n")
	body.WriteString("return nil\n")
	body.WriteString("}\n")

	return body.String()
}
