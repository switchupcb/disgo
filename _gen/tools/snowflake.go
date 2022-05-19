package tools

import (
	"fmt"
	"go/format"
	"strings"
)

const (
	commentSnowflake    = "// Snowflake represents a Discord API Snowflake."
	definitionSnowflake = "type Snowflake uint64"

	escapedCommentSnowflake    = "// $nowflake represents a Discord API $nowflake."
	escapedDefinitionSnowflake = "type $nowflake uint64"
)

// Snowflake replaces dasgo references of Snowflake with string.
func Snowflake(data []byte) ([]byte, error) {
	content := string(data)

	// escape the definition of a Snowflake.
	content = strings.Replace(content, commentSnowflake, escapedCommentSnowflake, 1)
	content = strings.Replace(content, definitionSnowflake, escapedDefinitionSnowflake, 1)

	// replace Snowflake.
	content = strings.ReplaceAll(content, "Snowflake", "string")

	// unescape the definition of a Snowflake.
	content = strings.Replace(content, escapedCommentSnowflake, commentSnowflake, 1)
	content = strings.Replace(content, escapedDefinitionSnowflake, definitionSnowflake, 1)

	// gofmt
	contentdata := []byte(content)
	fmtdata, err := format.Source(contentdata)
	if err != nil {
		return contentdata, fmt.Errorf("an error occurred while formatting the generated code.\n%w", err)
	}

	return fmtdata, nil
}
