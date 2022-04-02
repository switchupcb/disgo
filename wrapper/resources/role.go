package resource

// Role Object
// https://discord.com/developers/docs/topics/permissions#role-object
type Role struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Color        uint     `json:"color"`
	Hoist        bool     `json:"hoist"`
	Icon         string   `json:"string"`
	UnicodeEmoji string   `json:"string"`
	Position     int      `json:"position"`
	Permissions  int      `json:"permissions"`
	Managed      bool     `json:"managed"`
	Mentionable  bool     `json:"mentionable"`
	Tags         RoleTags `json:"tags"`
}

// Role Tags Structure
// https://discord.com/developers/docs/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	BotID             int64 `json:"bot_id"`
	IntegrationID     int64 `json:"integration_id"`
	PremiumSubscriber bool  `json:"premium_subscriber"`
}
