package router

import (
	"crud/internal/handler"
	"net/http"
)

func NewRouter(userHandler *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userHandler.GetUser)
	mux.HandleFunc("POST /users", userHandler.CreateUser)

	return mux
}
