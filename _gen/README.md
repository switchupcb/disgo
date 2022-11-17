# Generator

Disgo uses generators to easily update and maintain over 10,000 lines of code. 

## Build

Use `go build -o gen` to build the executable file for the generator. This may require you to set the `GOWORK` environment variable to `off`.

## Dasgo

Disgo sources Discord API objects from [dasgo](https://github.com/switchupcb/dasgo). All updates to Discord API Structs should be made in that repository.

| Step      | Description                                                                                |
| :-------- | :----------------------------------------------------------------------------------------- |
| download  | Download an updated version of `dasgo` for modification (`-d`).                            |
| endpoints | `dasgo` endpoints _(which must be removed)_ are converted into `disgo` endpoint functions. |
| xstruct   | `dasgo` structs are extracted into one file. Uses option to include `var` and `const`.     |
| typefix   | `Snowflake`, `Nonce`, and `Value` fields are converted to `string`.                        |

## Disgo

Disgo generates code for features using [copygen](https://github.com/switchupcb/copygen). **This requires corresponding `setup.go` files to be updated.** Use the `diff` from Git to update those files accordingly.

| Step         | Description                                                        |
| :----------- | :----------------------------------------------------------------- |
| `send.go`    | Uses copygen to generate request `Send()` functions.               |
| `handle.go`  | Uses copygen to generate request **event handling** functionality. |
| `command.go` | Uses copygen to generate request `Command()` functions.            |
| Clean        | Cleans the generated code.                                         |