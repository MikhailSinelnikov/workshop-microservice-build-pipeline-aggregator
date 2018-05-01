# Setup dev envioronment

Install Go (should be available in path)

Install `protoc` (should be available in path)

Other necessary tools are installed via `go get` and are part of the `Makefile`

# Build demo from scratch

`git clone` project

`cd` to the project directory

`make tools-update` to install or update required tools

`make deps-update` to get vendor dependencies

`make gen` to generate protobuf/grpc code

`make build` to build server binary

`docker build -t aggregator .` to build docker container

`./target/server` to run server from a binary

`docker run -p 11000:11000 aggregator` to run server from a container

Example deployment to Kubernetes cluster with Isto from local machine:

```
kubectl apply -f kubernetes/service.yaml
cat kubernetes/deployment.yaml | sed -e 's^{{image}}^IMG^g' | sed -e 's^{{version}}^VER^g' | istioctl kube-inject -f - | kubectl apply -f -
```

# Development

## Change API

Edit `api/protobuf-spec/*.proto`

`make gen` to update generated protobuf/grpc code

Build server as usual.

## Change dependencies

Edit `Gopkg.toml`

`make deps-update` to update vendor dependencies

## Test with local REPL client

See https://github.com/njpatel/grpcc for more details

```
# host CLI
docker run -ti -v "$(pwd):/proto" --net=host therealplato/grpcc-container:latest bash

# in container CLI
cd /proto
grpcc -p api/protobuf-spec/aggregator.proto -i -a 127.0.0.1:11000

# in grpcc REPL
# select package and run
client.aggregate({},pr)
```

## Keep Go build `Makefile` reusable for local and Docker builds

`Makefile` is used to build a binary executable from `go` sources.

The same `Makefile` may be used for
- local build environment setup (`tools-update`, `deps-update` targets),
- generated sources update (`gen` target),
- local build (`build` target)
- Docker image build (`tools-update`, `deps-update`, and `build` targets are
  run from the `Dockerfile`)

When changing `Makefile` make sure that all these scenarios are still supported,
which is especially important for fast local development.

## Keep Docker build `Dockerfile` cacheable

`Dockerfile` used in the project defines multi-stage build.

Some efforts were applied to ensure that docker build can cache layers and cached
layers can be reused as much as possible:
- `.dockerignore` only includes files necessary to build Docker image,
- `Dockerfile` adds source and build files gradually, so that on each step only
  necessary files are added.

When modifying code, make sure that iterative builds optimization goal is
maintained, which is especially important for fast local development.
