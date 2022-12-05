# Generator

Disgo uses generators to easily update and maintain over 10,000 lines of code. 

## Build

Use `go build -o gen` from the [`./_gen`](/_gen) directory to build the executable file for the generator. This may require you to set the `GOWORK` environment variable to `off`.

### Dependencies

_The following dependencies are used to download and unzip the latest version of [`dasgo`](https://github.com/switchupcb/dasgo) to `./gen/input/dasgo-10` (git-ignored)._

- [curl](https://curl.se/): [download](https://curl.se/download.html)
- [`unzip`](https://linux.die.net/man/1/unzip)

_The following dependencies are used during code generation._

- [`xstruct`](/tools/xstruct.exe)
- [copygen](https://github.com/switchupcb/copygen): `go install github.com/switchupcb/copygen@latest`

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

| Step           | Description                                                        |
| :------------- | :----------------------------------------------------------------- |
| `send.go`      | Uses copygen to generate request `Send()` functions.               |
| `handle.go`    | Uses copygen to generate request **event handling** functionality. |
| `sendevent.go` | Uses copygen to generate request `SendEvent()` functions.          |
| Clean          | Cleans the generated code.                                         |

# Bundle

## Build

A bundler is used to package the API Wrapper into the `disgo` package (`disgo.go`). Use `go build` from the [`./_gen/bundle`](/_gen/bundle) directory to build the executable file for the bundler. This may require you to set the `GOWORK` environment variable to `off`.

### Dependencies

- [bundle](https://pkg.go.dev/golang.org/x/tools/cmd/bundle): `go install golang.org/x/tools/cmd/bundle@latest`
- [fieldalignment](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment): `go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest`

### Comments

The `fieldalignment` tool is prone to removing comments from fieldaligned fields. Use the following steps to add removed comments from fields back.

1. Create a file in [`_gen/bundle/tools/replaced`](/_gen/bundle/tools/replaced/).
2. Add the fieldaligned type (with removed comments) to the file.
3. Add a boundary (`---`) on its own line.
4. Add the fieldaligned type (with comments added) to the file.

_Files must end with a single trailing newline._

Run `bundle` again to replace the fieldaligned type (without comments) in the `Replace` step.

