package router

import (
	"crud/internal/handler"
	"crud/internal/middleware"
	"net/http"
)

func NewRouter(userHandler *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userHandler.GetUser)
	mux.HandleFunc("POST /users", middleware.RequireValidation(userHandler.RegistrationHandler))

	return mux
}
