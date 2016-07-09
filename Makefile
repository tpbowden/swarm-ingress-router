GOOS=linux
GOARCH=amd64
TAG?=latest

install-deps:
	@glide install

compile: install-deps
	@echo "Building binary"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -v .
	@echo "Done"

build-image: compile
	@docker build -t tpbowden/ingress-router:$(TAG) .

release: build-image
	@docker push tpbowden/ingress-router:$(TAG)
