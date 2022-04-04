package resources

// Activity Object
// https://discord.com/developers/docs/topics/gateway#activity-object-activity-structure
type Activity struct {
	Name          string              `json:"name,omitempty"`
	Type          *Flag               `json:"type,omitempty"`
	URL           string              `json:"url,omitempty"`
	CreatedAt     int                 `json:"created_at,omitempty"`
	Timestamps    *ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID Snowflake           `json:"application_id,omitempty"`
	Details       string              `json:"details,omitempty"`
	State         string              `json:"state,omitempty"`
	Emoji         *Emoji              `json:"emoji,omitempty"`
	Party         *ActivityParty      `json:"party,omitempty"`
	Assets        *ActivityAssets     `json:"assets,omitempty"`
	Secrets       *ActivitySecrets    `json:"secrets,omitempty"`
	Instance      bool                `json:"instance,omitempty"`
	Flags         BitFlag             `json:"flags,omitempty"`
	Buttons       []Button            `json:"buttons,omitempty"`
}

// ActivityTimestamps Struct
// https://discord.com/developers/docs/game-sdk/activities#data-models-activitytimestamps-struct
type ActivityTimestamps struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// ActivityAssets Struct
// https://discord.com/developers/docs/game-sdk/activities#data-models-activitytimestamps-struct
type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

// ActivityParty Struct
// https://discord.com/developers/docs/game-sdk/activities#data-models-activityparty-struct
type ActivityParty struct {
	ID   string `json:"id,omitempty"`
	Size []int  `json:"size,omitempty"`
}

// PartySize Struct
// https://discord.com/developers/docs/game-sdk/activities#data-models-partysize-struct
type PartySize struct {
	CurrentSize int32 `json:"current_size,omitempty"`
	MaxSize     int32 `json:"max_size,omitempty"`
}

// ActivitySecrets Struct
// https://discord.com/developers/docs/game-sdk/activities#data-models-activitysecrets-struct
type ActivitySecrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

// ActivityType Enum
// https://discord.com/developers/docs/game-sdk/activities#data-models-activitytype-enum
const (
	FlagActivityTypeEnumPlaying   = 0
	FlagActivityTypeEnumStreaming = 1
	FlagActivityTypeEnumListening = 2
	FlagActivityTypeEnumWatching  = 3
	FlagActivityTypeEnumCustom    = 4
	FlagActivityTypeEnumCompeting = 5
)

// ActivityJoinRequestReply Enum
// https://discord.com/developers/docs/game-sdk/activities#data-models-activityjoinrequestreply-enum
const (
	FlagActivityJoinRequestReplyEnumNo     = 0
	FlagActivityJoinRequestReplyEnumYes    = 1
	FlagActivityJoinRequestReplyEnumIgnore = 2
)

// ActivityActionType Enum
// https://discord.com/developers/docs/game-sdk/activities#data-models-activityactiontype-enum
const (
	FlagActivityActionTypeEnumJoin     = 1
	FlagActivityActionTypeEnumSpectate = 2
)
