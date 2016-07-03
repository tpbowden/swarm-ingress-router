FROM mhart/alpine-node:4.4

ARG HTTP_PROXY=""
ARG HTTPS_PROXY=""

RUN http_proxy=${HTTP_PROXY} apk add --update git

WORKDIR /src

ADD package.json ./
RUN http_proxy=${HTTP_PROXY} https_proxy=${HTTPS_PROXY} npm install

ADD . ./


ENTRYPOINT ["node"]

CMD ["./index.js"]

EXPOSE 8080 9090
