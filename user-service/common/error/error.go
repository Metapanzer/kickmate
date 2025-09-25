package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// ValidationResponse represents a single field validation error
type ValidationResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

// ErrValidator is a map to store validation tags
var ErrValidator = map[string]string{}

// ErrValidationResponse converts validation errors into a slice of ValidationResponse to show each field error
func ErrValidationResponse(err error) (validationResponses []ValidationResponse) {
	var fieldErr validator.ValidationErrors
	if errors.As(err, &fieldErr) {
		for _, err := range fieldErr {
			switch err.Tag() {
			case "required":
				validationResponses = append(validationResponses, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is required", err.Field()),
				})
			case "email":
				validationResponses = append(validationResponses, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a valid email address", err.Field()),
				})
			default:
				errValidator, ok := ErrValidator[err.Tag()]
				if ok {
					count := strings.Count(errValidator, "%s")
					if count == 1 {
						validationResponses = append(validationResponses, ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field()),
						})
					} else {
						validationResponses = append(validationResponses, ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field(), err.Param()),
						})
					}
				} else {
					validationResponses = append(validationResponses, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("something wrong on %s: %s", err.Field(), err.Tag()),
					})
				}
			}
		}
	}
	return validationResponses
}

// WrapError logs the error and returns it
func WrapError(err error) error {
	logrus.Errorf("error: %v", err)
	return err
}
