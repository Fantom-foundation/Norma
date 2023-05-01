Project Norma
=============

Project of integrating Carmen Storage and Tosca VM into go-opera.

# Building and Running

## Requirements

For building the project, the following tools are required:
* Go: version 1.20 or later; we recommend to use your system's package manager; alternatively, you can follow Go's [installation manual](https://go.dev/doc/install) or; if you need to maintain multiple versions, [this tutorial](https://go.dev/doc/manage-install) describes how to do so
* Docker: we recommend to use your system's package manager or the installation manuals listed in the [Using Docker](#using-docker) section below
* GNU make, or comatible

Optinally, before running `go generate ./...`, make sure you have installed:
* GoMock: `go install github.com/golang/mock/mockgen@v1.6.0`
* Solidity Compiler (solc) - see [Installing the Solidity Compiler](https://docs.soliditylang.org/en/latest/installing-solidity.html)
* go-ethereum's abigen - you need to compile [go-ethereum](https://github.com/ethereum/go-ethereum/)
  and copy it from `build/bin/abigen` into your PATH


## Building

To build the project, run
```
make -j
```
This will build the required docker images (make sure you have Docker access permissions!) and the Norma go application. To run tests, use
```
make test
```
To clean up a build, use `make clean`.

## Running

To run Norma, you can run the `norma` executable created by the build process:
```
build/norma <cmd> <args...>
```
To list the available commands, run
```
build/norma
```


# Developer Information

## Using Docker

Some experiments simulate network using Docker. For a local development the Docker must be installed:
* MacOS: https://docs.docker.com/desktop/install/mac-install/
* Linux: https://docs.docker.com/engine/install/ubuntu/

### Permissions
After installation, make sure your user has the needed permissions to run docker containers on your system. You can test this by running
```
docker images
```
If you get an error stating a lack of permissions, you might have to add your non-root user to the docker group (see [this stackoverflow post](https://stackoverflow.com/questions/48957195/how-to-fix-docker-got-permission-denied-issue) for details):
```
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker
```
If the `newgrp docker` command is not working, a `reboot` might help.


### Building
The experiments use the docker image that wraps the forked Opera/Norma client. The image is build as part of 
the build process, and can be explicitly triggered:
```
make build-docker-image
```

### Commands
During the development, a few Docker commands can come handy:
```
docker run -i -t -d opera         // runs container with Opera in background (without -d it would run in foreground)
docker ps                         // shows running container
docker exec -it <ID> /bin/sh      // opens interactive shell inside the container, the ID is obtained by previous command
docker logs <ID>                  // prints stdout (log) of the container
docker stop <ID>                  // stop (kills) the container
docker rm -f $(docker ps -a -q)   // stop and clean everything 
```

## Known Norma Restrictions

Known restrictions
 - only one node will be a validator, and it is the first node to be started; this node must life until the end of the scenario
 - currently, all transactions are send to the validator node
