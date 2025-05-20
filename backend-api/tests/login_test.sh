#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _dir=$(readlink -f `dirname "$0"`)


api=${1:-http://127.0.0.1:4011}
config=${2:-configs/account.yaml}

email=$(yq .account.email $config)
password=$(yq .account.password $config)

data=$(jq -n --arg email "$email" --arg password "$password" '{email:$email,password:$password}')

curl -X POST $api/api/v1/open/account/login -d "$data"
