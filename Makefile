run:
	ADDR=:8080 JWT_SECRET=devsecret-change-me DATABASE_URL=$$(grep DATABASE_URL .env.local | cut -d= -f2-) go run ./cmd/server

docker-up:
	docker compose up --build

migrate-up:
	migrate -path migrations -database "$$DATABASE_URL" up

migrate-down:
	migrate -path migrations -database "$$DATABASE_URL" down 1

test:
	go test ./...
