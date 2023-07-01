# This is a multi stages Dockerfile, which builds go-opera
# from the client/ directory first, and runs the binary then.
#
# This Dockerfile requires running installation of Docker,
# and then the image is build by typing
# > docker build . -t <image-name>
#

# The build is done in independent stages, to allow for
# caching of the intermediate results.

# Stage 1: build the Carmen C++ library
FROM golang:1.20.3 AS carmen-build
# Install Carmen prerequisities
RUN apt-get update && apt-get install -y clang
RUN go install github.com/bazelbuild/bazelisk@v1.15.0 && ln -s /go/bin/bazelisk /bin/bazel

COPY client/carmen/ /client/carmen

# Build Carmen C++ library
WORKDIR /client/carmen/go/lib
RUN /bin/bash ./build_libcarmen.sh

# Stage 2: build the client
FROM golang:1.20.3 AS client-build
WORKDIR /client

# Copy go.mod and go.sum files to cache dependencies
COPY client/go.* ./

# We also need to copy carmen and tosca directories, as they are
# replaced in the go.mod file with local paths
COPY client/carmen/ ./carmen
COPY client/tosca/ ./tosca

# Download dependencies
RUN go mod download

# Copy the source code
COPY client/ ./

# The built carmen library is needed to build the client
COPY --from=carmen-build /client/carmen/go/lib/libcarmen.so ./carmen/go/lib/libcarmen.so

# Build Opera with caching
RUN --mount=type=cache,target=/root/.cache/go-build make opera

# This results in an image that contains the Opera binary
#
# The container can be run by typing:
# > docker run -i -t opera
# or inspected opening terminal in it by
# > docker run -i -t opera /bin/sh
#
# Opera run can be customised by passing the environment variables, .e.g.:
#
# > docker run -e VALIDATOR_NUMBER=2 -e VALIDATORS_COUNT=5 -i -t opera
#
FROM debian:bookworm
COPY --from=client-build /client/build/opera .
COPY --from=carmen-build /client/carmen/go/lib/libcarmen.so .

ENV VALIDATOR_NUMBER=1
ENV VALIDATORS_COUNT=1
ENV STATE_DB_IMPL="geth"
ENV LD_LIBRARY_PATH=./

EXPOSE 6060
EXPOSE 18545
EXPOSE 18546

COPY scripts/run_opera.sh .
CMD /bin/bash run_opera.sh
