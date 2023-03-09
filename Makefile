all: npm build

npm:
	npm install --prefix web

generate:
	go generate ./...

build: build-frontend build-backend

snapshot: build-frontend
	goreleaser release --snapshot --rm-dist

dev-frontend:
	npm run dev --prefix web

dev-backend: generate
	go run --tags dev .

build-frontend:
	npm run build --prefix web

build-backend:
	go build -o bin/