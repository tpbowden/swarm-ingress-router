#!/bin/bash

echo "Creating networks..."
docker network create --driver=overlay frontends
docker network create --driver=overlay router-management
echo "Done"


echo "Creating Redis backend..."
docker service create --name router-storage --network router-management redis:3.2-alpine
echo "Done"


echo "Creating collector service..."
docker service create --name router-backend --constraint node.role==manager --mount \
  target=/var/run/docker.sock,source=/var/run/docker.sock,type=bind --network router-management \
  tpbowden/swarm-ingress-router:latest -r router-storage:6379 collector
echo "Done"

echo "Creating router service..."
docker service create --name router -p 80:8080 -p 443:8443 --network frontends \
  --network router-management tpbowden/swarm-ingress-router:latest -r \
  router-storage:6379 server -b 0.0.0.0
echo "Done"
