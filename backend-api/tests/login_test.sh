#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _dir=$(readlink -f `dirname "$0"`)


apiUrl=${1:-http://127.0.0.1:4011/local}
config=${2:-configs/account.yaml}

email=$(yq .account.email $config)
password=$(yq .account.password $config)

data=$(jq -n --arg email "$email" --arg password "$password" '{email:$email,password:$password,platform:"web"}')

curl -i -X POST "$apiUrl/api/v1/open/account/login?platform=web" -d "$data" | jq


exit

curl -i -X GET -H "Authorization: Bearer $token" $apiUrl/api/v1/auth/hello

curl -i -X POST -H "Authorization: Bearer $token" $apiUrl/api/v1/auth/account/logout
