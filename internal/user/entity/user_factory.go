package entity

import (
	"errors"
	"strings"
	"time"

	sharedtypes "github.com/Jose-Salazar-27/go-university-server/internal/shared/types"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/pkg/types"
)

// UserFactory builds User entities with required validation and sensible defaults.
type UserFactory struct{}

func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

// Create constructs a User with required fields. It sets CreatedAt/UpdatedAt to now
// and IsActive to true by default.
func (f *UserFactory) Create(email, passwordHash, firstName, lastName string, userType types.UserType) (*User, error) {
	//TODO: implment regex here
	if !(strings.Contains(email, "@")) {
		return nil, errors.New("invalid email")
	}

	if !userType.IsValid() {
		return nil, errors.New("invalid type of user")
	}

	now := time.Now()
	u := &User{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		UserType:     userType,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return u, nil
}

// CreateWithID constructs a User and assigns the provided ID if it's not empty.
func (f *UserFactory) CreateWithID(id sharedtypes.ID, email, passwordHash, firstName, lastName string, userType types.UserType) (*User, error) {
	u, err := f.Create(email, passwordHash, firstName, lastName, userType)
	if err != nil {
		return nil, err
	}
	if !id.IsEmpty() {
		u.ID = id
	}
	return u, nil
}

// HydrateUser method return and User struct pointer without any validation since this method should be use
// only to handle data fetched from persistence
func (f UserFactory) HydrateUser(
	id string,
	email string,
	password string,
	firstname string,
	lastname string,
	userType string,
	avatarURl string,
	isActive bool,
	CreatedAt time.Time,
	UpdatedAt time.Time,
) *User {
	return &User{
		ID:           sharedtypes.MustGetID(id),
		Email:        email,
		PasswordHash: password,
		FirstName:    firstname,
		LastName:     lastname,
		UserType:     types.UserType(userType),
		AvatarUrl:    avatarURl,
		IsActive:     isActive,
		CreatedAt:    CreatedAt,
		UpdatedAt:    UpdatedAt,
	}
}
