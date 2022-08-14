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
	content.WriteString("func IntegrationTestFlow(t *testing.T) {\n\n")
	content.WriteString("\tbot := &Client{\n\t")
	content.WriteString("\tAuthentication: BotToken(os.Getenv(\"TOKEN\")),\n")
	content.WriteString("\tConfig:         DefaultConfig(),\n\t}\n")
	content.WriteString("\t	eg, ctx := errgroup.WithContext(context.Background())\n\n\t")

	for i := range gen.Functions {
		content.WriteString(Function(&gen.Functions[i]) + "\n")
	}

	content.WriteString("\n}")
	return content.String(), nil
}

// Function provides generated code for a function.
func Function(function *models.Function) string {
	var fn strings.Builder

	if strings.Contains(function.Name, "Create") {
		fn.WriteString(generateComment(function) + "\n")
		fn.WriteString(generateRequest(function))
		fn.WriteString(generateErrHandle(function) + "\n")
	} else {
		fn.WriteString(generateErrGroup(function))
		fn.WriteString(generateComment(function) + "\n")
		fn.WriteString(generateRequest(function))
		fn.WriteString(generateErrHandle(function) + "\n})\n")
	}

	return fn.String()
}

// generateComment generates a function comment.
func generateComment(function *models.Function) string {
	var toComment strings.Builder
	for i, toType := range function.To {
		if i+1 == len(function.To) {
			toComment.WriteString(toType.Name())
			break
		}

		toComment.WriteString(toType.Name() + ", ")
	}

	var fromComment strings.Builder
	for i, fromType := range function.From {
		if i+1 == len(function.From) {
			fromComment.WriteString(fromType.Name())
			break
		}

		fromComment.WriteString(fromType.Name() + ", ")
	}

	return "// " + function.Name + " sends a  " + function.Name + "  reqeust to the Discord Gateway."
}

// generateSignature generates a function's signature.
func generateRequest(function *models.Function) string {
	return "request" + function.Name + " := " + function.Name + generateEndpointCall(function.From[0].Field) + "\n"
}

func generateErrHandle(function *models.Function) string {

	if strings.Contains(function.Name, "Create") {
		return "if _, err := request" + function.Name + ".Send(bot); err != nil {\n\treturn\n}\n"
	} else {
		return "if _, err := request" + function.Name + ".Send(bot); err != nil {\n\treturn err\n}\nreturn nil\n"
	}
}

func generateErrGroup(function *models.Function) string {
	return "eg.Go( func() error {\n"
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

	// Assign fields to ToType(s).
	for i, toType := range function.To {
		body.WriteString(generateAssignment(toType))
		if i+1 != len(function.To) {
			body.WriteString("\n")
		}
	}

	return body.String()
}

// generateAssignment generates assignments for a to-type.
func generateAssignment(toType models.Type) string {
	var assign strings.Builder
	assign.WriteString("// " + toType.Name() + " fields\n")

	for _, toField := range toType.Field.AllFields(nil, nil) {
		if toField.From != nil {
			assign.WriteString(toField.FullVariableName("") + " = ")

			fromField := toField.From
			if fromField.Options.Convert != "" {
				assign.WriteString(fromField.Options.Convert + "(" + fromField.FullVariableName("") + ")\n")
			} else {
				switch {
				case toField.FullDefinition() == fromField.FullDefinition():
					assign.WriteString(fromField.FullVariableName("") + "\n")
				case toField.FullDefinition()[1:] == fromField.FullDefinition():
					assign.WriteString("&" + fromField.FullVariableName("") + "\n")
				case toField.FullDefinition() == fromField.FullDefinition()[1:]:
					assign.WriteString("*" + fromField.FullVariableName("") + "\n")
				}
			}
		}
	}

	return assign.String()
}

// generateReturn generates a return statement for the function.
func generateReturn(function *models.Function) string {
	return "}"
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
				parameters.WriteString("ApplicationID: os.Getenv(\"APPLICATIONID\"")
			} else {
				parameters.WriteString(subfield.Name + ": os.Getenv(\"" + strings.ToUpper(subfield.Name) + "\")")
			}

			tagCount++
		}
	}

	return "{" + parameters.String() + "}"
}
