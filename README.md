<h1 align="center">Zeet CLI</h1>

<p align="center">
  <img width="150" height="150" src="./logo.svg" alt="The doctl mascot." />
</p>


## Usage
*For detailed docs, check out the [CLI Docs](https://docs.zeet.co/cli)*

```
Zeet CLI

Usage:
  zeet [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config:set  Add or modify a CLI config variable
  deploy      Deploy a project
  env:get     Retrieve an environment variable for a project
  env:set     Add or modify an environment variable for a project
  help        Help about any command
  job:run     Executes a command on a project, in a new instance
  login       Login to Zeet. You'll be prompted for a token (from https://zeet.co/account/api) if it's not passed via --token.
  logs        Logs the output for a given project
  restart     Restart a project
  status      Gets the status for a given project

Flags:
  -c, --config string   Config file (default "/Users/h/Library/Application Support/zeet/config.yaml")
  -h, --help            help for zeet

Use "zeet [command] --help" for more information about a command.
```

## Installing

#### MacOS
Using [Homebrew](https://brew.sh/):
```
brew install zeet-dev/tap/zeet
```

#### Other operating systems
Download the latest release from the [Releases](https://github.com/zeet-dev/cli/releases) page.

### Authenticating
Create an API key by going to https://zeet.co/account/api, or to [Dashboard](https://zeet.co/dashboard) >
Team Settings > API Keys > New API Key.

Then, run `zeet login` and paste your key.
