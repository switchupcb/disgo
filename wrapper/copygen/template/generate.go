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
	fn.WriteString(generateReturn(function))
	return fn.String()
}

////////////////////////////////////////////////////////////////////////////////
// Signature
////////////////////////////////////////////////////////////////////////////////

// generateComment generates a function comment.
func generateComment(function *models.Function) string {
	return "// Send sends a " + function.From[0].Field.FullDefinitionWithoutPointer() + " to Discord and returns a " + function.To[0].Field.FullDefinitionWithoutPointer() + "."
}

// generateSignature generates a function's signature.
func generateSignature(function *models.Function) string {
	return "func (r " + function.From[0].Field.FullDefinition() + ") Send(bot *Client) (" + generateResultParameters(function) + ") {"
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

////////////////////////////////////////////////////////////////////////////////
// Body
////////////////////////////////////////////////////////////////////////////////

// generateBody generates the body of a function.
func generateBody(function *models.Function) string {
	request := function.From[0].Field.FullDefinitionWithoutPointer()
	response := function.To[0].Field.FullDefinition()

	var body strings.Builder
	body.WriteString("var result " + response + "\n")
	body.WriteString("body, err := json.Marshal(r)\n")
	body.WriteString("if err != nil {\n")
	body.WriteString(generateMarshalErrReturn(function, request))
	body.WriteString("}\n")
	body.WriteString("\n")
	body.WriteString("err = SendRequest(result, bot.client, TODO, " + generateEndpointCall(function.From[0].Field) + ", body)\n")
	body.WriteString("if err != nil {\n")
	body.WriteString(generateSendRequestErrReturn(function, request))
	body.WriteString("}\n")

	return body.String()
}

// generateEndpointCall generates the endpoint function call (parameter) for a SendRequestJSON call.
func generateEndpointCall(request *models.Field) string {
	var parameters strings.Builder

	tagCount := 0
	for _, subfield := range request.Fields {

		// subfields with no tags are endpoints.
		if len(subfield.Tags) == 0 {
			if tagCount != 0 {
				parameters.WriteString(", ")
			}

			if subfield.Name == "ApplicationID" {
				parameters.WriteString("bot.ApplicationID")
			} else {
				parameters.WriteString("r." + subfield.Name)
			}

			tagCount++
		}
	}

	return "Endpoint" + request.Definition[1:] + "(" + parameters.String() + ")"
}

////////////////////////////////////////////////////////////////////////////////
// Return
////////////////////////////////////////////////////////////////////////////////

// generateMarshalErrReturn generates a return statement for the function.
func generateMarshalErrReturn(function *models.Function, request string) string {
	errorf := "fmt.Errorf(ErrSendMarshal" + ",\"" + request + "\"" + ", err)"
	switch len(function.To) {
	case 1:
		return "return " + errorf + "\n"
	case 2:
		return "return nil, " + errorf + "\n"
	default:
		return "return nil, " + errorf + "\n"
	}
}

// generateSendRequestErrReturn generates a return statement for the function.
func generateSendRequestErrReturn(function *models.Function, request string) string {
	errorf := "fmt.Errorf(ErrSendRequest" + ",\"" + request + "\"" + ", err)"
	switch len(function.To) {
	case 1:
		return "return " + errorf + "\n"
	case 2:
		return "return nil, " + errorf + "\n"
	default:
		return "return nil, " + errorf + "\n"
	}
}

// generateReturn generates a return statement for the function.
func generateReturn(function *models.Function) string {
	switch len(function.To) {
	case 1:
		return "\nreturn nil\n}\n"
	case 2:
		return "\nreturn result, nil\n}\n"
	default:
		return "\nreturn result, nil\n}\n"
	}
}
