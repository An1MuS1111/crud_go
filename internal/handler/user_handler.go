package handler

import (
	"crud/internal/service"
	"encoding/json"
	"fmt"
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
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if createUser.Name == "" || createUser.Email == "" {
		http.Error(w, "name and email can't be empty", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), createUser.Name, createUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
