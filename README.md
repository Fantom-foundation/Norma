Norma project
=============

Project of integrating Carmen Storage and Tosca VM into the go-opera.

## Requirements

Before you run `go generate ./...`, make sure you have installed:
* GoMock: `go install github.com/golang/mock/mockgen@v1.6.0`
* Solidity Compiler (solc) - see [Installing the Solidity Compiler](https://docs.soliditylang.org/en/latest/installing-solidity.html)
* go-ethereum's abigen - you need to compile [go-ethereum](https://github.com/ethereum/go-ethereum/)
  and copy it from `build/bin/abigen` into your PATH

## Using Docker

Some experiments simulate network using Docker. For a local development the Docker must be installed:
* MacOS: https://docs.docker.com/desktop/install/mac-install/
* Linux: https://docs.docker.com/engine/install/ubuntu/

The experiments use the docker image that wraps the forked Opera/Norma client. The image is build as part of 
the build process, and can be explicitly triggered:
```
make build-docker-image
```

During the development, a few Docker commands can come handy:
```
docker run -i -t -d opera  // runs container with Opera in background (without -d it would run in foreground)
docker ps                   // shows running container
docker exec -it <ID> /bin/sh    // opens interactive shell inside the container, the ID is obtained by previous command
docker logs <ID>            // prints stdout (log) of the container
docker stop <ID>            // stop (kills) the container 
```
