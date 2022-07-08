package wrapper

import (
	"net/url"

	"github.com/gorilla/schema"
)

var (
	// qsEncoder is used to create URL Query Strings from objects.
	qsEncoder = schema.NewEncoder()
)

// init runs at the start of the program.
func init() {
	// use `url` tags for the URL Query String encoder and decoder.
	qsEncoder.SetAliasTag("url")
}

// EndpointQueryString returns a URL Query String from a given object.
func EndpointQueryString(dst any) (string, error) {
	params := url.Values{}
	err := qsEncoder.Encode(dst, params)
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	return params.Encode(), nil
}
