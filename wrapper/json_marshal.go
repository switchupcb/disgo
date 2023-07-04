package wrapper

import (
	"bytes"
	"strconv"

	json "github.com/goccy/go-json"
)

var (
	byteEmptySlice = []byte("[]")
)

func (t Flags) MarshalJSON() ([]byte, error) {
	if len(t) == 0 {
		return byteEmptySlice, nil
	}

	var output bytes.Buffer
	output.WriteByte('[')

	stop := len(t) - 1
	for i, f := range t {
		output.WriteString(strconv.Itoa(int(f)))

		if i == stop {
			break
		}

		output.WriteByte(',')
	}

	output.WriteByte(']')

	return output.Bytes(), nil
}

func (r *BulkOverwriteGlobalApplicationCommands) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ApplicationCommands) //nolint:wrapcheck
}

func (r *BulkOverwriteGuildApplicationCommands) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ApplicationCommands) //nolint:wrapcheck
}

func (r *ModifyGuildChannelPositions) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Parameters) //nolint:wrapcheck
}

func (r *ModifyGuildRolePositions) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Parameters) //nolint:wrapcheck
}
