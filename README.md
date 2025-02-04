# Kabisa Assignment Backend - Quote API

This application wraps an api from [dummyjson.com](https://dummyjson.com/quotes) with endpoints for the following features:

- Retrieve a random quote
- Play a guessing game

The definition of the api can be found in openapi.yaml.

## Guessing game

To start a quessing game, you first need to post an empty request to `/quote-game`. You will receive three quotes and three authors, both sorted alphabetically. The goal of the game is to match which author wrote which quote. This response needs to be send within a minute to `/quote-game/{id}/answer`. For the exact JSON objects needed for this game, please refer to openapi.yaml.

## How to run

Easiest is to download the executable from the [releases](https://github.com/pietdevries94/Kabisa/releases) page and run the executable. By default the application doesn't produce any extra files, so it doesn't matter where you run it from.

You can persist the sqlite database and logs by specifying them in the configuration.

## Configuration

The application can be configured using environment variables, but for the sake of usability .env files are also supported.

| Environment variable            | Description                                                                                                                                                         | Default                      | Example                       |
| ------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------- | ----------------------------- |
| KABISAQUOTE_LISTEN_ADDRESS      | The address the application listens on                                                                                                                              | `127.0.0.1:3333`             | `:8080`                       |
| KABISAQUOTE_HTTP_CLIENT_TIMEOUT | The timeout for the HTTP client (used to fetch quotes)                                                                                                              | `10`                         | `60`                          |
| KABISAQUOTE_LOG_LEVEL           | The log level for the application                                                                                                                                   | `info`                       | `debug`                       |
| KABISAQUOTE_LOG_FILE_PATH       | The path to the log file. An empty string disables logging to a file                                                                                                | ``                           | `default.log`                 |
| KABISAQUOTE_SQLITE_DSN          | The DSN for the SQLite database, by default it's in memory. It's highly recommeded to use `?cache=shared` to prevent database locking issues with parallel requests | `file::memory:?cache=shared` | `file:quotes.db?cache=shared` |

## How to build

To build this application, you need to have one of the following two installed:

- Go 1.23.5+
- Nix package manager with flakes enabled (see below)
- make (needs to be installed on Windows. Easiest with `choco install make`)

Open a terminal in this repository and run:

```bash
make
```

This generates code (if needed), lints the codebase, tests the codebase and builds the application in the path `bin/api` (`bin\api.exe` on Windows) for your current OS and architecture.

There are targets to build cross-compile the application. The following are added as they are common combinations, but it's trivial to extend this if needed.

- `make build-linux-amd64` > `bin/api-linux-amd64`
- `make build-windows-amd64` > `bin/api-windows-amd64.exe`
- `make build-darwin-amd64` > `bin/api-darwin-amd64`
- `make build-darwin-arm64` > `bin/api-darwin-arm64`

Finally, there is also a target that builds all the targets above: `make build-all`

### Nix

[Nix](https://nixos.org/) is a package-manager that creates reproducable builds and is declarative. This repository has a valid Nix [flake](https://wiki.nixos.org/wiki/Flakes) to start a development environment with a pinned version of the underlying package repository, creating a reproduceable development environment.

However, in case of this project, the only dependency managed by Nix is Golang, so if you're not already using Nix, I recommend manually installing Golang.
