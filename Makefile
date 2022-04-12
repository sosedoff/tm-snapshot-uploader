PACKAGE ?= github.com/sosedoff/tm-snapshot-uploader

.PHONY: build
build:
	go build -o ./build/tm-snapshot-uploader

.PHONY: docker
docker:
	docker build --build-arg=PACKAGE=$(PACKAGE) -t tm-snapshot-uploader .
