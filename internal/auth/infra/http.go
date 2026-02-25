package infra

import (
	"net/http"

	"github.com/Jose-Salazar-27/go-university-server/internal/auth/application"
	"github.com/gofiber/fiber/v3"
)

type userHandler struct {
	interactor application.UserInteractor
}

func NewUserHandler(uc application.UserInteractor) *userHandler {
	return &userHandler{uc}
}

func (h userHandler) CreateUser(c fiber.Ctx) error {
	var req application.CreateUserInput

	if err := c.Bind().Body(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	data, err := h.interactor.Create(req)
	if err != nil {
		return c.Status(http.StatusOK).JSON(data)
	}

	return nil
}
