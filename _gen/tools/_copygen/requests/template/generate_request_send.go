// DO NOT CHANGE PACKAGE

// Package template provides a template used by copygen to generate custom code.
package template

import (
	"fmt"
	"hash"
	"hash/fnv"
	"sort"
	"strconv"
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
	content.WriteString("import json \"github.com/goccy/go-json\"\n")

	var funcs strings.Builder
	for i := range gen.Functions {
		funcs.WriteString(Function(&gen.Functions[i]) + "\n")
	}

	// call generateRouteIDs after the routeidMap is populated.
	content.WriteString(generateRouteIDs() + "\n")
	content.WriteString(funcs.String())

	return content.String(), nil
}

////////////////////////////////////////////////////////////////////////////////
// Route IDs
////////////////////////////////////////////////////////////////////////////////

// routeid represents the internal Rate Limit Bucket ID of a Route (endpoint + HTTP Method).
//
// Global Rate Limit reserves 0.
// OAuth reserves 1.
var routeid = 2

// routeidMap represents a map of Routes to RouteIDs (map[string]uint8).
//
// "" resolves to the Global Rate Limit "Route".
var routeidMap = map[string]int{
	"":      0,
	"OAuth": 1,
}

// generateRouteIDs generates the RouteIDs map.
func generateRouteIDs() string {
	var decl strings.Builder
	decl.WriteString("var (\n")
	decl.WriteString("// RouteIDs represents a map of Routes to Route IDs (map[string]uint8).\n")
	decl.WriteString("RouteIDs = map[string]uint8 {\n")

	// sort the map by value.
	keys := make([]string, 0, len(routeidMap))
	for route := range routeidMap {
		keys = append(keys, route)
	}

	sort.Slice(keys, func(i, j int) bool { return routeidMap[keys[i]] < routeidMap[keys[j]] })

	// populate the written map.
	for _, route := range keys {
		decl.WriteString(fmt.Sprintf("\"%s\": %d,\n", route, routeidMap[route]))
	}

	decl.WriteString("}\n")
	decl.WriteString(")")
	return decl.String()
}

////////////////////////////////////////////////////////////////////////////////
// Functions
////////////////////////////////////////////////////////////////////////////////

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
	tags := requestTags(request)

	var body strings.Builder

	// Write logged variables.
	body.WriteString("var err error\n")
	body.WriteString("xid := xid.New().String()\n")
	body.WriteString("routeid, resourceid := " + generateHashCall(request, routeid) + "\n")

	endpoint := generateEndpointCall(request)
	if len(tags["url"]) != 0 {
		body.WriteString("query, err := EndpointQueryString(r)\n")
		body.WriteString("if err != nil {\n")
		body.WriteString(generateQueryStringErrReturn(function, requestName) + "\n")
		body.WriteString("}\n")
		endpoint = endpoint + "+" + "\"?\"" + "+ query"
	}

	body.WriteString("endpoint := " + endpoint + "\n")
	body.WriteString("\n")

	// Write the function body.
	//
	// marshal the request.
	httpbody := "nil"
	if len(tags["json"]) != 0 {
		body.WriteString("body, err := json.Marshal(r)\n")
		body.WriteString("if err != nil {\n")
		body.WriteString(generateMarshalErrReturn(function, requestName) + "\n")
		body.WriteString("}\n")
		body.WriteString("\n")
		httpbody = "body"
	}

	// add file upload support (if applicable).
	contentType := generateContentType(tags)
	if contentType == "varies" {
		contentType = "contentType"

		body.WriteString("contentType := ContentTypeJSON\n")
		body.WriteString("if len(r.Files) != 0 {\n")
		body.WriteString("var multipartErr error\n")
		body.WriteString("if contentType, body, multipartErr = createMultipartForm(body, r.Files...); multipartErr != nil {\n")
		body.WriteString(generateMultipartErrReturn(function, requestName) + "\n")
		body.WriteString("}\n")
		body.WriteString("}\n")
		body.WriteString("\n")
	} else if contentType == "ContentTypeMultipartForm" {
		contentType = "contentType"

		body.WriteString("var contentType []byte\n")
		body.WriteString("var multipartErr error\n")
		body.WriteString("if contentType, body, multipartErr = createMultipartForm(body, &r.File); multipartErr != nil {\n")
		body.WriteString(generateMultipartErrReturn(function, requestName) + "\n")
		body.WriteString("}\n")
		body.WriteString("\n")
	}

	// declare result.
	result := "result"

	switch {
	case response.FullDefinition() == "error":
		result = "nil"

	case response.IsSlice():
		body.WriteString("result := make(" + response.FullDefinition() + ", 0)\n")
		result = "&" + result

	default:
		body.WriteString("result := new(" + response.FullDefinitionWithoutPointer() + ")\n")
	}

	// send the request.
	body.WriteString("err = SendRequest(bot, xid, routeid, resourceid, " +
		generateHTTPMethod(function) + ", endpoint, " +
		contentType + ", " + httpbody + ", " + result + ")\n",
	)
	body.WriteString("if err != nil {\n")
	body.WriteString(generateSendRequestErrReturn(function, requestName) + "\n")
	body.WriteString("}\n")

	// map the route to the route id.
	routeidMap[requestName] = routeid

	// increment route for the next request (if applicable).
	routeid++

	return body.String()
}

var hasher hash.Hash32 = fnv.New32a()

// generateHashCall generates the hashing function call for the Send() function.
func generateHashCall(request *models.Field, routeid int) string {
	var parameters strings.Builder
	routeidstring := strconv.Itoa(routeid)

	parameters.WriteString("\"" + routeidstring + "\"")

	tagCount := 0
	for _, subfield := range request.Fields {

		// subfields without marshalled tags (i.e NOT `-`) are endpoint parameters (i.e `GuildID`).
		if fieldTags(subfield) == 0 {
			if subfield.Name == "ApplicationID" || subfield.Definition != "string" {
				continue
			} else {
				hasher.Write([]byte(subfield.Name))

				identifier := string(hasher.Sum(nil))
				parameters.WriteString(fmt.Sprintf(",\"%x\" + %s", identifier, "r."+subfield.Name))

				hasher.Reset()
			}

			tagCount++
		}
	}

	return "RateLimitHashFuncs[" + routeidstring + "]" + "(" + parameters.String() + ")"
}

// generateHTTPMethod generates the HTTP method type for a SendRequest(..., method) call.
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

// generateEndpointCall generates the endpoint function call for a SendRequest(..., endpoint) call.
//
// NOTE: must be run after the ./_gen clean operation.
func generateEndpointCall(request *models.Field) string {
	var parameters strings.Builder

	tagCount := 0
	for _, subfield := range request.Fields {
		if subfield.Definition != "string" {
			continue
		}

		// subfields without marshalled tags (i.e NOT `-`) are endpoint parameters (i.e `GuildID`).
		if fieldTags(subfield) == 0 {
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

// generateContentType generates the content type for a SendRequest(..., content type) call.
func generateContentType(tags map[string][]string) string {
	switch {
	case len(tags["dasgo"]) != 0:
		if tags["dasgo"][0] == "file" {
			return "ContentTypeMultipartForm"
		}

		return "varies"

	case len(tags["json"]) != 0:
		return "ContentTypeJSON"

	case len(tags["url"]) != 0:
		return "ContentTypeURLQueryString"

	default:
		return "nil"
	}
}

// requestTags determines the tags used in a request struct ('s subfields).
func requestTags(request *models.Field) map[string][]string {
	tags := make(map[string][]string)

	for _, subfield := range request.Fields {
		// check subfields of embedded fields.
		if subfield.Embedded {
			embedTags := requestTags(subfield)

			for tagKey, tagNames := range embedTags {
				tags[tagKey] = append(tags[tagKey], tagNames...)
			}
		}

		// tagKey (i.e `json`, `url`, `dasgo`)
		for tagKey := range subfield.Tags {
			// tagName (i.e `name`, `type`)
			for tagName := range subfield.Tags[tagKey] {
				if tagName != "-" {
					tags[tagKey] = append(tags[tagKey], tagName)
				}
			}
		}
	}

	return tags
}

// fieldTags determines the tags a field uses.
func fieldTags(field *models.Field) int {
	tags := 0

	// tagKey (i.e `json`, `url`, `dasgo`)
	for tagKey := range field.Tags {
		// tagName (i.e `name`, `type`)
		for tagName := range field.Tags[tagKey] {
			if tagName != "-" {
				tags++
			}
		}
	}

	return tags
}

////////////////////////////////////////////////////////////////////////////////
// Return
////////////////////////////////////////////////////////////////////////////////

var requestError = `ErrorRequest{
	ClientID:      bot.ApplicationID,
	CorrelationID: xid,
	RouteID:       routeid,
	ResourceID:    resourceid,
	Endpoint:      %s,
	Err:           %s,
}`

// generateQueryStringErrReturn generates a return statement for the function.
func generateQueryStringErrReturn(function *models.Function, request string) string {
	err := fmt.Sprintf(requestError, "\"\"", "err")
	switch len(function.To) {
	case 1:
		return "return " + err
	case 2:
		return "return nil, " + err
	default:
		return "return nil, " + err
	}
}

// generateMarshalErrReturn generates a return statement for the function.
func generateMarshalErrReturn(function *models.Function, request string) string {
	err := fmt.Sprintf(requestError, "endpoint", "fmt.Errorf(errSendMarshal, err)")
	switch len(function.To) {
	case 1:
		return "return " + err
	case 2:
		return "return nil, " + err
	default:
		return "return nil, " + err
	}
}

// generateMultipartErrReturn generates a return statement for the function.
func generateMultipartErrReturn(function *models.Function, request string) string {
	err := fmt.Sprintf(requestError, "\"\"", "err")
	switch len(function.To) {
	case 1:
		return "return " + err
	case 2:
		return "return nil, " + err
	default:
		return "return nil, " + err
	}
}

// generateSendRequestErrReturn generates a return statement for the function.
func generateSendRequestErrReturn(function *models.Function, request string) string {
	err := fmt.Sprintf(requestError, "endpoint", "err")
	switch len(function.To) {
	case 1:
		return "return " + err
	case 2:
		return "return nil, " + err
	default:
		return "return nil, " + err
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
