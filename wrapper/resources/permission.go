package resources

// Bitwise Permission Flags
// https://discord.com/developers/docs/topics/permissions#permissions-bitwise-permission-flags
const (
	FlagFlagsPermissionBitwiseCREATE_INSTANT_INVITE      = 1 << 0
	FlagFlagsPermissionBitwiseKICK_MEMBERS               = 1 << 1
	FlagFlagsPermissionBitwiseBAN_MEMBERS                = 1 << 2
	FlagFlagsPermissionBitwiseADMINISTRATOR              = 1 << 3
	FlagFlagsPermissionBitwiseMANAGE_CHANNELS            = 1 << 4
	FlagFlagsPermissionBitwiseMANAGE_GUILD               = 1 << 5
	FlagFlagsPermissionBitwiseADD_REACTIONS              = 1 << 6
	FlagFlagsPermissionBitwiseVIEW_AUDIT_LOG             = 1 << 7
	FlagFlagsPermissionBitwisePRIORITY_SPEAKER           = 1 << 8
	FlagFlagsPermissionBitwiseSTREAM                     = 1 << 9
	FlagFlagsPermissionBitwiseVIEW_CHANNEL               = 1 << 10
	FlagFlagsPermissionBitwiseSEND_MESSAGES              = 1 << 11
	FlagFlagsPermissionBitwiseSEND_TTS_MESSAGES          = 1 << 12
	FlagFlagsPermissionBitwiseMANAGE_MESSAGES            = 1 << 13
	FlagFlagsPermissionBitwiseEMBED_LINKS                = 1 << 14
	FlagFlagsPermissionBitwiseATTACH_FILES               = 1 << 15
	FlagFlagsPermissionBitwiseREAD_MESSAGE_HISTORY       = 1 << 16
	FlagFlagsPermissionBitwiseMENTION_EVERYONE           = 1 << 17
	FlagFlagsPermissionBitwiseUSE_EXTERNAL_EMOJIS        = 1 << 18
	FlagFlagsPermissionBitwiseVIEW_GUILD_INSIGHTS        = 1 << 19
	FlagFlagsPermissionBitwiseCONNECT                    = 1 << 20
	FlagFlagsPermissionBitwiseSPEAK                      = 1 << 21
	FlagFlagsPermissionBitwiseMUTE_MEMBERS               = 1 << 22
	FlagFlagsPermissionBitwiseDEAFEN_MEMBERS             = 1 << 23
	FlagFlagsPermissionBitwiseMOVE_MEMBERS               = 1 << 24
	FlagFlagsPermissionBitwiseUSE_VAD                    = 1 << 25
	FlagFlagsPermissionBitwiseCHANGE_NICKNAME            = 1 << 26
	FlagFlagsPermissionBitwiseMANAGE_NICKNAMES           = 1 << 27
	FlagFlagsPermissionBitwiseMANAGE_ROLES               = 1 << 28
	FlagFlagsPermissionBitwiseMANAGE_WEBHOOKS            = 1 << 29
	FlagFlagsPermissionBitwiseMANAGE_EMOJIS_AND_STICKERS = 1 << 30
	FlagFlagsPermissionBitwiseUSE_APPLICATION_COMMANDS   = 1 << 31
	FlagFlagsPermissionBitwiseREQUEST_TO_SPEAK           = 1 << 32
	FlagFlagsPermissionBitwiseMANAGE_EVENTS              = 1 << 33
	FlagFlagsPermissionBitwiseMANAGE_THREADS             = 1 << 34
	FlagFlagsPermissionBitwiseCREATE_PUBLIC_THREADS      = 1 << 35
	FlagFlagsPermissionBitwiseCREATE_PRIVATE_THREADS     = 1 << 36
	FlagFlagsPermissionBitwiseUSE_EXTERNAL_STICKERS      = 1 << 37
	FlagFlagsPermissionBitwiseSEND_MESSAGES_IN_THREADS   = 1 << 38
	FlagFlagsPermissionBitwiseUSE_EMBEDDED_ACTIVITIES    = 1 << 39
	FlagFlagsPermissionBitwiseMODERATE_MEMBERS           = 1 << 40
)

// Overwrite Object
// https://discord.com/developers/docs/resources/channel#overwrite-object
type PermissionOverwrite struct {
	ID    Snowflake `json:"id,omitempty"`
	Type  *Flag     `json:"type,omitempty"`
	Deny  string    `json:"deny,omitempty"`
	Allow string    `json:"allow,omitempty"`
}

const (
	FlagPermissionOverwriteTypeRole   = 0
	FlagPermissionOverwriteTypeMember = 1
)
