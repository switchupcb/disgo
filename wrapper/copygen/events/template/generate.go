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
	functions := make([]*models.Function, len(gen.Functions))
	for i := range gen.Functions {
		functions[i] = &gen.Functions[i]
	}

	var content strings.Builder
	content.WriteString(string(gen.Keep) + "\n")
	content.WriteString(generateHandlers(functions) + "\n")
	content.WriteString(generateHandle(functions) + "\n")
	content.WriteString(generatehandle(functions) + "\n")
	return content.String(), nil
}

////////////////////////////////////////////////////////////////////////////////
// Handlers Struct
////////////////////////////////////////////////////////////////////////////////

// generateHandlers generates the Handlers struct.
func generateHandlers(functions []*models.Function) string {
	var strct strings.Builder
	strct.WriteString("// Handlers represents a bot's event handlers.\n")
	strct.WriteString("type Handlers struct {\n")

	// write the body.
	for _, function := range functions {
		strct.WriteString(function.Name + " []func(*" + function.Name + ")\n")
	}

	// add manual fields.
	strct.WriteString("mu sync.RWMutex\n")

	strct.WriteString("}\n")
	return strct.String()
}

////////////////////////////////////////////////////////////////////////////////
// Handle Func
////////////////////////////////////////////////////////////////////////////////

// generateHandle provides generated code for the end-user Handle function.
func generateHandle(functions []*models.Function) string {
	var fn strings.Builder
	fn.WriteString("// Handle adds an event handler for the given event to the bot.\n")
	fn.WriteString("func (bot *Client) Handle(eventname string, function interface{}) {\n")
	fn.WriteString("bot.Handlers.mu.Lock()\n")
	fn.WriteString("defer bot.Handlers.mu.Unlock()\n")
	fn.WriteString("\n")
	fn.WriteString("switch eventname {\n")

	// write cases.
	cases := len(functions)
	for i, function := range functions {
		fn.WriteString("case FlagGatewayEventName" + function.Name + ":\n")
		fn.WriteString("if f, ok := function.(func(*" + function.Name + ")); ok {\n")
		fn.WriteString("bot.Handlers." + function.Name + " = append(bot.Handlers." + function.Name + ", f)\n")
		fn.WriteString("return\n")
		fn.WriteString("}\n")

		if i+1 != cases {
			fn.WriteString("\n")
		}
	}
	fn.WriteString("}\n")

	fn.WriteString("\nlog.Printf(\"Event Handler for %s was not added.\", eventname)\n")
	fn.WriteString("}\n")
	return fn.String()
}

////////////////////////////////////////////////////////////////////////////////
// handle Func
////////////////////////////////////////////////////////////////////////////////

// generatehandle provides generated code for the event handler handle function.
func generatehandle(functions []*models.Function) string {
	var fn strings.Builder
	fn.WriteString("// handle handles an event using its name and data.\n")
	fn.WriteString("func (bot *Client) handle(eventname string, data json.RawMessage) {\n")
	fn.WriteString("bot.Handlers.mu.RLock()\n")
	fn.WriteString("defer bot.Handlers.mu.RUnlock()\n")
	fn.WriteString("\n")
	fn.WriteString("switch eventname {\n")

	// write cases.
	cases := len(functions)
	for i, function := range functions {
		fn.WriteString(generateCase(function.Name))

		if i+1 != cases {
			fn.WriteString("\n")
		}
	}

	fn.WriteString("}\n")
	fn.WriteString("}\n")
	return fn.String()
}

// generateCase generates the switch case statement for the handle function.
func generateCase(eventname string) string {
	var c strings.Builder
	c.WriteString("case FlagGatewayEventName" + eventname + ":\n")
	c.WriteString("var event *" + eventname + "\n")
	c.WriteString("if err := json.Unmarshal(data, event); err != nil {\n")
	c.WriteString("log.Printf(ErrLogEventUnmarshal, eventname, err)\n")
	c.WriteString("return\n")
	c.WriteString("}\n\n")

	// call the handlers.
	c.WriteString("for _, handler := range bot.Handlers." + eventname + " {\n")
	c.WriteString("go handler(event)\n")
	c.WriteString("}\n")
	return c.String()
}
