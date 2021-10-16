NAME := reciva-web-remote
TAGS := "prod"

build: frontend backend

release: frontend backend-linux backend-mac backend-windows 

frontend:
	npm run build --prefix web

backend:
	go build -tags=$(TAGS)

dev-frontend:
	npm run dev --prefix web

dev-backend:
	go run .

npm:
	npm i --prefix web

backend-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(NAME)-linux-amd64 -tags=$(TAGS)
	GOOS=linux GOARCH=arm go build -o bin/$(NAME)-linux-arm -tags=$(TAGS)
	GOOS=linux GOARCH=arm64 go build -o bin/$(NAME)-linux-arm64 -tags=$(TAGS)

backend-mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(NAME)-darwin-amd64 -tags=$(TAGS)
	GOOS=darwin GOARCH=arm64 go build -o bin/$(NAME)-darwin-arm64 -tags=$(TAGS)

backend-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(NAME)-windows-amd64.exe -tags=$(TAGS)
