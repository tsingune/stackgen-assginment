package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// ValidatedContextKey is the key used to store validated data in context
	ValidatedContextKey contextKey = "validated"
)

var validate = validator.New()

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// ValidateRequest validates the request body against the provided struct
func ValidateRequest(next http.HandlerFunc, model interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new instance of the model type
		val := model

		// Decode request body
		if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate the struct
		if err := validate.Struct(val); err != nil {
			var errors ValidationErrors

			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, ValidationError{
					Field:   err.Field(),
					Message: getErrorMsg(err),
				})
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": errors,
			})
			return
		}

		// Store validated model in context
		ctx := context.WithValue(r.Context(), ValidatedContextKey, val)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// getErrorMsg returns a human-readable error message for validation errors
func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is below minimum"
	case "max":
		return "Value exceeds maximum"
	case "datetime":
		return "Invalid datetime format"
	default:
		return "Invalid value"
	}
}

// GetValidated retrieves the validated model from the request context
func GetValidated(r *http.Request) interface{} {
	return r.Context().Value(ValidatedContextKey)
}
