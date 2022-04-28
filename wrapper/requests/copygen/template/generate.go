// DO NOT CHANGE PACKAGE

// Package template provides a template used by copygen to generate custom code.
package template

import (
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

// requestDefinition gets the definition for a request.
func requestDefinition(field *models.Field) string {
	return field.Pointer + field.Definition
}

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
	return "// Send sends a " + requestDefinition(function.From[0].Field) + " to Discord and returns a " +
		function.To[0].Field.FullNameWithoutPointer("") + "."
}

// generateSignature generates a function's signature.
func generateSignature(function *models.Function) string {
	return "func (r " + requestDefinition(function.From[0].Field) + ") Send(bot *Client) (" + generateResultParameters(function) + ") {"
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
	request := requestDefinition(function.From[0].Field)
	response := function.To[0].Field.FullName()

	var body strings.Builder
	body.WriteString("var result " + response + "\n")
	body.WriteString("body, err := json.Marshal(r)\n")
	body.WriteString("if err != nil {\n")
	body.WriteString(generateMarshalErrReturn(function, request))
	body.WriteString("}\n")
	body.WriteString("\n")
	body.WriteString("err = http.SendRequestJSON(bot.client, bot.ctx, http.POST, " + generateEndpointCall(function.From[0].Field) + ", body)\n")
	body.WriteString("if err != nil {\n")
	body.WriteString(generateSendRequestErrReturn(function, request))
	body.WriteString("}\n")
	body.WriteString("\n")
	body.WriteString("err = ParseResponseJSON(bot.ctx, result)\n")
	body.WriteString("if err != nil {\n")
	body.WriteString(generateParseResponseErrReturn(function, request))
	body.WriteString("}\n")
	return body.String()
}

const defaultEndpointParameter = "strconv.FormatUint(uint64(bot.ApplicationID), 10)"

// generateEndpointCall generates the endpoint function call (parameter) for a SendRequestJSON call.
func generateEndpointCall(request *models.Field) string {
	parameters := defaultEndpointParameter

	for _, field := range request.Fields {
		if len(field.Tags) == 0 {
			parameters += ", strconv.FormatUint(uint64(r." + field.Name + "), 10)"
		}
	}

	return "Endpoint" + request.Definition + "(" + parameters + ")"
}

////////////////////////////////////////////////////////////////////////////////
// Return
////////////////////////////////////////////////////////////////////////////////

// generateMarshalErrReturn generates a return statement for the function.
func generateMarshalErrReturn(function *models.Function, request string) string {
	switch len(function.To) {
	case 1:
		return "return fmt.Errorf(\"an error occurred while marshalling a " + request + ": \\n%w\", err)\n"
	case 2:
		return "return nil, fmt.Errorf(\"an error occurred while marshalling a " + request + ": \\n%w\", err)\n"
	default:
		return "return nil, fmt.Errorf(\"an error occurred while marshalling a " + request + ": \\n%w\", err)\n"
	}
}

// generateSendRequestErrReturn generates a return statement for the function.
func generateSendRequestErrReturn(function *models.Function, request string) string {
	switch len(function.To) {
	case 1:
		return "return fmt.Errorf(\"an error occurred while sending a " + request + ": \\n%w\", err)\n"
	case 2:
		return "return nil, fmt.Errorf(\"an error occurred while sending a " + request + ": \\n%w\", err)\n"
	default:
		return "return nil, fmt.Errorf(\"an error occurred while sending a " + request + ": \\n%w\", err)\n"
	}
}

// generateSendRequestErrReturn generates a return statement for the function.
func generateParseResponseErrReturn(function *models.Function, request string) string {
	switch len(function.To) {
	case 1:
		return "return fmt.Errorf(\"an error occurred while parsing the response of a " + request + ": \\n%w\", err)\n"
	case 2:
		return "return nil, fmt.Errorf(\"an error occurred while parsing the response of a " + request + ": \\n%w\", err)\n"
	default:
		return "return nil, fmt.Errorf(\"an error occurred while parsing the response of a " + request + ": \\n%w\", err)\n"
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
