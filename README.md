[![pipeline status](https://gitlab.com/banterbus/banter-bus-management-api/badges/master/pipeline.svg?style=flat-square)](https://gitlab.com/banterbus/banter-bus-management-api/commits/master)
[![coverage report](https://gitlab.com/banterbus/banter-bus-management-api/badges/master/coverage.svg?style=flat-square)](https://gitlab.com/banterbus/banter-bus-management-api/commits/master)

# Banter Bus Management API

The repo contains all the code related "REST" API for the Banter Bus application.

## Usage

```bash
git clone git@gitlab.com:banter-bus/banter-bus-management-api.git
```

### devcontainers

The preferred method is to use devcontainers, you will need:

- [VSCode](https://code.visualstudio.com/)
- [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

```bash
```

### Other

```bash
go mod download
```

#### Docker

To start the docker-compose containers:

```bash
make start
```

To just start the database, for example when debugging in VSCode. You can do:

```bash
make start-db
```

## Database Client

We are using the NoSQL database client, which provides an easy to use GUI at `localhost:3000`. It allows us to check the state of the database without needing
to use the CLI, which is very handy whilst testing/debugging.
