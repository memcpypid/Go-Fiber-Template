.PHONY: run migrate seed build mod-tidy

run:
	go run cmd/server/main.go

migrate:
	go run cmd/migrate/main.go

seed:
	go run cmd/seed/main.go

build:
	go build -o tmp/main cmd/server/main.go

mod-tidy:
	go get github.com/gofiber/fiber/v2
	go mod tidy
