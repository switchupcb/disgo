# HTTP

HTTP is a protocol used to communicate resources. HTTP is used to create and receive thousands of requests per second from the Discord API's application commands. Disgo uses `fiber` to manage HTTP.
 

## Libraries

| Library                                         | Description                          | Last Commit (as of March 28, 2022) |
| :---------------------------------------------- | :----------------------------------- | :--------------------------------- |
| [gnet](https://github.com/panjf2000/gnet)       | Non-blocking Event-Driven Networking | 2 days                             |
| [fasthttp](https://github.com/valyala/fasthttp) | Zero-allocation-hotpath HTTP         | 8 days                             |
| [echo](https://github.com/labstack/echo)        | `net/http` Web Framework             | 20 days                            |
| [fiber](https://github.com/gofiber/fiber)       | `fasthttp` Web Framework             | ***5 days**                        |
| [gearbox](https://github.com/gogearbox/gearbox) | `fasthttp` Web Framework             | 4 months                           |

In modern Go, performance differences between HTTP frameworks depend on the underlying HTTP client: `fasthttp` is up to 10x faster than `net/http`, so any web framework that uses it (or allows you to swap from `net/http`) will maintain a higher performance. While `gnet` is the [6th fastest of all languages](https://www.techempower.com/benchmarks/#section=data-r20&hw=ph&test=plaintext&s=1) _(faster than `fasthttp` for plaintext)_, using event-driven code (as opposed to Go's goroutine features) is not worth it.

`echo` [dropped support for `fasthttp` in v3](https://github.com/labstack/echo/issues/1617#issuecomment-771379610) which no longer makes it one of the fastest web frameworks for Go. `gearbox` is the [6th fastest framework of all languages](https://web-frameworks-benchmark.netlify.app/result), but it's [extremely lightweight](https://gogearbox.com/). `fiber` is not much slower than `gearbox`, more maintained, and contains [many features](https://docs.gofiber.io/) useful for serving HTTP requests.

### Source

| Name                          | URL                                                                                           | Date           |
| :---------------------------- | :-------------------------------------------------------------------------------------------- | :------------- |
| Web Frameworks Benchmark (Go) | https://web-frameworks-benchmark.netlify.app/result?l=go                                      | March 24, 2022 |
| Best JSON response/s          | https://www.techempower.com/benchmarks/#section=data-r20&hw=ph&test=json&s=1&l=zijocf-sf      | Feb 8, 2021    |
| Best Plaintext response/s     | https://www.techempower.com/benchmarks/#section=data-r20&hw=ph&test=plaintext&s=1&l=zijocf-sf | Feb 8, 2021    |

