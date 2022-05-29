# URL

URL encoded query strings are involved in a number of Discord API HTTP requests. Disgo uses `gorilla/schema` to encode and decode URL Query String (Parameters). This may change to a custom library or custom functions in the future to avoid the use of reflection.