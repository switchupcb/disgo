# JSON

JSON is a human-readable serialization format used to exchange data. Disgo serializes and deserializes a large amount of data due to its functionality as an API Wrapper. Disgo uses `go-json` for JSON serialization.

## Libraries

| Library                                             | Description                          | Last Commit (as of March 20, 2022) |
| :-------------------------------------------------- | :----------------------------------- | :--------------------------------- |
| [*go-json](https://github.com/goccy/go-json)        | Serialize & Deserialize              | 1 day                              |
| [simdjson-go](https://github.com/minio/simdjson-go) | SIMD JSON Serializer & Deserializer  | 10 days                            |
| [sonic](https://github.com/bytedance/sonic)         | JIT, SIMD Serializer & Deserializer  | 20 days                            |
| [gjson/sjson](https://github.com/tidwall/gjson)     | Get & Set JSON                       | 1 month                            |
| [*jsoniter](https://github.com/json-iterator/go)    | Serializing & Deserializing          | 2021                               |
| [jsonparser](https://github.com/buger/jsonparser)   | Deserializing (selective; no schema) | 2021                               |
| [easyjson](https://github.com/mailru/easyjson)      | Serializing & Deserializing          | 2021                               |

_* indicates a library that functions as a drop-in replacement for `encoding/json`._

`sonic` claims to have a higher encode, decode, and get-one speed than all libraries except for `simdjson-go`; which is based upon the [C++ simdjson library](https://github.com/simdjson/simdjson). `simdjson-go` places a [requirement on the CPU](https://github.com/minio/simdjson-go#requirements) for parsing (non-serialized data). `sonic` places a [requirement on the CPU and OS](https://github.com/bytedance/sonic#requirement) for use. `go-json` is the slowest of the three, but functions as a drop-in replacement for the standard library. **This library aims to function without requirements, so `go-json` is used**. Developers that require faster serialization are advanced enough to implement faster serialization libraries on their own.

### Benchmark (Raw)

Raw benchmarks compare the claimed values from individual libraries _(as each likely requires a specific implementation)_.

| Library                | Small (KB) | Small (ns/op) | Medium (KB) | Medium (ns/op)      | Large (KB)   | Large (ns/op) |
| :--------------------- | :--------- | :------------ | :---------- | :------------------ | :----------- | :------------ |
| easyjson               | .08        | 67 MB/s       | ?           | ?                   | 500          | 125 MB/s      |
| gjson                  | ?          | ?             | .319        | 311                 | ?            | ?             |
| jsoniter               | ?          | ?             | ?           | ?                   | 24.9K        | 37760014      |
| jsonparser             | .19        | 1367          | 2.4         | 15955               | 24           | 85308         |
| simdjson-go _(decode)_ | ?          | ?             | ?           | ?                   | 31K Compress | 1242 MB/s     |
| go-json  _(decode)_    | .123       | 427           | 1.52        | 3054                | 25.3         | 44780         |
| sonic _(decode)_       | .4         | 256 MB/s      | 13          | 42688 (305.36 MB/s) | 635          | 594 MB/s      |

_*`go-json` [claims](https://github.com/goccy/go-json/pull/254) to have a (near-equivalent) encode speed, and 2x decode speed of sonic. In contrast, `sonic` [claims (6 months later)](https://github.com/bytedance/sonic#benchmarks) to have a (near equivalent; higher) decode speed, 1.7x (Generic) - 3x (Parallel) encode MB/s speed of go-json. `sonic` is owned by ByteDance's TikTok._

### Source

| Compare                                          | URL                                            | Date              |
| :----------------------------------------------- | :--------------------------------------------- | :---------------- |
| sonic, gojson, gjson, jsoniter, _(with MB/s)_    | https://github.com/bytedance/sonic#benchmarks  | Jan 25, 2022      |
| gojson, sonic                                    | https://github.com/goccy/go-json/pull/254      | Jun 23, 2021      |
| simdjson-go, jsoniter                            | https://github.com/minio/simdjson-go           | ***May 17, 2021** |
| _All Previous (decode)_                          | https://github.com/buger/jsonparser#benchmarks | Jun 20, 2021      |
| gjson, easyjson, jsonparser, jsoniter _(decode)_ | https://github.com/tidwall/gjson#performance   | Apr 10, 2017      |
| jsonparser vs. jsoniter vs. easyjson             | https://github.com/json-iterator/go-benchmark  | Dec 7, 2016       |
| easyjson vs. ffjson                              | https://github.com/mailru/easyjson#benchmarks  | Feb 28, 2016      |
