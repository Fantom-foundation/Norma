# Copyright 2024 Fantom Foundation
# This file is part of Norma System Testing Infrastructure for Sonic.
#
# Norma is free software: you can redistribute it and/or modify
# it under the terms of the GNU Lesser General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Norma is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU lesser General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public License
# along with Norma. If not, see <http://www.gnu.org/licenses/>.

BUILD_DIR := $(CURDIR)/build

.PHONY: all test clean

all: build-sonic-docker-image norma normatool

pull-hello-world-image:
	DOCKER_BUILDKIT=1 docker image pull hello-world

pull-alpine-image:
	DOCKER_BUILDKIT=1 docker image pull alpine

pull-prometheus-image:
	DOCKER_BUILDKIT=1 docker image pull prom/prometheus:v2.44.0

build-sonic-docker-image:
	DOCKER_BUILDKIT=1 docker build . -t sonic

generate-abi: load/contracts/abi/Counter.abi load/contracts/abi/ERC20.abi load/contracts/abi/Store.abi load/contracts/abi/UniswapV2Pair.abi load/contracts/abi/UniswapRouter.abi # requires installed solc and Ethereum abigen - check README.md

load/contracts/abi/Counter.abi: load/contracts/Counter.sol
	solc --evm-version london -o ./load/contracts/abi --overwrite --pretty-json --optimize --optimize-runs 200 --abi --bin ./load/contracts/Counter.sol
	abigen --type Counter --pkg abi --abi load/contracts/abi/Counter.abi --bin load/contracts/abi/Counter.bin --out load/contracts/abi/Counter.go

load/contracts/abi/ERC20.abi: load/contracts/ERC20.sol
	solc --evm-version london -o ./load/contracts/abi --overwrite --pretty-json --optimize --optimize-runs 200 --abi --bin ./load/contracts/ERC20.sol
	abigen --type ERC20 --pkg abi --abi load/contracts/abi/ERC20.abi --bin load/contracts/abi/ERC20.bin --out load/contracts/abi/ERC20.go

load/contracts/abi/Store.abi: load/contracts/Store.sol
	solc --evm-version london -o ./load/contracts/abi --overwrite --pretty-json --optimize --optimize-runs 200 --abi --bin ./load/contracts/Store.sol
	abigen --type Store --pkg abi --abi load/contracts/abi/Store.abi --bin load/contracts/abi/Store.bin --out load/contracts/abi/Store.go

load/contracts/abi/UniswapV2Pair.abi: load/contracts/UniswapV2Pair.sol
	solc --evm-version london -o ./load/contracts/abi --overwrite --pretty-json --optimize --optimize-runs 200 --abi --bin ./load/contracts/UniswapV2Pair.sol
	abigen --type UniswapV2Pair --pkg abi --abi load/contracts/abi/UniswapV2Pair.abi --bin load/contracts/abi/UniswapV2Pair.bin --out load/contracts/abi/UniswapV2Pair.go

load/contracts/abi/UniswapRouter.abi: load/contracts/UniswapRouter.sol
	solc --evm-version london -o ./load/contracts/abi --overwrite --pretty-json --optimize --optimize-runs 200 --abi --bin ./load/contracts/UniswapRouter.sol
	abigen --type UniswapRouter --pkg abi --abi load/contracts/abi/UniswapRouter.abi --bin load/contracts/abi/UniswapRouter.bin --out load/contracts/abi/UniswapRouter.go

generate-mocks: # requires installed mockgen
	go generate ./...

norma: pull-prometheus-image build-sonic-docker-image
	go build -o $(BUILD_DIR)/norma ./driver/norma

normatool: 
	go build -o $(BUILD_DIR)/normatool ./driver/normatool

test: pull-hello-world-image pull-alpine-image pull-prometheus-image build-sonic-docker-image
	go test ./... -v

clean:
	rm -rvf $(CURDIR)/build
