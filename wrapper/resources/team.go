package resources

// Team Object
// https://discord.com/developers/docs/topics/teams#data-models-team-object
type Team struct {
	Icon        string        `json:"icon,omitempty"`
	ID          Snowflake     `json:"id,omitempty"`
	Members     []*TeamMember `json:"members,omitempty"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	OwnerUserID Snowflake     `json:"owner_user_id,omitempty"`
}

// Team Member Object
// https://discord.com/developers/docs/topics/teams#data-models-team-member-object
type TeamMember struct {
	MembershipState Flag      `json:"membership_state,omitempty"`
	Permissions     []string  `json:"permissions,omitempty"`
	TeamID          Snowflake `json:"team_id,omitempty"`
	User            *User     `json:"user,omitempty"`
}

// Membership State Enum
// https://discord.com/developers/docs/topics/teams#data-models-membership-state-enum
const (
	FlagEnumStateMembershipINVITED  = 1
	FlagEnumStateMembershipACCEPTED = 2
)
