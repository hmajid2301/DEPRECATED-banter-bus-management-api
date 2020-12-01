[![pipeline status](https://gitlab.com/banterbus/banter-bus-server/badges/master/pipeline.svg?style=flat-square)](https://gitlab.com/banterbus/banter-bus-server/commits/master)
[![coverage report](https://gitlab.com/banterbus/banter-bus-server/badges/master/coverage.svg?style=flat-square)](https://gitlab.com/banterbus/banter-bus-server/commits/master)

# Banter Bus Server

The repo contains all the code related to the game server for the Banter Bus application.

## Usage

Install dependencies locally:

```bash
go mod download
```

### Docker

To start the docker-compose containers:

```bash
make start
```

### Database

To just start the database, for example when debugging in VSCode. You can do:

```bash
make start-db
```

There is a MongoDB client (GUI) at `localhost:3000`
