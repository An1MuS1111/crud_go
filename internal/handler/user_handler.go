package handler

import (
	"crud/internal/middleware"
	"crud/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")

	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	user, err := h.service.FindByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	fmt.Println(user)
	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// add user with all the required fields
func (h *UserHandler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve validated user from context

	user, ok := r.Context().Value(middleware.UserContextKey).(middleware.UserInput)
	if !ok {
		log.Fatalln("Error: user data not found in the context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Register the user using the service
	registeredUser, err := h.service.Register(r.Context(), user.Name, user.Email)
	if err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response content type
	w.Header().Set("Content-Type", "application/json")
	// Encode the response
	if err := json.NewEncoder(w).Encode(registeredUser); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
