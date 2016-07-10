# Swarm ingress router

Route DNS names to labelled Swarm services using Docker 1.12's internal service load balancing

* No external load balancer or config files needed making for easy deployments
* Automatic service discovery and load balancing handled by Docker
* Scaled and maintained by the Swarm for high resilience and performance
* Incredibly lightweight image (less than 15MB image after decompression)

## Installation

Fix of all you will need to create a network for your frontend services to run on:

    docker network create --driver=overlay frontends

Then you have to start the router on this network. It also needs to be able to communicate with Docker:

    docker service create --name router -p 80:8080 \
      --mount target=/var/run/docker.sock,source=/var/run/docker.sock,type=bind \
      --network frontends tpbowden/swarm-ingress-router:latest

Now you can start your frontend services and they will be available on all of your Swarm nodes:

    docker service create --name frontend --label ingress=true --label ingress.dnsname=example.local \
      --label ingress.targetport=80 --network frontends nginx:stable-alpine

If you now add a DNS record for `example.local` pointing to your Docker node you will be routed to the service.
The service must be restricted to run only on master nodes (as it has to query for services).

## Usage

In order for the router to pick up a service, the service must have the following labels:

* `ingress=true`
* `ingress.dnsname=<your service's external DNS name>`
* `ingress.targetport=<your service's externally-facing port>`

You do not need to publish this port externally as long as your services are both on a shared network.

## Todo

* Allow services to listen for HTTPS and decrypt using certificates stored elsewhere (probably Vault)
* Better logging
* Command line arguments for log level, bind address etc
* Use Docker events to sync and long polling as a fallback
