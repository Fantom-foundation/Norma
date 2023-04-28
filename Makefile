BUILD_DIR := $(CURDIR)/build

.PHONY: all test clean

all: build-opera-docker-image norma

build-opera-docker-image:
	docker build . -t opera

norma: build-opera-docker-image
	go build -o $(BUILD_DIR)/norma ./driver/norma

test:
	go test ./... -v

clean:
	rm -rvf $(CURDIR)/build