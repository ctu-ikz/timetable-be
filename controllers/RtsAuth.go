package controllers

import (
	"github.com/ctu-ikz/timetable-be/services"
	"github.com/gorilla/mux"
)

func StartAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", services.PostUser).Methods("POST")
	router.HandleFunc("/auth/user/{id}", services.GetUserByID).Methods("GET")
	router.HandleFunc("/auth/login", services.LoginUser).Methods("POST")
	router.HandleFunc("/auth/refresh", services.RefreshTokens).Methods("POST")
}
