package requests

import "github.com/switchupcb/disgo/wrapper/resources"

// Get Guild Audit Log
// GET /guilds/{guild.id}/audit-logs
// https://discord.com/developers/docs/resources/audit-log#get-guild-audit-log
type GetGuildAuditLog struct {
	UserID     resources.Snowflake `json:"user_id"`
	ActionType resources.Flag      `json:"action_type"`
	Before     resources.Snowflake `json:"before,omitempty"`
	Limit      resources.Flag      `json:"limit,omitempty"`
}
