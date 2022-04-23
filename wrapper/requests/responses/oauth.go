package responses

import (
	"time"

	"github.com/switchupcb/disgo/wrapper/resources"
)

// Current Authorization Information Response Structure
// https://discord.com/developers/docs/topics/oauth2#get-current-authorization-information
type CurrentAuthorizationInformation struct {
	Application *resources.Application `json:"application"`
	Scopes      []*int                 `json:"scopes"`
	Expires     *time.Time             `json:"expires"`
	User        *resources.User        `json:"user"`
}
