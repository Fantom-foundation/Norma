# syntax=docker/dockerfile:1

# This is a multi stages Dockerfile, which builds go-opera
# from the client/ directory first, and runs the binary then.
#
# This Dockerfile requires running installation of Docker,
# and then the image is build by typing
# > docker build . -t <image-name>
#

# The build is done in independent stages, to allow for
# caching of the intermediate results.

# Stage 0: load dependencies
FROM golang:1.22 AS client-dependencies

WORKDIR /client
COPY client/go.mod ./go.mod
RUN go mod download

WORKDIR /
COPY go.mod go.mod
RUN go mod download

# Stage 1: build the client
FROM client-dependencies AS client-build

WORKDIR /client

# Copy the client code into the image.
COPY client/ ./

# Build sonic with caching
RUN --mount=type=cache,target=/root/.cache/go-build make sonicd sonictool

# Build norma itself
WORKDIR /norma
COPY . ./
RUN --mount=type=cache,target=/root/.cache/go-build make normatool

# This results in an image that contains the sonic binary
#
# The container can be run by typing:
# > docker run -i -t sonic
# or inspected opening terminal in it by
# > docker run -i -t sonic /bin/sh
#
# sonic run can be customised by passing the environment variables, .e.g.:
#
# > docker run -e VALIDATOR_NUMBER=2 -e VALIDATORS_COUNT=5 -i -t sonic
#
FROM debian:bookworm

RUN apt-get update && \
    apt-get install iproute2 iputils-ping -y

COPY --from=client-build /client/build/sonicd /client/build/sonictool ./
COPY --from=client-build /norma/build/normatool ./

ENV STATE_DB_IMPL="geth"
ENV VM_IMPL="geth"
ENV LD_LIBRARY_PATH=./
ENV TINI_KILL_PROCESS_GROUP=1

EXPOSE 5050
EXPOSE 6060
EXPOSE 18545
EXPOSE 18546

COPY genesis/example-genesis.json ./genesis.json
COPY scripts/run_sonic_privatenet.sh ./run_sonic.sh
COPY scripts/set_genesis.sh ./set_genesis.sh

CMD ["./run_sonic.sh"]
