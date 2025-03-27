#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _path=$(dirname $0)


cd ${_path}

set -x

for f in $(ls * | grep -v git_hooks.sh); do
    ln -frs $f ${_wd}/.git/hooks/
done
