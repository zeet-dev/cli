<h1 align="center">Zeet CLI</h1>

<p align="center">
  <img width="150" height="150" src="./logo.svg" alt="The doctl mascot." />
</p>


## Usage
*For detailed docs, check out the [CLI Docs](https://zeet.co/cli/getting-started/)*

```
Zeet CLI

Usage:
  zeet [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  deploy      Deploy
  help        Help about any command
  login       Login to Zeet
  logs        View the logs for a project

Flags:
  -c, --config string      Config file (default "/Users/h/Library/Application Support/zeet/config.yaml")
  -v, --debug              Enable verbose debug logging
  -h, --help               help for zeet
  -s, --server string      Zeet API Server (default "https://anchor.zeet.co")
      --ws-server string   Zeet Websocket/Subscriptions Server (default "wss://anchor.zeet.co")

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
Create an API key by going to [Dashboard](https://zeet.co/dashboard) >
Team Settings > API Keys > New API Key.

Then, run `zeet login` and paste your key.