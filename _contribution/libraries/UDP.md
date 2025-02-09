# UDP

UDP is a protocol used to establish low-latency connections between multiple parties. UDP is used to transmit voice data over a Discord UDP Connection. Disgo uses the standard library `net`] to manage UDP traffic.

## Libraries

| Library                                   | Description                                     |
| :---------------------------------------- | :---------------------------------------------- |
| [net](https://pkg.go.dev/net)             | Portable interface for network I/O              |
| [gnet](https://github.com/panjf2000/gnet) | Non-blocking, event-driven networking framework |
