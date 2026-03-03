package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/Jose-Salazar-27/go-university-server/internal/user/dto"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/entity"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserHandler struct {
	repository entity.UserRepository
	factory    *entity.UserFactory
}

func NewCreateUserHandler(r entity.UserRepository, f *entity.UserFactory) *CreateUserHandler {
	return &CreateUserHandler{repository: r, factory: f}
}

func (h CreateUserHandler) Handle(ctx context.Context, command *CreateUserCommand) (*dto.CreateUserResponseDto, error) {
	// hash password, doesn't save it as plain text
	bytes, err := bcrypt.GenerateFromPassword([]byte(command.Password), 14)
	if err != nil {
		return nil, errors.New("cannot process password")
	}

	user, err := h.factory.CreateWithID(
		command.ID,
		command.Email,
		string(bytes),
		command.FirstName,
		command.LastName,
		command.UserType,
	)
	if err != nil {
		return nil, err
	}

	tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.repository.Create(tctx, user); err != nil {
		return nil, err
	}

	return &dto.CreateUserResponseDto{ID: user.ID.String()}, nil
}
