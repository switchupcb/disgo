package wrapper

import (
	"fmt"

	json "github.com/goccy/go-json"
)

/**unmarshal_embed.go contains custom UnmarshalJSON() functions.

Structs that contain an embedded field - which implements UnmarshalJSON() - will
use the embedded field's implementation of UnmarshalJSON(). These structs must
also implement UnmarshalJSON() to prevent null pointer dereferences.

--------------------------------------------------------------------------------------*/

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
