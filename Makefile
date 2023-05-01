BUILD_DIR := $(CURDIR)/build

.PHONY: all test clean

all: build-opera-docker-image norma

pull-hello-world-image:
	docker image pull hello-world

build-opera-docker-image:
	docker build . -t opera

norma: build-opera-docker-image
	go build -o $(BUILD_DIR)/norma ./driver/norma

test: pull-hello-world-image build-opera-docker-image
	go test ./... -v

clean:
	rm -rvf $(CURDIR)/build
