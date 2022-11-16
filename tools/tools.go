package tools

import (
	"encoding/base64"
	"net/http"
)

// DataURI returns a Data URI from the given HTTP Content Type Header and base64 encoded data.
//
// https://en.wikipedia.org/wiki/Data_URI_scheme
func DataURI(contentType, base64EncodedData string) string {
	return "data:" + contentType + ";base64," + base64EncodedData
}

// ImageDataURI returns a Data URI from the given image data.
func ImageDataURI(img []byte) string {
	return DataURI(http.DetectContentType(img), base64.StdEncoding.EncodeToString(img))
}
