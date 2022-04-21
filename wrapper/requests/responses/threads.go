package responses

import (
	"github.com/switchupcb/disgo/wrapper/resources"
)

// List Active Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListActiveThreads struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
}
