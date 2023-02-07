package main

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

type ValidationError struct {
	validator.FieldError
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Errors() []string {
	errs := []string{}

	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	validate := validator.New()

	validate.RegisterValidation("sku", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
		matches := re.FindAllString(fl.Field().String(), -1)

		if len(matches) != 1 {
			return false
		}

		return true
	})

	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i).(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}

	var validationErrors []ValidationError

	for _, err := range errs {
		validationErrors = append(validationErrors, ValidationError{
			err.(validator.FieldError),
		})
	}

	return validationErrors
}
