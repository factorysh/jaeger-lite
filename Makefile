GIT_VERSION?=$(shell git describe --tags --always --abbrev=42 --dirty)

build: bin vendor
	go build \
		-ldflags "-X github.com/factorysh/jaeger-lite/version.version=$(GIT_VERSION)" \
		-o bin/jaeger-lite \
		.

bin:
	mkdir -p bin

vendor:
	dep ensure

clean:
	rm -rf bin vendor

image: docker-build
	docker build -t jaeger-lite .

docker-build:
	docker run -t --rm \
	-v `pwd`:/go/src/github.com/factorysh/jaeger-lite/ \
	-w /go/src/github.com/factorysh/jaeger-lite/ \
	-e CGO_ENABLED=0 \
	bearstech/golang-dev \
	make build

.PHONY: demo
demo: image
	cd demo \
	&& docker-compose down \
	&& docker-compose up -d traefik \
	&& sleep 1 \
	&& docker-compose up client \
	&& sleep 10 \
	&& docker-compose up promclient | grep apdex