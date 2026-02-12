package domain

// UserType represents the type of user in the system
type UserType string

const (
	UserTypeStudent   UserType = "student"
	UserTypeProfessor UserType = "professor"
	UserTypeAdmin     UserType = "admin"
)

// IsValid checks if the user type is valid
func (ut UserType) IsValid() bool {
	switch ut {
	case UserTypeStudent, UserTypeProfessor, UserTypeAdmin:
		return true
	default:
		return false
	}
}

// String returns the string representation of UserType
func (ut UserType) String() string {
	return string(ut)
}

// ToUserType converts a string to UserType, returns UserTypeStudent as default
func ToUserType(s string) UserType {
	switch s {
	case "student":
		return UserTypeStudent
	case "professor":
		return UserTypeProfessor
	case "admin":
		return UserTypeAdmin
	default:
		return UserTypeStudent // default value
	}
}
