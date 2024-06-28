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

# Stage 1: build the client
FROM golang:1.22 AS client-build

WORKDIR /client

# Copy the client code into the image.
COPY client/ ./

# Build sonic with caching
RUN --mount=type=cache,target=/root/.cache/go-build make sonicd sonictool

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
COPY --from=client-build /client/build/sonicd .
COPY --from=client-build /client/build/sonictool .

ENV STATE_DB_IMPL="geth"
ENV VM_IMPL="geth"
ENV LD_LIBRARY_PATH=./

EXPOSE 5050
EXPOSE 6060
EXPOSE 18545
EXPOSE 18546

COPY genesis/example-genesis.json ./genesis.json
COPY scripts/run_sonic_privatenet.sh ./run_sonic.sh
CMD ["/bin/bash", "run_sonic.sh"]
