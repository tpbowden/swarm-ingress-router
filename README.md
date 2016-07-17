# Swarm ingress router

Route DNS names to labelled Swarm services using Docker 1.12's internal service load balancing

* No external load balancer or config files needed making for easy deployments
* Integrated TLS decryption for services which provide a certificate and key
* Automatic service discovery and load balancing handled by Docker
* Scaled and maintained by the Swarm for high resilience and performance
* Incredibly lightweight image (less than 20MB after decompression)

## Installation

First of all you will need to create a network for your frontend services to run on:

    docker network create --driver=overlay frontends

Then you have to start the router on this network.
It must listen on the standard HTTP/HTTPS ports and be able to communicate with the Docker daemon:

    docker service create --name router -p 80:8080 -p 443:8443 \
      --mount target=/var/run/docker.sock,source=/var/run/docker.sock,type=bind \
      --network frontends tpbowden/swarm-ingress-router:latest

Now you can start your frontend service and it will be available on all of your Swarm nodes:

    docker service create --name frontend --label ingress=true --label ingress.dnsname=example.local \
      --label ingress.targetport=80 --network frontends --label ingress.tls=true --label ingress.forcetls=true \
      --label ingress.cert="$(cat fixtures/cert.crt)" --label ingress.key="$(cat fixtures/key.key)" \
      nginx:stable-alpine

If you now add a DNS record for `example.local` pointing to your Docker node you will be routed to the service.
The service must be restricted to run only on master nodes (as it has to query for services).

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
