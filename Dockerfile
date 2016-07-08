FROM alpine:3.4

COPY ./ingress-router /bin/
EXPOSE 8080
ENTRYPOINT ["/bin/ingress-router"]
