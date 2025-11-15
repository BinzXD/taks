package helper

import (
	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) []string {
	var errors []string

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			errors = append(errors, e.Field()+" is "+e.Tag())
		}
	} else {
		errors = append(errors, err.Error())
	}

	return errors
}
