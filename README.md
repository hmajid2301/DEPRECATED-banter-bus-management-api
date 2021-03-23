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

Before you can use devcontainers you need to run a script to generate the correct
files. This script allows you to specify things such as which shell to use bash/fish
and if you want to add any extra fish extensions installed using fisher.

```bash
make devcontainer OPTIONS="--help"
usage:
Create devcontainer files

Options:
  -fisher value
        List of fisher extensions to install.
  -shell string
        The shell to use within the devcontainer. [fish, bash]
```

```bash
make devcontainer OPTIONS="-shell fish -fisher dracula/fish"
```

This generates some of the key files we need to use devcontainers including:

- devcontainer.json
- docker-compose.dev.yml

This allows different users customise their own dev environment and use what they prefer.

Finally, use the command palette (CTRL + SHIFT + P) and run `Remote Containers: Rebuild and Reopen in Container`.

> Remote the first time you do this it can take some time. As VSCode needs to build and run the Docker image/container.

The main reason to use devcontainers is that it means that the development environemnt of all the developers
should be very similar. It will have all the dependencies installed. It will also create the other Docker containers
related to the database when the devcontainer is created.

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
