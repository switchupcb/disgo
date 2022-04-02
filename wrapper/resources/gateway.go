package resource

// Presence Update Event Fields
// https://discord.com/developers/docs/topics/gateway#presence-update-presence-update-event-fields
type PresenceUpdate struct {
	User         *User        `json:"user"`
	GuildID      int64        `json:"guild_id"`
	Status       string       `json:"status"`
	Activities   []*Activity  `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
}

// Client Status Object
// https://discord.com/developers/docs/topics/gateway#client-status-object
type ClientStatus struct {
	Desktop string `json:"desktop"`
	Mobile  string `json:"mobile"`
	Web     string `json:"web"`
}

// Activity Object
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-structure
type Activity struct {
	Name          string      `json:"name"`
	Type          uint8       `json:"type"`
	URL           string      `json:"url,omitempty"`
	CreatedAt     int         `json:"created_at"`
	Timestamps    *Timestamps `json:"timestamps,omitempty"`
	ApplicationID int64       `json:"application_id,omitempty"`
	Details       string      `json:"details,omitempty"`
	State         string      `json:"state,omitempty"`
	Emoji         *Emoji      `json:"emoji,omitempty"`
	Party         *Party      `json:"party,omitempty"`
	Assets        *Assets     `json:"assets,omitempty"`
	Secrets       *Secrets    `json:"secrets,omitempty"`
	Instance      bool        `json:"instance,omitempty"`
	Flags         uint8       `json:"flags,omitempty"`
	Buttons       []Button    `json:"buttons,omitempty"`
}
