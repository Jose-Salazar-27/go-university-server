package infra

import (
	"net/http"

	"github.com/Jose-Salazar-27/go-university-server/internal/auth/application"
	"github.com/gofiber/fiber/v3"
)

type createUserHandler struct {
	userCreator *application.CreateUserInteractor
}

func NewCreateUserHandler(uc *application.CreateUserInteractor) *createUserHandler {
	return &createUserHandler{uc}
}

func (h createUserHandler) Handle(c fiber.Ctx) error {
	var req application.CreateUserInput

	if err := c.Bind().Body(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	data, err := h.userCreator.Execute(req)
	if err != nil {
		return c.Status(http.StatusOK).JSON(data)
	}

	return nil
}
