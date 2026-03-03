package endpoint

import (
	"net/http"

	"github.com/Jose-Salazar-27/go-university-server/internal/user/cmd"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/dto"
	"github.com/gofiber/fiber/v3"
)

type createUserEnpoint struct {
	// dependencies goes here
	Router fiber.Router
}

func (ep createUserEnpoint) MapEnpoint() {
	ep.Router.Post("", ep.handler())
}

func (ep createUserEnpoint) handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.Context()

		request := &dto.CreateUserDto{}
		if err := c.Bind().Body(request); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		command, err := cmd.NewCreateUserCommand(
			request.Email,
			request.Password,
			request.FirstName,
			request.LastName,
			request.UserType,
		)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

}
