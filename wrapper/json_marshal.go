package wrapper

import (
	json "github.com/goccy/go-json"
)

func (r *BulkOverwriteGlobalApplicationCommands) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ApplicationCommands) // nolint:wrapcheck
}

func (r *BulkOverwriteGuildApplicationCommands) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ApplicationCommands) // nolint:wrapcheck
}

func (r *ModifyGuildChannelPositions) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Parameters) // nolint:wrapcheck
}

func (r *ModifyGuildRolePositions) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Parameters) // nolint:wrapcheck
}
