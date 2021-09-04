package validate

import (
	"context"

	"github.com/go-playground/validator/v10"
)

// RegisterValidationFn type allows the register of custom validators.
type RegisterValidationFn func(*validator.Validate)

var v *validator.Validate

// Setup creates the validate and calls the given functions to register additional validators.
func Setup(registerValidations ...RegisterValidationFn) {
	v = validator.New()

	for _, fn := range registerValidations {
		fn(v)
	}
}

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func Struct(s interface{}) error {
	return v.Struct(s)
}

// Var validates a single variable using tag style validation.
func Var(ctx context.Context, field interface{}, tag string) error {
	return v.VarCtx(ctx, field, tag)
}
