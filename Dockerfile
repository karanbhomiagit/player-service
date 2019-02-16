FROM golang:1.11.2-alpine3.8
MAINTAINER M. - Karan Bhomia

ENV SOURCES /go/src/github.com/karanbhomiagit/player-service/

RUN apk update -qq && apk add git

COPY . ${SOURCES}
RUN go get github.com/stretchr/testify
RUN cd ${SOURCES} && CGO_ENABLED=0 go install

ENV PORT 8080
ENV MAX_CONCURRENT_REQUESTS 100

EXPOSE 3000

ENTRYPOINT player-service