#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _dir=$(readlink -f `dirname "$0"`)


export USER_UID=$(id -u) USER_GID=$(id -g)

mkdir -p configs logs data/postgres data/redis

envsubst < ${_dir}/compose.databases.yaml > compose.databases.yaml
