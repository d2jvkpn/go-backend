#!/bin/bash
set -eu -o pipefail; _wd=$(pwd); _dir=$(readlink -f `dirname "$0"`)


app=$1
language=${2}

#### 1. create the project
# --template=vue # vue javascript project
if [[ "$language" == "ts" || "$language" == "typescript" ]]; then
    npm init vite@latest $app -- --template=vue-ts
elif [[ "$language" == "js" || "$language" == "javascript" ]]; then
    npm init vite@latest $app -- --template=vue
else
    >&2 echo 'unknown language: js/javascript ts/typescript'
    exit 1
fi

cd "$app"

npm install
#npm run format
npm fund

#### 2. add packages
npm install --save-dev @types/node
npm install --save-dev vite-plugin-vue-devtools
npm install --save-dev unplugin-vue-components # unplugin-auto-import
npm install --save-dev tailwindcss @tailwindcss/vite postcss autoprefixer
# ?? error: npx tailwindcss init -p

# https://tailwindcss.com/docs/installation/using-vite
npm install vue-router
npm install element-plus @element-plus/icons-vue
npm install pinia pinia-plugin-persistedstate

#### 3. setup env
cat > env <<"EOF"
# path: .env
PORT=3001

VITE_ENV=local
VITE_BASE_PATH=/local
VITE_API_URL=http://localhost:3011
EOF

cp env .env

#### 4. setup makefile
cat > Makefile <<"EOF"
#!/bin/make

include .env

SHELL = /bin/bash

run:
	# node node_modules/vite/bin/vite.js --help
	npm run dev -- --port=$(PORT) --host=0.0.0.0 # --mode dev

build:
	rm -rf target/dist
	# node node_modules/vite/bin/vite.js build --help
	npm run build -- --base=$(VITE_BASE_PATH) --outDir=target/dist$(VITE_BASE_PATH) # --mode dev
	ls -alt target/dist
EOF

#### 5. setup gitignore
cat >> .gitignore <<EOF

####
.env
.env.*
target/
cache/

docker-compose.yaml
docker-compose.yml
compose.yaml
compose.yml
EOF
