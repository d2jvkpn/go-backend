#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _dir=$(readlink -f `dirname "$0"`)


go install github.com/swaggo/swag/cmd/swag@latest
