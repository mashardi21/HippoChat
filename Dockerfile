FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /HippoChat
RUN mkdir /HippoChat/proto
RUN mkdir /HippoChat/settings

WORKDIR /HippoChat

COPY ./proto/service.pb.go /HippoChat/proto
COPY ./proto/service_grpc.pb.go /HippoChat/proto
COPY ./settings/settings.go /HippoChat/settings
COPY ./main.go /HippoChat

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o HippoChatServer .

CMD ./HippoChatServer
