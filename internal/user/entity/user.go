package entity

import (
	"context"
	"time"

	sharedtypes "github.com/Jose-Salazar-27/go-university-server/internal/shared/types"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/pkg/types"
)

type UserRepository interface {
	Create(ctx context.Context, u *User) error
}

type User struct {
	ID           sharedtypes.ID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	UserType     types.UserType
	AvatarUrl    string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
