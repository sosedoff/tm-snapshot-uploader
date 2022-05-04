PACKAGE ?= github.com/sosedoff/tm-snapshot-uploader
DOCKER_IMAGE ?= sosedoff/tm-snapshot-uploader

.PHONY: build
build:
	go build -o ./build/tm-snapshot-uploader

.PHONY: docker
docker:
	docker build --build-arg=PACKAGE=$(PACKAGE) -t $(DOCKER_IMAGE) .

.PHONY: docker-push
docker-push: docker
	docker push $(DOCKER_IMAGE)
