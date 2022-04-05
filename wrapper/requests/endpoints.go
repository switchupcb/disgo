package requests

// Discord API Endpoints
// defined using a parser, then appending the Base URL.
const (
	EndpointGetGlobalApplicationCommands   = "https://discord.com/api/v9/applications/{application.id}/commands"
	EndpointCreateGlobalApplicationCommand = "https://discord.com/api/v9/applications/{application.id}/commands"
	EndpointGetGlobalApplicationCommand    = "https://discord.com/api/v9/applications/{application.id}/commands/{command.id}"
	EndpointEditGlobalApplicationCommand   = "https://discord.com/api/v9/applications/{application.id}/commands/{command.id}"
)
