// Code generated by github.com/switchupcb/copygen
// DO NOT EDIT.
package requests

import (
	"fmt"
	"strconv"

	json "github.com/goccy/go-json"
	"github.com/switchupcb/disgo/wrapper/pkg/http"
	"github.com/switchupcb/disgo/wrapper/resources"
)

// Send sends a EditGlobalApplicationCommand to Discord and returns a ApplicationCommand.
func (r *EditGlobalApplicationCommand) Send(bot Client) (*resources.ApplicationCommand, error) {
	// Use pre-defined edge cases for certain fields (if necessary),
	// to provide feedback to the client prior to a request being sent.
	// if else return

	// send the request.
	var result *resources.ApplicationCommand
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while sending an EditGlobalApplicationCommand request\n%w", err)
	}

	err = http.SendRequestJSON(bot.client, bot.ctx, http.POST, EndpointEditGlobalApplicationCommand(strconv.FormatUint(uint64(bot.ApplicationID), 10), strconv.FormatUint(uint64(r.CommandID), 10)), body)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while sending an EditGlobalApplicationCommand request\n%w", err)
	}

	err = ParseResponseJSON(bot.ctx, result)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while sending an EditGlobalApplicationCommand request\n%w", err)
	}

	return result, nil
}

// Send command repeated for every request...
