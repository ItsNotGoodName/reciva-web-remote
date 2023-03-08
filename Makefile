all: npm build

npm:
	npm install

build: build-frontend build-backend

snapshot: build-frontend build-snapshot

dev-frontend:
	npm run dev

dev-backend:
	go run --tags dev .

build-frontend:
	npm run build

build-backend:
	go build -o bin/

build-snapshot:
	goreleaser release --snapshot --rm-dist