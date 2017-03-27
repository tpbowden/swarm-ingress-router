GOOS=linux
GOARCH=amd64
TAG?=dev
CERTIFICATE=`cat fixtures/cert.crt`
KEY=`cat fixtures/key.key`

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
	@docker tag tpbowden/swarm-ingress-router:$(TAG) tpbowden/swarm-ingress-router:latest
	@docker push tpbowden/swarm-ingress-router:latest

release-dev: build-image
	@docker push tpbowden/swarm-ingress-router:dev

test:
	@go test -cover `go list ./... | grep -v '/vendor/'`

dev-compile:
	@docker build -t ingress-router-dev -f Dockerfile.build .
	@docker run --name ingress-router-dev -it ingress-router-dev make compile
	@docker cp ingress-router-dev:/go/src/github.com/tpbowden/swarm-ingress-router/swarm-ingress-router .
	@docker rm ingress-router-dev
	@docker rmi ingress-router-dev

dev: dev-compile
	@docker build -t swarm-ingress-router-dev:latest .
	@docker swarm init
	@docker network create --driver=overlay frontends
	@docker network create --driver=overlay router-management
	@docker service create --name router-storage --network router-management redis:3.2-alpine
	@docker service create --name router-backend --constraint node.role==manager --mount \
		target=/var/run/docker.sock,source=/var/run/docker.sock,type=bind,readonly --network router-management \
		swarm-ingress-router-dev:latest -r router-storage:6379 collector
	@docker service create --name router -p 80:8080 -p 443:8443 --network frontends \
		--network router-management --user nobody swarm-ingress-router-dev:latest -r \
		router-storage:6379 server -b 0.0.0.0
	@docker service create --name frontend --label ingress=true --label ingress.dnsname=example.local \
		--label ingress.targetport=80 --network frontends --label ingress.tls=true --label \
		ingress.forcetls=true --label ingress.cert="$(CERTIFICATE)" --label \
		ingress.key="$(KEY)" nginx:stable-alpine
	@echo "Development services running, run make cleanup to stop, or make reload to load current code"

cleanup:
	@docker swarm leave --force
	@docker network rm frontends router-management
	@docker rmi swarm-ingress-router-dev
	@rm ./swarm-ingress-router

reload: cleanup dev

