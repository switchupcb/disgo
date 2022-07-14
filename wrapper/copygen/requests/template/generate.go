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
	return "// Send sends a " + function.From[0].Field.FullDefinitionWithoutPointer() + " request to Discord and returns a " + function.To[0].Field.FullDefinitionWithoutPointer() + "."
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
	request := function.From[0].Field
	requestName := request.FullDefinitionWithoutPointer()
	response := function.To[0].Field

	// write the function body.
	var body strings.Builder
	errDecl := ":="

	// marshal the request.
	httpbody := "nil"
	uniquetags := uniqueTags(request)
	if uniquetags["json"] != 0 {
		body.WriteString("body, err " + errDecl + " json.Marshal(r)\n")
		body.WriteString("if err != nil {\n")
		body.WriteString(generateMarshalErrReturn(function, requestName))
		body.WriteString("}\n")
		body.WriteString("\n")
		httpbody = "body"
		errDecl = "="
	}

	// call the endpoint's function.
	endpoint := generateEndpointCall(function.From[0].Field)
	if uniquetags["url"] != 0 {
		body.WriteString("query, err := EndpointQueryString(r)\n")
		body.WriteString("if err != nil {\n")
		body.WriteString(generateQueryStringErrReturn(function, requestName))
		body.WriteString("}\n")
		body.WriteString("\n")
		endpoint = endpoint + "+" + "\"?\"" + "+ query"
		errDecl = "="
	}

	// declare result.
	result := "result"

	switch {
	case response.FullDefinition() == "error":
		result = "nil"
	case response.IsSlice():
		body.WriteString("var result " + response.FullDefinition() + "\n")
	default:
		body.WriteString("result := new(" + response.FullDefinitionWithoutPointer() + ")\n")
	}

	// send the request.
	body.WriteString("err " + errDecl + " SendRequest(bot, " + generateHTTPMethod(function) + ", " +
		endpoint + ", " + generateContentType(uniquetags) + ", " + httpbody + ", " + result + ")\n")
	body.WriteString("if err != nil {\n")
	body.WriteString(generateSendRequestErrReturn(function, requestName))
	body.WriteString("}\n")

	return body.String()
}

// generateHTTPMethod generates the HTTP method type for a SendRequest call.
func generateHTTPMethod(function *models.Function) string {
	http := function.Options.Custom["http"][0]

	var method string
	switch http {
	case "GET":
		method = "Get"
	case "POST":
		method = "Post"
	case "PUT":
		method = "Put"
	case "PATCH":
		method = "Patch"
	case "DELETE":
		method = "Delete"
	}

	return "fasthttp.Method" + method
}

// generateEndpointCall generates the endpoint function call (parameter) for a SendRequest call.
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

// generateContentType generates the content type for a SendRequest call.
func generateContentType(tags map[string]int) string {
	switch {
	case tags["dasgo"] != 0:
		return "contentTypeMulti"
	case tags["json"] != 0:
		return "contentTypeJSON"
	case tags["url"] != 0:
		return "contentTypeURL"
	default:
		return "nil"
	}
}

// uniqueTags determines the unique tags of a request.
func uniqueTags(request *models.Field) map[string]int {
	uniquetags := make(map[string]int)

	for _, subfield := range request.Fields {
		for k := range subfield.Tags {
			uniquetags[k]++
		}
	}

	return uniquetags
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

// generateQueryStringErrReturn generates a return statement for the function.
func generateQueryStringErrReturn(function *models.Function, request string) string {
	errorf := "fmt.Errorf(ErrQueryString" + ",\"" + request + "\"" + ", err)"
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
