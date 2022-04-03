package resources

// Emoji Object
// https://discord.com/developers/docs/resources/emoji#emoji-object-emoji-structure
type Emoji struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Roles         []int64 `json:"roles,omitempty"`
	User          *User   `json:"user,omitempty"`
	RequireColons bool    `json:"require_colons,omitempty"`
	Managed       bool    `json:"managed,omitempty"`
	Animated      bool    `json:"animated,omitempty"`
	Available     bool    `json:"available,omitempty"`
}

// Reaction Object
// https://discord.com/developers/docs/resources/channel#reaction-object
type Reaction struct {
	count int    `json:"count"`
	me    bool   `json:"me"`
	emoji *Emoji `json:"emoji"`
}
