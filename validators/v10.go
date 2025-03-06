package validators

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()

	// Register a function to get the field name from the `json` tag
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Add your custom validations here
	// Example: v.RegisterValidation("customtag", customValidationFunc)

	return &CustomValidator{
		validator: v,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Convert validator errors to echo HTTP errors
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := make(map[string]string)

		for _, e := range validationErrors {
			errorMessages[e.Field()] = formatValidationError(e)
		}

		return echo.NewHTTPError(http.StatusBadRequest, errorMessages)
	}
	return nil
}

// formatValidationError formats the validation error message
func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value must be at least " + e.Param()
	case "max":
		return "Value cannot be longer than " + e.Param()
	// Add more cases as needed
	default:
		return "Validation failed on " + e.Tag()
	}
}
