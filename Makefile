
run:
	@go run cmd/api/main.go

# start db & app
docker-up:
	@docker compose up -d --build

# Dev Mode: start db only
docker-db:
	@docker compose up -d db

docker-clear:
	@docker compose down -v

docker-config:
	@docker compose config