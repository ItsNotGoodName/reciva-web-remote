NPM_PREFIX := podman run --rm -it -p 3000:3000 -v "$(shell pwd)/left/web:/work" -w /work docker.io/library/node:16.13

all: npm build

npm:
	$(NPM_PREFIX) npm install

build: build-frontend build-backend

login:
	$(NPM_PREFIX) bash

dev-frontend:
	$(NPM_PREFIX) npm run dev

dev-backend:
	go run --tags dev .

build-frontend:
	$(NPM_PREFIX) npm run build

build-backend:
	go build -o bin/

snapshot:
	goreleaser release --snapshot --rm-dist