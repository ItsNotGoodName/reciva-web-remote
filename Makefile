build: frontend backend

frontend:
	npm run build --prefix web

backend:
	go build -tags="prod"

dev-frontend:
	npm run dev --prefix web

dev-backend:
	go run .

npm:
	npm i --prefix web