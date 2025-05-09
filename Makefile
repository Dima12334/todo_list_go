migration:
	@migrate create -ext sql -dir ./migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

swag:
	swag init -g internal/app/app.go

test:
	go test -v ./...
