// Package requests contains the setup information for copygen generated code.
package requests

import (
	"github.com/switchupcb/disgo/wrapper/requests"
	"github.com/switchupcb/disgo/wrapper/resources"
)

// Copygen defines the functions that will be generated.
type Copygen interface {
	SendGetGlobalApplicationCommands(*requests.GetGlobalApplicationCommands) ([]*resources.ApplicationCommand, error)
	SendCreateGlobalApplicationCommand(*requests.CreateGlobalApplicationCommand) (*resources.ApplicationCommand, error)
	SendGetGlobalApplicationCommand(*requests.GetGlobalApplicationCommand) (*resources.ApplicationCommand, error)
	SendEditGlobalApplicationCommand(*requests.EditGlobalApplicationCommand) (*resources.ApplicationCommand, error)
}
