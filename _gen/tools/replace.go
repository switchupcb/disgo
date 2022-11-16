package tools

import (
	"fmt"
	"go/format"
	"strings"

	"golang.org/x/exp/slices"
)

// TypeFix replaces dasgo references to a type with another type.
func TypeFix(data []byte) ([]byte, error) {
	content := string(data)

	content = snowflake(content)

	const (
		commentNonce    = "// Nonce represents a Discord nonce (integer or string)."
		definitionNonce = "type Nonce interface{}"

		commentValue    = "// Value represents a value (string, integer, or double)."
		definitionValue = "type Value interface{}"
	)

	content = field(content, "Nonce", "json.Number", []string{commentNonce, definitionNonce}...)
	content = field(content, "Value", "json.Number", []string{commentValue, definitionValue}...)

	// gofmt
	contentdata := []byte(content)
	fmtdata, err := format.Source(contentdata)
	if err != nil {
		return contentdata, fmt.Errorf("an error occurred while formatting the generated code.\n%w", err)
	}

	return fmtdata, nil
}

// field replaces dasgo references of the given "field" with the replaced field.
//
// skip skips lines which match any string in the given parameters.
func field(content string, field string, replace string, skip ...string) string {
	var keep strings.Builder

	for _, line := range strings.Split(content, "\n") {
		// remove the skipped lines.
		if slices.Contains(skip, line) {
			continue
		}

		// replace lines with two occurrences or more of the field (i.e Nonce Nonce).
		lineFields := strings.Fields(line)
		count := 0
		for _, lineField := range lineFields {
			if lineField == field {
				count++
			}
		}

		if count == 2 {
			keep.WriteString(
				fmt.Sprintf("%s %s %s\n", field, replace, strings.Join(lineFields[2:], " ")),
			)

			continue
		}

		keep.WriteString(line + "\n")
	}

	return keep.String()
}

// snowflake replaces dasgo references of Snowflake with string.
func snowflake(content string) string {
	const (
		commentSnowflake    = "// Snowflake represents a Discord API Snowflake."
		definitionSnowflake = "type Snowflake uint64"

		escapedCommentSnowflake    = "// $nowflake represents a Discord API $nowflake."
		escapedDefinitionSnowflake = "type $nowflake uint64"
	)

	// escape the definition of a Snowflake.
	content = strings.Replace(content, commentSnowflake, escapedCommentSnowflake, 1)
	content = strings.Replace(content, definitionSnowflake, escapedDefinitionSnowflake, 1)

	// replace Snowflake.
	content = strings.ReplaceAll(content, "Snowflake", "string")

	// unescape the definition of a Snowflake.
	content = strings.Replace(content, escapedCommentSnowflake, commentSnowflake, 1)
	content = strings.Replace(content, escapedDefinitionSnowflake, definitionSnowflake, 1)

	return content
}
