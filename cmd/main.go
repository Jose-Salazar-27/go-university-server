package main

import (
	"log"

	"github.com/Jose-Salazar-27/go-university-server/internal/shared/kernel/infra/httpx"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	mig, err := migrate.New(
		"file://migrations",
		"postgres://admin:admin123@localhost:5432/university_platform?sslmode=disable")

	if err != nil {
		log.Fatalf("cannot initialize migrations \n got error: %s \n", err.Error())
	}

	if err := mig.Up(); err != nil {
		handleMigrationError(err)
	}

	log.Println("migrations executed successfully")

	app := fiber.New(fiber.Config{
		StructValidator: httpx.NewRequestValidator(),
	})

	log.Fatal(app.Listen(":3000"))
}

func handleMigrationError(err error) {
	switch err {
	case migrate.ErrNoChange:
		log.Println("migrations are up to day. no changes made")

	default:
		log.Fatalf("cannot run migrations \n got error: %s", err.Error())

	}
}
