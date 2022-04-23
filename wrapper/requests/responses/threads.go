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

// List Public Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListPublicArchivedThreads struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
}

// List Private Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListPrivateArchivedThreads struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
}

// List Joined Private Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListJoinedPrivateArchivedThreads struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
}

// create threads for things that return more than 1 object
// create new files if not already existed under responses package
