build: bin vendor
	go build -o bin/jaeger-lite .

bin:
	mkdir -p bin

vendor:
	dep ensure

clean:
	rm -rf bin vendor

image:
	docker run -t --rm \
	-v `pwd`:/go/src/github.com/factorysh/jaeger-lite/ \
	-w /go/src/github.com/factorysh/jaeger-lite/ \
	golang \
	make build