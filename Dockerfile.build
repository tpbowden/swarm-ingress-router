FROM golang:1.7-alpine

ENV CGO_ENABLED=0

RUN apk add -U make git curl
RUN go get github.com/Masterminds/glide
WORKDIR /go/src/github.com/tpbowden/swarm-ingress-router

COPY . ./
