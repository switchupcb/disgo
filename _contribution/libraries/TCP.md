# TCP

TCP is a protocol used to exchange messages over a network. TCP is used to receive events from the Discord API's Gateway (WebSocket Connection). Disgo uses a `nhooyr/websocket` fork to manage TCP.

## Libraries

| Library                                                   | Description        | Last Commit (as of March 28, 2022) |
| :-------------------------------------------------------- | :----------------- | :--------------------------------- |
| [gorilla/websocket](https://github.com/gorilla/websocket) | WebSocket          | 1 month                            |
| [gobwas/ws](https://github.com/gobwas/ws)                 | RFC-6455 WebSocket | 8 months                           |
| [nhooyr/websocket](https://github.com/nhooyr/websocket)   | WebSocket          | 11 months                          |

### Source

| Name                        | URL                                                                                 | Date         |
| :-------------------------- | :---------------------------------------------------------------------------------- | :----------- |
| **nhooyr**, gorlla, gobwas  | https://github.com/nhooyr/websocket/commit/edda9c633d5c78c7d38fcc952b4105dd4ccfb619 | Nov 9, 2019  |
| Gorilla Websocket           | https://github.com/gorilla/websocket/pull/542#issue-496671284                       | Sep 21, 2019 |
| A Million Websockets and Go | https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/           | Aug 2, 2017  |

