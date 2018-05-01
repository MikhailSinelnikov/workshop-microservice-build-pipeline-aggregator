# build stage
FROM golang:1.9.6-alpine3.7 AS build-env

RUN \
  apk update && \
  apk add git make

ADD Makefile /go/src/github.com/kublr/workshop-microservice-build-pipeline-aggregator/Makefile
WORKDIR /go/src/github.com/kublr/workshop-microservice-build-pipeline-aggregator

RUN make tools-update

ADD . /go/src/github.com/kublr/workshop-microservice-build-pipeline-aggregator

RUN make deps-update

RUN make build

# final stage
FROM alpine:3.7
COPY --from=build-env /go/src/github.com/kublr/workshop-microservice-build-pipeline-aggregator/target/server /opt/aggregator/server
ENTRYPOINT ["/opt/aggregator/server"]
EXPOSE 10000
