DATABASE_URI?="postgres://admin:admin123@localhost:5432/university_platform?sslmode=disable"

dev:
	go run ./cmd/main.go

migrate-up:
	migrate -path ./migrations -database $(DATABASE_URI) up

migrate-down:
	migrate -path ./migrations -database $(DATABASE_URI) down

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

fix-migration:
	migrate -path ./migrations -database $(DATABASE_URI) force 00001 
