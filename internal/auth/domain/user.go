package domain

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Jose-Salazar-27/go-university-server/internal/shared/kernel/valueobject"
)

var (
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrEmptyFirstName    = errors.New("first name cannot be empty")
	ErrEmptyLastName     = errors.New("last name cannot be empty")
	ErrInvalidUserType   = errors.New("invalid user type")
	ErrEmptyPasswordHash = errors.New("password hash cannot be empty")
)

type UserRepository interface {
	Create(ctx context.Context, u *User) (err error)
}

// User represents a user entity in the domain
type User struct {
	ID           valueobject.ID `json:"id"`
	Email        string         `json:"email"`
	PasswordHash string         `json:"-"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	UserType     UserType       `json:"user_type"`
	AvatarURL    *string        `json:"avatar_url,omitempty"`
	IsActive     bool           `json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// NewUser creates a new User instance with validation
func NewUser(
	email string,
	passwordHash string,
	firstName string,
	lastName string,
	userType UserType,
	avatarURL *string,
) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, fmt.Errorf("email validation failed: %w", err)
	}

	if passwordHash == "" {
		return nil, ErrEmptyPasswordHash
	}

	if strings.TrimSpace(firstName) == "" {
		return nil, ErrEmptyFirstName
	}

	if strings.TrimSpace(lastName) == "" {
		return nil, ErrEmptyLastName
	}

	if !userType.IsValid() {
		return nil, ErrInvalidUserType
	}

	now := time.Now()

	user := &User{
		ID:           valueobject.NewID(),
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		FirstName:    strings.TrimSpace(firstName),
		LastName:     strings.TrimSpace(lastName),
		UserType:     userType,
		AvatarURL:    avatarURL,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return user, nil
}

// UserFromPersistence creates a User instance from database records
// This method assumes data from database is already validated and doesn't perform additional validation
func UserFromPersistence(
	id valueobject.ID,
	email string,
	passwordHash string,
	firstName string,
	lastName string,
	userType UserType,
	avatarURL *string,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		ID:           id,
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		FirstName:    strings.TrimSpace(firstName),
		LastName:     strings.TrimSpace(lastName),
		UserType:     userType,
		AvatarURL:    avatarURL,
		IsActive:     isActive,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// Business Logic Methods

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(firstName, lastName string, avatarURL *string) error {
	if strings.TrimSpace(firstName) == "" {
		return ErrEmptyFirstName
	}

	if strings.TrimSpace(lastName) == "" {
		return ErrEmptyLastName
	}

	u.FirstName = strings.TrimSpace(firstName)
	u.LastName = strings.TrimSpace(lastName)
	u.AvatarURL = avatarURL
	u.UpdatedAt = time.Now()

	return nil
}

// UpdatePassword updates the user's password hash
func (u *User) UpdatePassword(passwordHash string) error {
	if passwordHash == "" {
		return ErrEmptyPasswordHash
	}

	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()

	return nil
}

// Deactivate deactivates the user
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// IsProfessor checks if the user is a professor
func (u *User) IsProfessor() bool {
	return u.UserType == UserTypeProfessor
}

// IsStudent checks if the user is a student
func (u *User) IsStudent() bool {
	return u.UserType == UserTypeStudent
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.UserType == UserTypeAdmin
}

// Validate performs comprehensive validation of the user
func (u *User) Validate() error {
	if err := u.ID.Validate(); err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	if err := validateEmail(u.Email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if u.PasswordHash == "" {
		return ErrEmptyPasswordHash
	}

	if strings.TrimSpace(u.FirstName) == "" {
		return ErrEmptyFirstName
	}

	if strings.TrimSpace(u.LastName) == "" {
		return ErrEmptyLastName
	}

	if !u.UserType.IsValid() {
		return ErrInvalidUserType
	}

	return nil
}

// Private helper functions
var emailRegex = regexp.MustCompile(
	`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`,
)

func validateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	if match := emailRegex.MatchString(email); !match {
		return errors.New("invalid email format")
	}
	return nil
}
