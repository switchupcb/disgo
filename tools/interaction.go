package tools

import (
	"github.com/switchupcb/disgo"
)

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
