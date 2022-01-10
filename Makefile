NAME := reciva-web-remote

npm:
	npm i --prefix web

build: generate backend

release: generate backend-linux backend-darwin backend-windows 

snapshot: 
	goreleaser release --snapshot --rm-dist

generate:
	go generate ./...

backend:
	go build -o bin/$(NAME)

dev-frontend:
	npm run dev --prefix web

dev-backend:
	go run --tags dev . 

backend-linux:
	GOOS=linux GOARCH=386 go build -o bin/$(NAME)-linux-386
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-linux-amd64
	GOOS=linux GOARCH=arm go build -o bin/$(NAME)-linux-arm
	GOOS=linux GOARCH=arm64 go build -o bin/$(NAME)-linux-arm64

backend-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/$(NAME)-darwin-arm64

backend-windows:
	GOOS=windows GOARCH=386 go build -o bin/$(NAME)-windows-386.exe
	GOOS=windows GOARCH=amd64 go build -o bin/$(NAME)-windows-amd64.exe
