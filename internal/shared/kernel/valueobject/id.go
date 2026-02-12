package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidID = errors.New("invalid UUID format")
	ErrEmptyID   = errors.New("ID cannot be empty")
)

// ID represents a domain identifier as a value object
type ID struct {
	value uuid.UUID
}

// NewID creates a new ID with a random UUID
func NewID() ID {
	return ID{value: uuid.New()}
}

// IDFromUUID creates an ID from an existing UUID
func IDFromUUID(u uuid.UUID) ID {
	return ID{value: u}
}

// IDFromString creates an ID from a UUID string
func IDFromString(s string) (ID, error) {
	if s == "" {
		return ID{}, ErrEmptyID
	}

	u, err := uuid.Parse(s)
	if err != nil {
		return ID{}, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}

	return ID{value: u}, nil
}

// String returns the UUID string representation
func (id ID) String() string {
	return id.value.String()
}

// UUID returns the underlying uuid.UUID
func (id ID) UUID() uuid.UUID {
	return id.value
}

// IsEmpty checks if the ID is empty (zero value)
func (id ID) IsEmpty() bool {
	return id.value == uuid.Nil
}

// Equals compares two IDs for equality
func (id ID) Equals(other ID) bool {
	return id.value == other.value
}

// Validate checks if the ID is valid
func (id ID) Validate() error {
	if id.IsEmpty() {
		return ErrEmptyID
	}
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (id ID) MarshalJSON() ([]byte, error) {
	if id.IsEmpty() {
		return []byte(`null`), nil
	}
	return []byte(`"` + id.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (id *ID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsed, err := IDFromString(str)
	if err != nil {
		return err
	}

	*id = parsed
	return nil
}
