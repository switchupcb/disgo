# JSON

JSON is a human-readable serialization format used to exchange data. Disgo serializes and deserializes a large amount of data due to its functionality as an API Wrapper. Disgo uses **TBA** for JSON serialization.

## Libraries

| Library                                                                               | Description                                                                        | Last Commit (as of March 20, 2022) |
| :------------------------------------------------------------------------------------ | :--------------------------------------------------------------------------------- | :--------------------------------- |
| [go-json](https://github.com/goccy/go-json)                                           | Fast JSON Serializing & Deserializing (Drop in replacement for `encoding/json`)    | 1 day                              |
| [sonic](https://github.com/bytedance/sonic)                                           | Fast JSON Serializing & Deserializing (HTML escaping NOT in conformity to RFC8259) | 20 days                            |
| [gjson](https://github.com/tidwall/gjson) + [sjson](https://github.com/tidwall/sjson) | Get and Set JSON                                                                   | 1 month                            |
| [easyjson](https://github.com/mailru/easyjson)                                        | Fast JSON Serializing & Deserializing                                              | 2021                               |
| [jsoniter](https://github.com/json-iterator/go)                                       | Fast JSON Serializing & Deserializing (Drop in replacement for `encoding/json`)    | 2021                               |
| [jsonparser](https://github.com/buger/jsonparser)                                     | Fast JSON Serializing & Deserializing (selective parser; no Schema)                | 2021                               |
| [ffjson](https://github.com/pquerna/ffjson )                                          | Fast JSON Serializing & Deserializing (Drop in replacement for `encoding/json`)    | 2019                               |


### Benchmark (Standardized)

Standardized benchmarks are standardized from raw claimed benchmarks by using ratios.

| Library | Small (<> KB) | Medium (<> KB) | Large (<> KB) |
| :------ | :------------ | :------------- | :------------ |

### Benchmark (Raw)

Raw benchmarks compare the claimed values from individual libraries _(as each likely requires a specific implementation)_.

| Library | Small (KB) | Small (ns/op) | Medium (KB) | Medium (ns/op) | Large (KB) | Large (ns/op) |
| :------ | :--------- | :------------ | :---------- | :------------- | :--------- | :------------ |

### Other Benchmarks

