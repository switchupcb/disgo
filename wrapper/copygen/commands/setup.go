package wrapper

import (
	disgo "github.com/switchupcb/disgo/wrapper"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	// opcode 1
	Heartbeat(*disgo.Heartbeat) error
	// opcode 2
	Identify(*disgo.Identify) error
	// opcode 3
	PresenceUpdate(*disgo.GatewayPresenceUpdate) error
	// opcode 4
	VoiceStateUpdate(*disgo.VoiceStateUpdate) error
	// opcode 6
	Resume(*disgo.Resume) error
	// opcode 8
	RequestGuildMembers(*disgo.RequestGuildMembers) error
}
