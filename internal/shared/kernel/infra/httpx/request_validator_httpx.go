package httpx

import "github.com/go-playground/validator/v10"

// https://docs.gofiber.io/guide/validation
type structValidator struct {
	validate *validator.Validate
}

func NewRequestValidator() *structValidator {
	return &structValidator{validate: validator.New()}
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}
