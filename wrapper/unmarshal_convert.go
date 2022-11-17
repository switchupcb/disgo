package wrapper

import (
	"fmt"
	"strconv"
)

/**unmarshal_convert.go contains type conversion functions for interfaces.

This enables users (developers) to easily type convert interfaces. */

const (
	errTypeConvert = "attempted to type convert InteractionData of type %v to type %s"
)

/* Nonce */

func (n Nonce) String() string {
	return string(n)
}

func (n Nonce) Int64() (int64, error) {
	return strconv.ParseInt(string(n), base10, bit64) //nolint:wrapcheck
}

/* Value */

func (n Value) String() string {
	return string(n)
}

func (n Value) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), bit64) //nolint:wrapcheck
}

func (n Value) Int64() (int64, error) {
	return strconv.ParseInt(string(n), base10, bit64) //nolint:wrapcheck
}

/* InteractionData */

// ApplicationCommand type converts an InteractionData field into an ApplicationCommandData struct.
func (i *Interaction) ApplicationCommand() *ApplicationCommandData {
	switch i.Data.InteractionDataType() {
	case FlagInteractionTypeAPPLICATION_COMMAND,
		FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		return i.Data.(*ApplicationCommandData) //nolint:forcetypeassert

	case FlagInteractionTypePING:
		panic(fmt.Sprintf(errTypeConvert, "Ping", "ApplicationCommandData"))

	case FlagInteractionTypeMESSAGE_COMPONENT:
		panic(fmt.Sprintf(errTypeConvert, "MessageComponentData", "ApplicationCommandData"))

	case FlagInteractionTypeMODAL_SUBMIT:
		panic(fmt.Sprintf(errTypeConvert, "ModalSubmitData", "ApplicationCommandData"))
	}

	panic(fmt.Sprintf(errTypeConvert, i.Data.InteractionDataType(), "ApplicationCommandData"))
}

// MessageComponent type converts an InteractionData field into a MessageComponentData struct.
func (i *Interaction) MessageComponent() *MessageComponentData {
	switch i.Data.InteractionDataType() {
	case FlagInteractionTypeMESSAGE_COMPONENT:
		return i.Data.(*MessageComponentData) //nolint:forcetypeassert

	case FlagInteractionTypePING:
		panic(fmt.Sprintf(errTypeConvert, "Ping", "MessageComponentData"))

	case FlagInteractionTypeAPPLICATION_COMMAND,
		FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		panic(fmt.Sprintf(errTypeConvert, "ApplicationCommandData", "MessageComponentData"))

	case FlagInteractionTypeMODAL_SUBMIT:
		panic(fmt.Sprintf(errTypeConvert, "ModalSubmitData", "MessageComponentData"))
	}

	panic(fmt.Sprintf(errTypeConvert, i.Data.InteractionDataType(), "MessageComponentData"))
}

// ModalSubmit type converts an InteractionData field into a ModalSubmitData struct.
func (i *Interaction) ModalSubmit() *ModalSubmitData {
	switch i.Data.InteractionDataType() {
	case FlagInteractionTypeMODAL_SUBMIT:
		return i.Data.(*ModalSubmitData) //nolint:forcetypeassert

	case FlagInteractionTypePING:
		panic(fmt.Sprintf(errTypeConvert, "Ping", "ModalSubmitData"))

	case FlagInteractionTypeAPPLICATION_COMMAND,
		FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		panic(fmt.Sprintf(errTypeConvert, "ApplicationCommandData", "ModalSubmitData"))

	case FlagInteractionTypeMESSAGE_COMPONENT:
		panic(fmt.Sprintf(errTypeConvert, "MessageComponentData", "ModalSubmitData"))
	}

	panic(fmt.Sprintf(errTypeConvert, i.Data.InteractionDataType(), "ModalSubmitData"))
}
