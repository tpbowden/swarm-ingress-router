FROM alpine:3.4

COPY ./swarm-ingress-router /bin/
EXPOSE 8080 8443
ENTRYPOINT ["/bin/swarm-ingress-router"]
CMD ["-b", "0.0.0.0"]

