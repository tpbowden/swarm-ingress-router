FROM mhart/alpine-node:4.4

RUN http_proxy=http://proxy.intra.bt.com:8080 apk add --update git

WORKDIR /src

ADD package.json ./
RUN http_proxy=http://proxy.intra.bt.com:8080 https_proxy=http://proxy.intra.bt.com:8080 npm install

ADD . ./


ENTRYPOINT ["node"]

CMD ["./index.js"]

EXPOSE 8080
