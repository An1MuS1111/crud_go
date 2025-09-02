package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
)

// contextKeyType is a custom type for context keys to avoid collisions
type contextKeyType string

const (
	userContextKey contextKeyType = "userInput"
)

// UserInput defines the structure for user input data
type UserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RequireValidation is a middleware that validates user input and attaches it to the request context
func RequireValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := UserInput{}

		// Decode JSON from request body
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Validate required fields
		if user.Name == "" {
			http.Error(w, "Error: name is required", http.StatusBadRequest)
			return
		}
		if user.Email == "" {
			http.Error(w, "Error: email is required", http.StatusBadRequest)
			return
		}

		// Validate email format
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(user.Email) {
			http.Error(w, "Error: invalid email format", http.StatusBadRequest)
			return
		}

		// Attach validated user to the request context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next(w, r.WithContext(ctx))
	}
}

// GetUserFromContext retrieves the validated UserInput from the request context
func GetUserFromContext(r *http.Request) (UserInput, bool) {
	user, ok := r.Context().Value(userContextKey).(UserInput)
	return user, ok
}
