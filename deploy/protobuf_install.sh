#!/bin/bash
set -eu -o pipefail # -x
_wd=$(pwd); _path=$(dirname $0 | xargs -i readlink -f {})

pb_version=${pb_version:-24.3}
pb_go_version=${pb_go_version:-1.31.0}

# https://github.com/protocolbuffers/protobuf/releases/download/v28.2/protobuf-28.2.zip
# mkdir -p target/protoc
# unzip protobuf-28.2.zip -f -d target/protoc

# https://github.com/protocolbuffers/protobuf-go/releases/download/v1.34.2/protoc-gen-go.v1.34.2.linux.amd64.tar.gz

pb_version=${pb_version:-28.2}
pb_go_version=${pb_go_version:-1.34.2}

mkdir -p ~/Apps/bin

if ! command protoc &> /dev/null; then
    curl -L -o ~/Downloads/protoc-${pb_version}-linux-x86_64.zip \
      https://github.com/protocolbuffers/protobuf/releases/download/v${pb_version}/protoc-${pb_version}-linux-x86_64.zip

    mkdir -p ~/Apps/protoc
    unzip ~/Downloads/protoc-${pb_version}-linux-x86_64.zip -d ~/Apps/protoc
fi

if ! command protoc-gen-go &> /dev/null; then
    curl -L -o ~/Downloads/protoc-gen-go.v${pb_go_version}.linux.amd64.tar.gz \
      https://github.com/protocolbuffers/protobuf-go/releases/download/v${pb_go_version}/protoc-gen-go.v${pb_go_version}.linux.amd64.tar.gz

    tar -xf ~/Downloads/protoc-gen-go.v${pb_go_version}.linux.amd64.tar.gz -C ~/Apps/bin
fi

if ! command protoc-gen-go-grpc &> /dev/null; then
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi
