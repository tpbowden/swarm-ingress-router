# Swarm ingress router

Routes DNS names to labelled Swarm services

## Installation

    docker service create --name router -p 8080 \
      --mount target=/var/run/docker.sock,source=/var/run/docker.sock,type=bind \
      tpbowden/ingress-router:latest

The service must be restricted to run only on master nodes

## Usage

In order for the router to pick up a service, the service must have the following labels:

* `ingress=true`
* `ingress.dnsname=<your service's external DNS name>`
* `ingress.targetport=<your service's externally-facing port>`

You do not need to publish this port externally
