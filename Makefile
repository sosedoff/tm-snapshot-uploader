PACKAGE ?= github.com/sosedoff/tm-snapshot-uploader
DOCKER_IMAGE ?= tm-snapshot-uploader

.PHONY: build
build:
	go build -o ./build/tm-snapshot-uploader

.PHONY: docker
docker:
	docker build --build-arg=PACKAGE=$(PACKAGE) -t $(DOCKER_IMAGE) .
