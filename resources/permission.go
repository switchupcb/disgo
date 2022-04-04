package resources

// Bitwise Permission Flags
// https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
const (
	FlagBitwisePermissionFlagsCREATE_INSTANT_INVITE      = 1 << 0
	FlagBitwisePermissionFlagsKICK_MEMBERS               = 1 << 1
	FlagBitwisePermissionFlagsBAN_MEMBERS                = 1 << 2
	FlagBitwisePermissionFlagsADMINISTRATOR              = 1 << 3
	FlagBitwisePermissionFlagsMANAGE_CHANNELS            = 1 << 4
	FlagBitwisePermissionFlagsMANAGE_GUILD               = 1 << 5
	FlagBitwisePermissionFlagsADD_REACTIONS              = 1 << 6
	FlagBitwisePermissionFlagsVIEW_AUDIT_LOG             = 1 << 7
	FlagBitwisePermissionFlagsPRIORITY_SPEAKER           = 1 << 8
	FlagBitwisePermissionFlagsSTREAM                     = 1 << 9
	FlagBitwisePermissionFlagsVIEW_CHANNEL               = 1 << 10
	FlagBitwisePermissionFlagsSEND_MESSAGES              = 1 << 11
	FlagBitwisePermissionFlagsSEND_TTS_MESSAGES          = 1 << 12
	FlagBitwisePermissionFlagsMANAGE_MESSAGES            = 1 << 13
	FlagBitwisePermissionFlagsEMBED_LINKS                = 1 << 14
	FlagBitwisePermissionFlagsATTACH_FILES               = 1 << 15
	FlagBitwisePermissionFlagsREAD_MESSAGE_HISTORY       = 1 << 16
	FlagBitwisePermissionFlagsMENTION_EVERYONE           = 1 << 17
	FlagBitwisePermissionFlagsUSE_EXTERNAL_EMOJIS        = 1 << 18
	FlagBitwisePermissionFlagsVIEW_GUILD_INSIGHTS        = 1 << 19
	FlagBitwisePermissionFlagsCONNECT                    = 1 << 20
	FlagBitwisePermissionFlagsSPEAK                      = 1 << 21
	FlagBitwisePermissionFlagsMUTE_MEMBERS               = 1 << 22
	FlagBitwisePermissionFlagsDEAFEN_MEMBERS             = 1 << 23
	FlagBitwisePermissionFlagsMOVE_MEMBERS               = 1 << 24
	FlagBitwisePermissionFlagsUSE_VAD                    = 1 << 25
	FlagBitwisePermissionFlagsCHANGE_NICKNAME            = 1 << 26
	FlagBitwisePermissionFlagsMANAGE_NICKNAMES           = 1 << 27
	FlagBitwisePermissionFlagsMANAGE_ROLES               = 1 << 28
	FlagBitwisePermissionFlagsMANAGE_WEBHOOKS            = 1 << 29
	FlagBitwisePermissionFlagsMANAGE_EMOJIS_AND_STICKERS = 1 << 30
	FlagBitwisePermissionFlagsUSE_APPLICATION_COMMANDS   = 1 << 31
	FlagBitwisePermissionFlagsREQUEST_TO_SPEAK           = 1 << 32
	FlagBitwisePermissionFlagsMANAGE_EVENTS              = 1 << 33
	FlagBitwisePermissionFlagsMANAGE_THREADS             = 1 << 34
	FlagBitwisePermissionFlagsCREATE_PUBLIC_THREADS      = 1 << 35
	FlagBitwisePermissionFlagsCREATE_PRIVATE_THREADS     = 1 << 36
	FlagBitwisePermissionFlagsUSE_EXTERNAL_STICKERS      = 1 << 37
	FlagBitwisePermissionFlagsSEND_MESSAGES_IN_THREADS   = 1 << 38
	FlagBitwisePermissionFlagsUSE_EMBEDDED_ACTIVITIES    = 1 << 39
	FlagBitwisePermissionFlagsMODERATE_MEMBERS           = 1 << 40
)

// Overwrite Object
// https://discord.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    Snowflake `json:"id,omitempty"`
	Type  Flag      `json:"type,omitempty"`
	Deny  string    `json:"deny,string,omitempty"`
	Allow string    `json:"allow,string,omitempty"`
}

const (
	FlagPermissionOverwriteTypeRole   = 0
	FlagPermissionOverwriteTypeMember = 1
)
