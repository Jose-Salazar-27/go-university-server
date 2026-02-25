package application

import (
	"context"
	"errors"
	"time"

	"github.com/Jose-Salazar-27/go-university-server/internal/auth/domain"
	shared "github.com/Jose-Salazar-27/go-university-server/internal/shared/kernel/domain"
	"golang.org/x/crypto/bcrypt"
)

type (
	CreateUserInput struct {
		Email     string          `json:"email" validate:"required,email"`
		Password  string          `json:"password" validate:"required,min=8"`
		FirstName string          `json:"first_name" validate:"required"`
		LastName  string          `json:"last_name" validate:"required"`
		UserType  domain.UserType `json:"user_type" validate:"required"`
	}

	CreateUserOutput struct {
		ID        string
		CreatedAt time.Time
	}
)

//go:generate mockgen -destination user_interactor_mock.go -package application . UserInteractor
type UserInteractor interface {
	Create(input CreateUserInput) (CreateUserOutput, error)
}

type userInteractor struct {
	repository domain.UserRepository
	storage    shared.Storage
}

func NewCreateUserInteractor(r domain.UserRepository, s shared.Storage) *userInteractor {
	return &userInteractor{r, s}
}

func (interactor userInteractor) Create(in CreateUserInput) (CreateUserOutput, error) {
	// hash password, doesn't save it as plain text
	bytes, err := bcrypt.GenerateFromPassword([]byte(in.Password), 14)
	if err != nil {
		return CreateUserOutput{}, errors.New("cannot process password")
	}

	user, err := domain.NewUser(in.Email, string(bytes), in.FirstName, in.LastName, in.UserType, nil)
	if err != nil {
		return CreateUserOutput{}, err
	}

	avatarUrl := interactor.storage.BuildObjectURL("users", user.ID.String(), "avatar.jpg")
	user.AvatarURL = &avatarUrl

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := interactor.repository.Create(ctx, user); err != nil {
		return CreateUserOutput{}, err
	}

	return CreateUserOutput{ID: user.ID.String(), CreatedAt: user.CreatedAt}, nil
}
