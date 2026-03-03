package types

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

// ID is a value object representing a UUID for entities.
type ID struct {
	value string
}

// NewID generates a new UUID-based UserID.
func NewID() (ID, error) {
	id := uuid.New().String()
	if id == "" {
		return ID{}, errors.New("failed to generate user id")
	}
	return ID{value: id}, nil
}

// ParseID validates and parses a string into a UserID.
func ParseID(id string) (ID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return ID{}, errors.New("invalid user id format")
	}
	return ID{value: id}, nil
}

func MustGetID(id string) ID {
	return ID{id}
}

func (id ID) String() string { return id.value }

func (id ID) Equals(other ID) bool { return id.value == other.value }

func (id ID) IsEmpty() bool { return id.value == "" }

func (id ID) MarshalJSON() ([]byte, error) { return json.Marshal(id.value) }

func (id *ID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsedID, err := ParseID(str)
	if err != nil {
		return err
	}

	*id = parsedID
	return nil
}
