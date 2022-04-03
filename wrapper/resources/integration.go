package resources

import "time"

// Integration Object
// https://discord.com/developers/docs/resources/guild#integration-object
type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           bool               `json:"enabled"`
	Syncing           bool               `json:"syncing"`
	RoleID            string             `json:"role_id"`
	EnableEmoticons   bool               `json:"enable_emoticons"`
	ExpireBehavior    uint8              `json:"expire_behavior"`
	ExpireGracePeriod int                `json:"expire_grace_period"`
	User              *User              `json:"user"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          time.Time          `json:"synced_at"`
	SubscriberCount   int                `json:"subscriber_count"`
	Revoked           bool               `json:"revoked"`
	Application       Application        `json:"application"`
}

// Integration Expire Behaviors
// https://discord.com/developers/docs/resources/guild#integration-object-integration-expire-behaviors
const (
	REMOVEROLE = 0
	KICK       = 1
)

// Integration Account Object
// https://discord.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Integration Application Object
// https://discord.com/developers/docs/resources/guild#integration-application-object-integration-application-structure
type IntegrationApplication struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Bot         User   `json:"bot"`
}
