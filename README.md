# Swarm ingress router

Route DNS names to labelled Swarm services using Docker 1.12's internal service load balancing

[![Build Status](https://travis-ci.org/tpbowden/swarm-ingress-router.svg?branch=master)](https://travis-ci.org/tpbowden/swarm-ingress-router) [![Go Report Card](https://goreportcard.com/badge/github.com/tpbowden/swarm-ingress-router)](https://goreportcard.com/report/github.com/tpbowden/swarm-ingress-router)

WARNING: This application depends on features in Docker's currently experimental realease 1.12. There are still
serveral bugs which can cause things to break, especially related to container IP address management across multiple
Swarm hosts.

## Features

* No external load balancer or config files needed making for easy deployments
* Integrated TLS decryption for services which provide a certificate and key
* Automatic service discovery and load balancing handled by Docker
* Scaled and maintained by the Swarm for high resilience and performance
* Incredibly lightweight image (less than 20MB after decompression)

## Installation

These are the manual steps to install the app. You can run `bootstrap.sh` to quickly create the required
services automatically.

First of all you will need to create a network for your frontend services to run on and one for storage:

    docker network create --driver=overlay frontends
    docker network create --driver=overlay router-management

Next you need to start Redis which will store service configuration

    docker service create --name router-storage --network router-management redis:3.2-alpine

Then you have to start the router's backend on management network. The service must be restricted to
run only on master nodes (as it has to query for services).

    docker service create --name router-backend --constraint node.role=manager --mount \
      target=/var/run/docker.sock,source=/var/run/docker.sock,type=bind --network router-management \
      tpbowden/swarm-ingress-router:latest -r router-storage:6379 collector

Now you can start the router's frontend on both the management and frontend network.
It must listen on the standard HTTP/HTTPS ports 

    docker service create --name router --mode global -p 80:8080 -p 443:8443 --network frontends \
      --network router-management tpbowden/swarm-ingress-router:latest -r \
      router-storage:6379 server -b 0.0.0.0

### Start a demo app

Finally, start your frontend service on the frontends network and it will be available on all of your Swarm nodes:

    docker service create --name frontend --label ingress=true --label ingress.dnsname=example.local \
      --label ingress.targetport=80 --network frontends --label ingress.tls=true --label \
      ingress.forcetls=true --label ingress.cert="$(cat fixtures/cert.crt)" --label \
      ingress.key="$(cat fixtures/key.key)" nginx:stable-alpine

If you add a DNS record for `example.local` pointing to your Docker node in `/etc/hosts` you will be
routed to the service.

## Usage

In order for the router to pick up a service, the service must have the following labels:

* `ingress=true`
* `ingress.dnsname=<your service's external DNS name>`
* `ingress.targetport=<your service's externally-facing port>`

For TLS you need the following lables:

* `ingress.tls=true`
* `ingress.cert="$(cat <your crt file>)"`
* `ingress.key="$(cat <your key file>"`

To force services to force TLS you can also use the following label:

* `ingress.forcetls=true`

You do not need to publish the service's port as a node port as long as it is exposed internally and on the same network
as the router.

## Todo

* Better logging
* Command line argument for log level
* Use Docker events to stay in sync and long polling as a fallback
* Create a docker-compose file which can be converted into a stack
