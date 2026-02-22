package auth

import (
	"database/sql"

	"github.com/Jose-Salazar-27/go-university-server/internal/auth/application"
	"github.com/Jose-Salazar-27/go-university-server/internal/auth/infra"
	"github.com/Jose-Salazar-27/go-university-server/internal/auth/infra/persistence"
	"github.com/gofiber/fiber/v3"
)

type Module struct {
	name   string
	engine *fiber.App
	db     *sql.DB
}

func (mod Module) ConfigureEnpoints() {
	group := mod.engine.Group(mod.name)

	userRepository := persistence.NewUserRepository(mod.db)
	h := infra.NewUserHandler(application.NewCreateUserInteractor(userRepository, nil))

	group.Post("", h.CreateUser)
}
