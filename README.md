# marvel-characters-api

This API serves as a gateway for fetching character data from Marvel's API.

## Prerequisites

Inspired by [3 Musketeers](https://3musketeers.io/), this project uses the following heavily for application development:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/)

There is a make target (`make build`) to build a binary which can be targeted for a particular execution environment. Environment variables are used for configuration.

Go to [Quick Start](#quick-start)!

## Usage

#### configure
```bash
$ make .env
```

* see generated `.env` file for configuration

#### run unit tests
```bash
$ make test
```

#### tidy dependencies
```bash
$ make deps
```

#### build binary (target execution environment can be specified using `GOOS` in `.env`)
```bash
$ make build
```

#### start the API inside a Docker container
```bash
$ make start
```

### Helpers during development:

#### format all .go files in project (using go fmt)
```bash
$ make fmt
```

#### generate test mocks for all interfaces in project
```bash
$ make genMocks
```

#### generate Swagger documentation
```bash
$ make genSwagger
```

### Quick Start
```bash
$ make .env
# set correct values for the following in .env
# MARVEL_API_BASE_URL
# MARVEL_API_KEY_PUBLIC
# MARVEL_API_KEY_PRIVATE
$ make start
# accessible endpoints:
# http://localhost:8080/characters
# http://localhost:8080/characters/{id}
```
