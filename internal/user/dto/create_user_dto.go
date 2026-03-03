package dto

type CreateUserDto struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserType  string `json:"user_type" validate:"required"`
}

type CreateUserResponseDto struct {
	ID string `json:"id"`
}
