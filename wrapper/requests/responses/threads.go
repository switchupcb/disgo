package responses

import (
	"github.com/switchupcb/disgo/wrapper/resources"
)

// List Active Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListActiveThreadsResponse struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
}

// List Public Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListPublicArchivedThreadsResponse struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
}

// List Private Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListPrivateArchivedThreadsResponse struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
}

// List Joined Private Archived Threads Response Body
// https://discord.com/developers/docs/resources/channel#list-active-threads-response-body
type ListJoinedPrivateArchivedThreadsResponse struct {
	Threads []*resources.Channel      `json:"threads"`
	Members []*resources.ThreadMember `json:"members"`
	HasMore bool                      `json:"has_more"`
}
