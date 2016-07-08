FROM alpine:3.4

COPY ./swarm-ingress-router /bin/
EXPOSE 8080
ENTRYPOINT ["/bin/swarm-ingress-router"]
