package data

import (
	"reflect"
	"time"

	sharedtypes "github.com/Jose-Salazar-27/go-university-server/internal/shared/types"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/entity"
	"github.com/Jose-Salazar-27/go-university-server/internal/user/pkg/types"
)

// UserModel maps the domain User entity to the database schema for GORM.
// The name is intentionally generic so callers don't need to know about GORM.
//
// Use FromEntity to populate a UserModel from a domain object; the
// implementation leverages reflection to avoid repetitive field-by-field code.
type UserModel struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"column:email;size:255;unique;not null"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null"`
	FirstName    string    `gorm:"column:first_name;size:100;not null"`
	LastName     string    `gorm:"column:last_name;size:100;not null"`
	UserType     string    `gorm:"column:user_type;size:20;not null"`
	AvatarUrl    *string   `gorm:"column:avatar_url;size:255"`
	IsActive     bool      `gorm:"column:is_active;default:true"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName sets the insert table name for this struct type
func (UserModel) TableName() string { return "users" }

// fromEntity populates a UserModel from a domain entity.User using reflection.
// Special cases (ID, UserType, AvatarUrl) are handled explicitly since their
// types differ between the two structs.
func fromEntity(e *entity.User) *UserModel {
	um := &UserModel{}
	ev := reflect.ValueOf(e).Elem()
	umv := reflect.ValueOf(um).Elem()

	for i := 0; i < ev.NumField(); i++ {
		field := ev.Type().Field(i)
		val := ev.Field(i)
		umField := umv.FieldByName(field.Name)
		if !umField.IsValid() || !umField.CanSet() {
			continue
		}

		switch field.Name {
		case "ID":
			if !val.Interface().(sharedtypes.ID).IsEmpty() {
				umField.SetString(val.Interface().(sharedtypes.ID).String())
			}
		case "UserType":
			umField.SetString(val.Interface().(types.UserType).String())
		case "AvatarUrl":
			str := val.String()
			if str != "" {
				umField.Set(reflect.ValueOf(&str))
			}
		default:
			// identical types can be set directly
			if umField.Type() == val.Type() {
				umField.Set(val)
			}
		}
	}

	return um
}
