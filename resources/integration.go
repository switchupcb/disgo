package resources

import "time"

// Integration Object
// https://discord.com/developers/docs/resources/guild#integration-object
type Integration struct {
	ID                Snowflake          `json:"id,omitempty"`
	Name              string             `json:"name,omitempty"`
	Type              string             `json:"type,omitempty"`
	Enabled           bool               `json:"enabled,omitempty"`
	Syncing           bool               `json:"syncing,omitempty"`
	RoleID            Snowflake          `json:"role_id,omitempty"`
	EnableEmoticons   bool               `json:"enable_emoticons,omitempty"`
	ExpireBehavior    Flag               `json:"expire_behavior,omitempty"`
	ExpireGracePeriod int                `json:"expire_grace_period,omitempty"`
	User              *User              `json:"user,omitempty"`
	Account           IntegrationAccount `json:"account,omitempty"`
	SyncedAt          time.Time          `json:"synced_at,omitempty"`
	SubscriberCount   int                `json:"subscriber_count,omitempty"`
	Revoked           bool               `json:"revoked,omitempty"`
	Application       Application        `json:"application,omitempty"`
}

// Integration Expire Behaviors
// https://discord.com/developers/docs/resources/guild#integration-object-integration-expire-behaviors
const (
	FlagIntegrationExpireBehaviorsREMOVEROLE = 0
	FlagIntegrationExpireBehaviorsKICK       = 1
)

// Integration Account Object
// https://discord.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Integration Application Object
// https://discord.com/developers/docs/resources/guild#integration-application-object-integration-application-structure
type IntegrationApplication struct {
	ID          Snowflake `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	Description string    `json:"description,omitempty"`
	Bot         User      `json:"bot,omitempty"`
}
