#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _dir=$(readlink -f `dirname "$0"`)


command -v yq > /dev/null

export IMAGE_Tag=$1 HTTP_Port=$2

export APP_Name=$(yq .app_name project.yaml) \
  IMAGE_Name=$(yq .image_name project.yaml) \
  USER_UID=$(id -u) \
  USER_GID=$(id -g)

mkdir -p data/$APP_Name
envsubst < ${_dir}/compose.backend-api.yaml > compose.backend-api.yaml

####
exit 0
docker-compose up -d
docker-compose logs

####
exit 0
container=xxxx
# 'docker-compose down' removes running containers only, not stopped containers
[ ! -z "$(docker ps --all --quiet --filter name=$container)" ] && docker rm -f $container

docker stop $container
docker rm $container
