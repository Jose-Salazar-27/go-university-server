package cmd

import (
	"errors"
	"fmt"
	"time"

	sharedtypes "github.com/Jose-Salazar-27/go-university-server/internal/shared/types"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/pkg/types"
)

type CreateUserCommand struct {
	ID        sharedtypes.ID `json:"id" validate:"required"`
	Email     string         `json:"email" validate:"required,email"`
	Password  string         `json:"password" validate:"required,min=8"`
	FirstName string         `json:"first_name" validate:"required"`
	LastName  string         `json:"last_name" validate:"required"`
	UserType  types.UserType `json:"user_type" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCreateUserCommand(email, password, firstname, lastname, utype string) (*CreateUserCommand, error) {
	newID, err := sharedtypes.NewID()
	if err != nil {
		return nil, fmt.Errorf("cannot create id. got error: %w", err)
	}

	userType := types.UserType(utype)
	if !userType.IsValid() {
		return nil, errors.New("invalid type of user")
	}

	now := time.Now()

	command := &CreateUserCommand{
		ID:        newID,
		Email:     email,
		Password:  password,
		FirstName: firstname,
		LastName:  lastname,
		UserType:  userType,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return command, nil
}
