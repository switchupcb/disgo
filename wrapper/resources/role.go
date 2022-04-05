package resources

// Role Object
// https://discord.com/developers/docs/topics/permissions#role-object
type Role struct {
	ID           Snowflake `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Color        uint      `json:"color,omitempty"`
	Hoist        bool      `json:"hoist,omitempty"`
	Icon         string    `json:"string,omitempty"`
	UnicodeEmoji string    `json:"string,omitempty"`
	Position     int       `json:"position,omitempty"`
	Permissions  string    `json:"permissions,omitempty"`
	Managed      bool      `json:"managed,omitempty"`
	Mentionable  bool      `json:"mentionable,omitempty"`
	Tags         RoleTags  `json:"tags,omitempty"`
}

// Role Tags Structure
// https://discord.com/developers/docs/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	BotID             Snowflake `json:"bot_id,omitempty"`
	IntegrationID     Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber bool      `json:"premium_subscriber,omitempty"`
}
