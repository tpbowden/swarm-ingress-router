GOOS=linux
GOARCH=amd64
TAG?=latest

all: install-deps test

install-deps:
	@glide install

compile: install-deps
	@echo "Building binary"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -v .
	@echo "Done"

release-compile:
	@docker build -t ingress-router-build -f Dockerfile.build .
	@docker run --name ingress-router-build -it ingress-router-build make compile
	@docker cp ingress-router-build:/go/src/github.com/tpbowden/swarm-ingress-router/swarm-ingress-router .
	@docker rm ingress-router-build
	@docker rmi ingress-router-build

build-image: release-compile
	@docker build -t tpbowden/swarm-ingress-router:$(TAG) .

release: build-image
	@docker push tpbowden/swarm-ingress-router:$(TAG)

test:
	@go test -cover `go list ./... | grep -v '/vendor/'`
