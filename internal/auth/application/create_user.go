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
		Email     string          `json:"email"`
		Password  string          `json:"password"`
		FirstName string          `json:"first_name"`
		LastName  string          `json:"last_name"`
		UserType  domain.UserType `json:"user_type"`
	}

	CreateUserOutput struct {
		ID        string
		CreatedAt time.Time
	}
)

type CreateUserInteractor struct {
	repository domain.UserRepository
	storage    shared.Storage
}

func NewCreateUserInteractor(r domain.UserRepository, s shared.Storage) *CreateUserInteractor {
	return &CreateUserInteractor{r, s}
}

func (interactor CreateUserInteractor) Execute(in CreateUserInput) (CreateUserOutput, error) {
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
