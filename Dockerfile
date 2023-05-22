# This is a two stages Dockerfile, which builds go-opera
# from the client/ directory first, and runs the binary then. 
#
# This Dockerfile requires running installation of Docker, 
# and then the image is build by typing
# > docker build . -t <image-name>
#
FROM golang:1.20.3 AS build

# Install Carmen prerequisities
RUN apt-get update && apt-get install -y clang
RUN go install github.com/bazelbuild/bazelisk@v1.15.0 && ln -s /go/bin/bazelisk /bin/bazel

COPY client/ /client

# Build Carmen C++ library
WORKDIR /client/carmen/go/lib
RUN /bin/bash ./build_libcarmen.sh

# Build Opera
WORKDIR /client
RUN make opera

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
COPY --from=build /client/build/opera .
COPY --from=build /client/carmen/go/lib/libcarmen.so .

ENV VALIDATOR_NUMBER=1
ENV VALIDATORS_COUNT=1
ENV LD_LIBRARY_PATH=./

EXPOSE 6060
EXPOSE 18545

COPY scripts/run_opera.sh .
CMD /bin/bash run_opera.sh
