package configuration

import (
	"github.com/Jose-Salazar-27/go-university-server/internal/user/cmd"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/dto"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/entity"
	"github.com/mehdihadeli/go-mediatr"
)

func ConfigUserMediator(repo entity.UserRepository, factory *entity.UserFactory) error {
	if err := mediatr.RegisterRequestHandler[*cmd.CreateUserCommand, *dto.CreateUserResponseDto](
		cmd.NewCreateUserHandler(repo, factory),
	); err != nil {
		return err
	}
	return nil
}
