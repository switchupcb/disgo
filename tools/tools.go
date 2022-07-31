package tools

import (
	"encoding/base64"
	"net/http"

	disgo "github.com/switchupcb/disgo/wrapper"
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

// OptionsToMap parses an array of options and suboptions into an OptionMap.
func OptionsToMap(
	optionMap map[string]*disgo.ApplicationCommandInteractionDataOption,
	options []*disgo.ApplicationCommandInteractionDataOption,
	amount int,
) map[string]*disgo.ApplicationCommandInteractionDataOption {
	if optionMap == nil {
		optionMap = make(map[string]*disgo.ApplicationCommandInteractionDataOption, amount)
	}

	// add suboptions (slice by value is the most performant)
	for _, option := range options {
		optionMap[option.Name] = option
		if len(option.Options) != 0 {
			OptionsToMap(optionMap, option.Options, amount)
		}
	}

	return optionMap
}

// NumOptions determines the amount of options (and suboptions) in a given array of options.
func NumOptions(options []*disgo.ApplicationCommandInteractionDataOption) int {
	amount := len(options)

	// count the amount of suboptions.
	for _, option := range options {
		if len(option.Options) != 0 {
			amount += NumOptions(option.Options)
		}
	}

	return amount
}
