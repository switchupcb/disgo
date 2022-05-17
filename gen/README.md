# Generator

Disgo uses generators to easily update and maintain over 8000 lines of code.

## Dasgo

Disgo sources Discord API objects from [dasgo](https://github.com/switchupcb/dasgo). All updates to Discord API Structs should be made in that repository.

| Step      | Description                                                                                |
| :-------- | :----------------------------------------------------------------------------------------- |
| download  | Download an updated version of `dasgo` for modification (`-d`).                            |
| endpoints | `dasgo` endpoints _(which must be removed)_ are converted into `disgo` endpoint functions. |
| nstruct   | `dasgo` structs are renamed _(via filename)_ to the `disgo` standard.                      |
| xstruct   | `dasgo` structs are extracted into one file. Uses option to include `var` and `const`.     |
| snowflake | `Snowflake` fields are converted to `string`.                                              |
| flagstd   | Flags in the extracted file are standardized.                                              |

## Disgo

Disgo generates code for features using [copygen](https://github.com/switchupcb/copygen). **This requires corresponding `setup.go` files to be updated.** Use the `diff` from Github to update those files accordingly.

| Step      | Description                                          |
| :-------- | :--------------------------------------------------- |
| `send.go` | Uses copygen to generate request `Send()` functions. |