FROM golang:1.14 as build

RUN apt-get update && apt-get install -y libsox-dev libsdl2-dev portaudio19-dev libopusfile-dev libopus-dev git

COPY ./ /go/src/github.com/itfactory-tm/itf-steunt-je

WORKDIR /go/src/github.com/itfactory-tm/itf-steunt-je

RUN go build -ldflags "-X main.revision=$(git rev-parse --short HEAD)" ./

FROM ubuntu:18.04

RUN apt-get update && apt-get install -y libsox-dev libsdl2-dev portaudio19-dev libopusfile-dev libopus-dev curl

RUN mkdir -p /go/src/github.com/itfactory-tm/itf-steunt-je/itf-steunt-je
WORKDIR /go/src/github.com/itfactory-tm/itf-steunt-je/itf-steunt-je
COPY ./audio /go/src/github.com/itfactory-tm/itf-steunt-je/itf-steunt-je/audio

COPY --from=build /go/src/github.com/itfactory-tm/itf-steunt-je/itf-steunt-je /usr/local/bin/

ENTRYPOINT /usr/local/bin/itf-steunt-je
