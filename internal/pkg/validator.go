package pkg

import "github.com/go-playground/validator/v10"

var (
	v *validator.Validate
)

func InitValidator() {
	v = validator.New()
}

func GetValidator() *validator.Validate {
	return v
}
