FROM mhart/alpine-node:4.4
ENV NODE_ENV=production

ARG HTTP_PROXY=""
ARG HTTPS_PROXY=""

WORKDIR /src

ADD package.json ./
RUN http_proxy=${HTTP_PROXY} https_proxy=${HTTPS_PROXY} npm install

ADD . ./

ENTRYPOINT ["bin/start"]

EXPOSE 8080 9090
