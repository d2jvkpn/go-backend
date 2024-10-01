#!/bin/bash
set -eu -o pipefail # -x
_wd=$(pwd); _path=$(dirname $0 | xargs -i readlink -f {})

export PATH="$HOME/Apps/bin:$(go env GOPATH)/bin:$PATH"

# go get google/protobuf/timestamp.proto

cat > proto/log.proto <<EOF
syntax = "proto3";
package proto;

option go_package = "./proto";
// import "google/protobuf/timestamp.proto";

message LogData {
	string appName = 1;
	string appVersion = 2;

	string requestId = 3; // uuid
	string requestAt = 4; // RFC3339Milli
	// google.protobuf.Timestamp requestAt = 4;
	string ip = 5;
	string msg = 6;
	string query = 7;
	int32 status_code = 8;
	string error = 9;

	int64 latency_milli = 10;
	map<string, string> identity = 11;
	bytes data = 12; // json bytes
}

message LogId {
	string id = 1;
}

service LogService {
	rpc PushLog(LogData) returns(LogId) {};
}
EOF

ls -al proto/

protoc --go-grpc_out=./ --go_out=./ --proto_path=./proto proto/*.proto

sed -i '/^\tmustEmbedUnimplemented/s#\t#\t// #' proto/*_grpc.pb.go

go fmt ./...
go vet ./...
