package resources

// Team Object
// https://discord.com/developers/docs/topics/teams#data-models-team-object
type Team struct {
	Icon        string        `json:"icon"`
	ID          string        `json:"id"`
	Members     []*TeamMember `json:"members"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	OwnerUserID int64         `json:"owner_user_id"`
}

// Team Member Object
// https://discord.com/developers/docs/topics/teams#data-models-team-member-object
type TeamMember struct {
	MembershipState uint8    `json:"membership_state"`
	Permissions     []string `json:"permissions"`
	TeamID          string   `json:"team_id"`
	User            *User    `json:"user"`
}

// Membership State Enum
// https://discord.com/developers/docs/topics/teams#data-models-membership-state-enum
const (
	INVITED  = 1
	ACCEPTED = 2
)
