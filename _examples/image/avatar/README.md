# Example: Avatar

This example adds and remove an avatar _(image)_ to your bot.

## Setup

**You must create a Discord Application in the [Discord Developer Portal](https://discord.com/developers/docs/getting-started#creating-an-app) to receive your Bot Token.** 

### Environment Variables

Assign an environment variable in the command line you will be running the program from.

#### Windows

```
set TOKEN=value
```

#### Mac/Linux

```
export TOKEN=value
``` 

**NEVER SHOW YOUR TOKEN TO THE PUBLIC.**

### Flags

This example uses [command line flags](https://pkg.go.dev/flag) to input the avatar file _(filename and location)_.

```
> avatar -h
Usage of avatar:
  -i string
        Set the location (Filepath or URL) of the avatar image using -i.
  -r    Use -r to remove the avatar image after successfully setting it.
```

## Usage

Use `go build` to build the executable binary. Use `avatar` to run it from the command line.

### Disclaimer

Image size must a size between 16 and 4096 _(2^4 - 2^12)_. For more information, view [image formatting](https://discord.com/developers/docs/reference#image-formatting).

### File

```
avatar -i avatar.png
```

_NOTE: `avatar.png` file must be in the current directory._


### URL

```
avatar -i send -c 1041179872518737990 -f https://assets-global.website-files.com/6257adef93867e50d84d30e2/636e0b52aa9e99b832574a53_full_logo_blurple_RGB.png
```