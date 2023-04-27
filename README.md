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
docker run -i -t -d opera         // runs container with Opera in background (without -d it would run in foreground)
docker ps                         // shows running container
docker exec -it <ID> /bin/sh      // opens interactive shell inside the container, the ID is obtained by previous command
docker logs <ID>                  // prints stdout (log) of the container
docker stop <ID>                  // stop (kills) the container
docker rm -f $(docker ps -a -q)   // stop and clean everything 
```

## Restrictons

Known restrictions
 - only one node will be a validator, and it is the first node to be started; this node must life until the end of the scenario
 - currently, all transactions are send to the validator node


## TODOs

Docker
- replace busy waiting for opera node to be online by probing (e.g, try to connect to RPC server)
- cleanup shutdown procedure; make sure opera node is really disconnecting from the network
- fix build for MacOS

Driver
- make executor spawn events in parallel
- collect all logs from the nodes
- remove nodes if driver crashes

Generator
- support slop-shaped traffic
- support wave-shaped traffic
- support >1000 Tx/s

Monitoring
- live monitoring of block processing in CLI tool
- get throughput metric (transactions/block and block time)
- get system metrics:
   - CPU load
   - Memory usage
   - Storage usage
- get Opera-internal metrics:
   - block finish / next block start delta time
   - block processing time
   - block commitment time
   - other processing phases

Client adaptations
- support Carment DB
- support automated transaction rate control

Long-term features
 - run nodes on different physical hosts
 - 
