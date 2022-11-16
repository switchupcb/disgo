package wrapper

import (
	"fmt"

	json "github.com/goccy/go-json"
)

/**unmarshal.go contains custom UnmarshalJSON() functions.

This enables json.Unmarshal() to unmarshal into types that contain fields that are interfaces.

In addition, structs that contain an embedded field - that implements UnmarshalJSON() - will
use the embedded field's implementation of UnmarshalJSON(). As a result, these structs must
also implement UnmarshalJSON() to prevent null pointer dereferences. */

/** Unused: Command, Event */

/** Nonce
Converted to string using ./_gen typefix operation.
Includes: CreateMessage, Message */

/** Value
Converted to string using ./_gen typefix operation.
Includes: ApplicationCommandOptionChoice, ApplicationCommandInteractionDataOption */

/** Component */

// unmarshalComponents unmarshals a JSON component array into a slice of Go Interface Components (with underlying structs).
func unmarshalComponents(b []byte) ([]Component, error) {
	// Components are always provided in a JSON array.
	// Create a variable (of type []struct) that can read all of the Component Types.
	var unmarshalledComponents []struct {

		// https://discord.com/developers/docs/interactions/message-components#component-object-example-component
		Type uint `json:"type"`
	}

	// unmarshal the JSON (components.{component.Type}) into unmarshalledComponents.
	if err := json.Unmarshal(b, &unmarshalledComponents); err != nil {
		return nil, fmt.Errorf(errUnmarshal, unmarshalledComponents, err)
	}

	// use the known component types to return a slice of Go Interface Components with underlying structs.
	components := make([]Component, len(unmarshalledComponents))
	for i, unmarshalledComponent := range unmarshalledComponents {
		var component Component

		// set the component (interface) to an underlying type.
		switch unmarshalledComponent.Type {
		case FlagComponentTypeActionRow:
			component = &ActionsRow{} //nolint:exhaustruct

		case FlagComponentTypeButton:
			component = &Button{} //nolint:exhaustruct

		case FlagComponentTypeSelectMenu,
			FlagComponentTypeUserSelect,
			FlagComponentTypeRoleSelect,
			FlagComponentTypeMentionableSelect,
			FlagComponentTypeChannelSelect:
			component = &SelectMenu{} //nolint:exhaustruct

		case FlagComponentTypeTextInput:
			component = &TextInput{} //nolint:exhaustruct

		default:
			return nil, fmt.Errorf(
				"attempt to unmarshal into unknown component type (%d)",
				unmarshalledComponent.Type,
			)
		}

		components[i] = component
	}

	return components, nil
}

func (r *EditOriginalInteractionResponse) UnmarshalJSON(b []byte) error {
	// The following pattern is present throughout this file
	// in order to prevent a stack overflow (of r.UnmarshalJSON()).
	type alias EditOriginalInteractionResponse
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(EditOriginalInteractionResponse)
	}

	*r = EditOriginalInteractionResponse(*d)

	return nil
}

func (r *CreateFollowupMessage) UnmarshalJSON(b []byte) error {
	type alias CreateFollowupMessage
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(CreateFollowupMessage)
	}

	*r = CreateFollowupMessage(*d)

	return nil
}

func (r *EditFollowupMessage) UnmarshalJSON(b []byte) error {
	type alias EditFollowupMessage
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(EditFollowupMessage)
	}

	*r = EditFollowupMessage(*d)

	return nil
}

func (r *EditMessage) UnmarshalJSON(b []byte) error {
	type alias EditMessage
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(EditMessage)
	}

	*r = EditMessage(*d)

	return nil
}

func (r *ForumThreadMessageParams) UnmarshalJSON(b []byte) error {
	type alias ForumThreadMessageParams
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(ForumThreadMessageParams)
	}

	*r = ForumThreadMessageParams(*d)

	return nil
}

func (r *ExecuteWebhook) UnmarshalJSON(b []byte) error {
	type alias ExecuteWebhook
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(ExecuteWebhook)
	}

	*r = ExecuteWebhook(*d)

	return nil
}

func (r *EditWebhookMessage) UnmarshalJSON(b []byte) error {
	type alias EditWebhookMessage
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(EditWebhookMessage)
	}

	*r = EditWebhookMessage(*d)

	return nil
}

func (r *ActionsRow) UnmarshalJSON(b []byte) error {
	type alias ActionsRow
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(ActionsRow)
	}

	*r = ActionsRow(*d)

	return nil
}

func (r *ModalSubmitData) UnmarshalJSON(b []byte) error {
	type alias ModalSubmitData
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(ModalSubmitData)
	}

	*r = ModalSubmitData(*d)

	return nil
}

func (r *Messages) UnmarshalJSON(b []byte) error {
	type alias Messages
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(Messages)
	}

	*r = Messages(*d)

	return nil
}

func (r *Modal) UnmarshalJSON(b []byte) error {
	type alias Modal
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(Modal)
	}

	*r = Modal(*d)

	return nil
}

func (r *Message) UnmarshalJSON(b []byte) error {
	type alias Message
	d := new(alias)

	var err error
	if d.Components, err = unmarshalComponents(b); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if err = json.Unmarshal(b, d); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(Message)
	}

	*r = Message(*d)

	return nil
}

/** InteractionData */

// unmarshalInteractionData unmarshals a JSON InteractionData object into
// a Go Interface InteractionData (with an underlying struct).
func unmarshalInteractionData(b json.RawMessage, x uint8) (InteractionData, error) {
	var interactionData InteractionData

	// use the known Interaction Data type to return
	// a Go Interface InteractionData with an underlying struct.
	switch x {
	case FlagInteractionTypePING:
		return nil, nil

	case FlagInteractionTypeAPPLICATION_COMMAND,
		FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		interactionData = &ApplicationCommandData{} //nolint:exhaustruct

	case FlagInteractionTypeMESSAGE_COMPONENT:
		interactionData = &MessageComponentData{} //nolint:exhaustruct

	case FlagInteractionTypeMODAL_SUBMIT:
		interactionData = &ModalSubmitData{} //nolint:exhaustruct
	}

	if interactionData == nil {
		return nil, fmt.Errorf("attempt to unmarshal into unknown interaction data type (%d)", x)
	}

	// unmarshal into the underlying struct.
	if err := json.Unmarshal(b, interactionData); err != nil {
		return nil, fmt.Errorf(errUnmarshal, interactionData, err)
	}

	return interactionData, nil
}

func (r *Interaction) UnmarshalJSON(b []byte) error {
	type alias Interaction

	var unmarshalledInteraction struct {
		alias
		Data json.RawMessage `json:"data,omitempty"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalledInteraction); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalledInteraction.alias.Data, err =
		unmarshalInteractionData(unmarshalledInteraction.Data, uint8(unmarshalledInteraction.Type)); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(Interaction)
	}

	*r = Interaction(unmarshalledInteraction.alias)

	return nil
}

/** InteractionCallbackData */

// unmarshalInteractionCallbackData unmarshals a JSON InteractionCallbackData object into
// a Go Interface InteractionCallbackData (with an underlying struct).
func unmarshalInteractionCallbackData(b []byte, x uint8) (InteractionCallbackData, error) {
	var interactionCallbackData InteractionCallbackData

	// use the known Interaction Callback Data type to return
	// a Go Interface InteractionCallbackData with an underlying struct.
	switch x {
	case FlagInteractionCallbackTypePONG:
		return nil, nil // Ping

	case FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
		FlagInteractionCallbackTypeUPDATE_MESSAGE:
		interactionCallbackData = &Messages{} //nolint:exhaustruct

	case FlagInteractionCallbackTypeDEFERRED_CHANNEL_MESSAGE_WITH_SOURCE:
		return nil, nil // Edit a followup response later.

	case FlagInteractionCallbackTypeDEFERRED_UPDATE_MESSAGE:
		return nil, nil // Edit the original response later.

	case FlagInteractionCallbackTypeAPPLICATION_COMMAND_AUTOCOMPLETE_RESULT:
		interactionCallbackData = &Autocomplete{} //nolint:exhaustruct

	case FlagInteractionCallbackTypeMODAL:
		interactionCallbackData = &Modal{} //nolint:exhaustruct
	}

	if interactionCallbackData == nil {
		return nil, fmt.Errorf(
			"attempt to unmarshal into unknown interaction callback data type (%d)",
			x)
	}

	// unmarshal into the underlying struct.
	if err := json.Unmarshal(b, interactionCallbackData); err != nil {
		return nil, fmt.Errorf(errUnmarshal, interactionCallbackData, err)
	}

	return interactionCallbackData, nil
}

func (r *InteractionResponse) UnmarshalJSON(b []byte) error {
	type alias InteractionResponse

	var unmarshalledInteractionResponse struct {
		alias
		Data json.RawMessage `json:"data,omitempty"`
	}

	var err error
	if err = json.Unmarshal(b, &unmarshalledInteractionResponse); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if unmarshalledInteractionResponse.alias.Data, err =
		unmarshalInteractionCallbackData(
			unmarshalledInteractionResponse.Data,
			uint8(unmarshalledInteractionResponse.Type)); err != nil {
		return fmt.Errorf(errUnmarshal, r, err)
	}

	if r == nil {
		r = new(InteractionResponse)
	}

	*r = InteractionResponse(unmarshalledInteractionResponse.alias)

	return nil
}

/** Structs that contain embedded fields that implement UnmarshalJSON() */

func (e *MessageCreate) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.Message); err != nil {
		return fmt.Errorf(errUnmarshal, e, err)
	}

	return nil
}

func (e *MessageUpdate) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.Message); err != nil {
		return fmt.Errorf(errUnmarshal, e, err)
	}

	return nil
}

func (e *InteractionCreate) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.Interaction); err != nil {
		return fmt.Errorf(errUnmarshal, e, err)
	}

	return nil
}

func (e *CreateInteractionResponse) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &e.InteractionResponse); err != nil {
		return fmt.Errorf(errUnmarshal, e, err)
	}

	return nil
}
