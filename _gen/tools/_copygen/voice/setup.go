package wrapper

import (
	disgo "github.com/switchupcb/disgo/wrapper"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	// opcode 0
	Identify(*disgo.VoiceIdentify) error
	// opcode 1
	SelectProtocol(*disgo.SelectProtocol) error
	// opcode 3
	Heartbeat(*disgo.VoiceHeartbeat) error
	// opcode 5
	Speaking(*disgo.Speaking) error
	// opcode 7
	Resume(*disgo.VoiceResume) error
}
