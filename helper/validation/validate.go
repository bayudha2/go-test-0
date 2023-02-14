package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate(p interface{}) ([]string, error) {
	validate := validator.New()
	if err := validate.Struct(p); err != nil {

		var errors = []string{}
		if errsObject, ok := err.(validator.ValidationErrors); ok {
			for _, err := range errsObject {
				switch err.Tag() {
				case "required":
					errors = append(errors, fmt.Sprintf("%s is required", err.Field()))
				case "min":
					errors = append(errors, fmt.Sprintf("%s value must greater than %s", err.Field(), err.Param()))
				case "max":
					errors = append(errors, fmt.Sprintf("%s value must less than %s", err.Field(), err.Param()))
				case "email":
					errors = append(errors, fmt.Sprintf("%s must be a email format", err.Field()))
				}
			}
		}
		return errors, err
	}

	return []string{}, nil
}
