# Run server in local development mode
# Requires to run a Postgres instance locally first
# brew services start postgresql@14
run:
	ADDR=:8080 JWT_SECRET=devsecret-change-me DATABASE_URL=$$(grep DATABASE_URL .env.local | cut -d= -f2-) PHOTO_UPLOAD_DIR=$$(grep PHOTO_UPLOAD_DIR .env.local | cut -d= -f2-) go run ./cmd/server

# Run server with Docker Compose
docker-up:
	docker compose up --build

ENV_DATABASE_URL=$(shell grep DATABASE_URL .env.local | cut -d= -f2-)

migrate-up:
	migrate -path migrations -database "$(ENV_DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(ENV_DATABASE_URL)" down 1

test:
	go test ./...
