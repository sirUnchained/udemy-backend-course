.PHONY: run
run:
	@go run ./cmd/api

.PHONY: seed
seed: 
	@go run cmd/migrate/seed/main.go

.PHONY: database
database:
	@docker compose up --build -d

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt