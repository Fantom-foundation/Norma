#!/bin/bash
set -e # exit script if anything fails

# Install Carmen requirements
go install github.com/bazelbuild/bazelisk@v1.15.0
ln -s /go/bin/bazelisk /bin/bazel
apt-get update
apt-get install -y clang

# Build Carmen C++ library
cd client/carmen/go/lib
./build_libcarmen.sh

# Build Opera
cd ../../..
make opera
