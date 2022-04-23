package responses

// Modify Current User Nick Response
// https://discord.com/developers/docs/topics/gateway#get-gateway-example-response
type ModifyCurrentUserNick struct {
	Nick *string `json:"nick,omitempty"`
}
