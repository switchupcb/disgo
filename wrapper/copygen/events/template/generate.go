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
	content.WriteString(generateRemove(functions) + "\n")
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
	fn.WriteString("func (bot *Client) Handle(eventname string, function interface{}) error {\n")
	fn.WriteString("bot.Handlers.mu.Lock()\n")
	fn.WriteString("defer bot.Handlers.mu.Unlock()\n")
	fn.WriteString("\n")
	fn.WriteString("switch eventname {\n")

	// write cases.
	cases := len(functions)
	for i, function := range functions {
		fn.WriteString(generateHandleCase(function.Name, generateIntentFlag(function.Options)))

		if i+1 != cases {
			fn.WriteString("\n")
		}
	}
	fn.WriteString("}\n")

	fn.WriteString("\nreturn fmt.Errorf(\"event handler for %s was not added.\", eventname)\n")
	fn.WriteString("}\n")
	return fn.String()
}

// generateHandleCase generates the switch case statement for the Handle function.
func generateHandleCase(eventname string, flag string) string {
	var c strings.Builder
	c.WriteString("case FlagGatewayEventName" + eventname + ":\n")

	// add automatic intent calculation.
	if flag != "" {
		c.WriteString("if !bot.Config.IntentSet[" + flag + "] {\n")
		c.WriteString("bot.Config.IntentSet[" + flag + "] = true\n")
		c.WriteString("bot.Config.Intents |= " + flag + "\n")
		c.WriteString("}\n\n")
	}

	// add the event handler.
	c.WriteString("if f, ok := function.(func(*" + eventname + ")); ok {\n")
	c.WriteString("bot.Handlers." + eventname + " = append(bot.Handlers." + eventname + ", f)\n")
	c.WriteString("return nil\n")
	c.WriteString("}\n")
	return c.String()
}

// generateIntentFlag generates the Intent Flag string for automatic intent calculation.
func generateIntentFlag(options models.FunctionOptions) string {
	if v, ok := options.Custom["intents"]; ok {
		intents := v[0]
		return strings.ReplaceAll(intents, " ", " | ")
	}

	return ""
}

////////////////////////////////////////////////////////////////////////////////
// Remove Func
////////////////////////////////////////////////////////////////////////////////

// generateRemove provides generated code for the end-user Remove function.
func generateRemove(functions []*models.Function) string {
	var fn strings.Builder
	fn.WriteString("// Remove removes the event handler at the given index from the bot.\n")
	fn.WriteString("// This function does NOT remove intents automatically.\n")
	fn.WriteString("func (bot *Client) Remove(eventname string, index int) error {\n")
	fn.WriteString("bot.Handlers.mu.Lock()\n")
	fn.WriteString("defer bot.Handlers.mu.Unlock()\n")
	fn.WriteString("\n")
	fn.WriteString("switch eventname {\n")

	// write cases.
	cases := len(functions)
	for i, function := range functions {
		fn.WriteString(generateRemoveCase(function.Name))

		if i+1 != cases {
			fn.WriteString("\n")
		}
	}
	fn.WriteString("}\n")

	fn.WriteString("\nreturn fmt.Errorf(\"event handler for %s at index %d was not removed.\", eventname, index)\n")
	fn.WriteString("}\n")
	return fn.String()
}

// generateRemoveCase generates the switch case statement for the Remove function.
func generateRemoveCase(eventname string) string {
	var c strings.Builder
	c.WriteString("case FlagGatewayEventName" + eventname + ":\n")

	// check the bounds of the handlers.
	c.WriteString("if len(bot.Handlers." + eventname + ") <= index {\n")
	c.WriteString("return fmt.Errorf(errRemoveInvalidEventHandler, eventname, index)\n")
	c.WriteString("}\n\n")

	// remove the event handler.
	c.WriteString("bot.Handlers." + eventname + " = " +
		"append(bot.Handlers." + eventname + "[:index], bot.Handlers." + eventname + "[index+1:]...)\n")
	c.WriteString("return nil\n")
	return c.String()
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
		fn.WriteString(generatehandleCase(function.Name))

		if i+1 != cases {
			fn.WriteString("\n")
		}
	}

	fn.WriteString("}\n")
	fn.WriteString("}\n")
	return fn.String()
}

// generatehandleCase generates the switch case statement for the handle function.
func generatehandleCase(eventname string) string {
	var c strings.Builder
	c.WriteString("case FlagGatewayEventName" + eventname + ":\n")
	c.WriteString("event := new(" + eventname + ")\n")
	c.WriteString("if err := json.Unmarshal(data, event); err != nil {\n")
	c.WriteString("log.Println(ErrorEvent{Event: FlagGatewayEventName" + eventname + ", Err: err, Action: ErrorEventActionUnmarshal}.Error())\n")
	c.WriteString("return\n")
	c.WriteString("}\n\n")

	// call the handlers.
	c.WriteString("for _, handler := range bot.Handlers." + eventname + " {\n")
	c.WriteString("go handler(event)\n")
	c.WriteString("}\n")
	return c.String()
}
