up:
	docker-compose up

up-with-build:
	docker-compose up --build

down:
	docker-compose down

test:
	docker-compose exec app go test -v ./...

migration:
	@docker-compose exec app migrate create -ext sql -dir ./migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@docker-compose exec app go run cmd/migrate/main.go up

migrate-down:
	@docker-compose exec app go run cmd/migrate/main.go down

swag:
	swag init -g internal/app/app.go
