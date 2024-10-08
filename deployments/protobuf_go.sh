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

enum event_level {
	debug = 0;
	info = 1;
	warning = 2;
	error = 3;
	critical = 4;
}

message LogRequest {
	string event_id = 1;
	string event_at = 2;
	event_level event_level = 3;
	string appName = 4;
	string appVersion = 5;

	string service = 6; // service name
	string id = 7; // uuid
	string at = 8; // RFC3339Milli
	// google.protobuf.Timestamp at = 8;
	string biz_name = 9; // POST@/api/v1/open/login
	map<string,string> biz_data = 10; // query, status, error, client
	map<string,string> identities = 11; // accountId, tokenId, ip, role
	string code = 12; // custom app code: ok, warn, error, panic
	double latency_milli = 13;

	repeated string labels = 14;
	bytes data = 15; // json bytes
}

message LogResponse {
	string event_id = 1;
	string id = 2;
}

service LogService {
	rpc Push(LogRequest) returns(LogResponse) {};
}
EOF

ls -al proto/

protoc --go-grpc_out=./ --go_out=./ --proto_path=./proto proto/*.proto

sed -i '/^\tmustEmbedUnimplemented/s#\t#\t// #' proto/*_grpc.pb.go

go fmt ./...
go vet ./...
