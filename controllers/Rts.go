package controllers

import (
	"net/http"

	"github.com/ctu-ikz/timetable-be/services"
	"github.com/gorilla/mux"
)

func StartRoutes(router *mux.Router) {
	StartAuthRoutes(router)
	router.HandleFunc("/", services.GetIndex).Methods("GET")
	router.Handle("/ping", JWTAuthMiddleware(http.HandlerFunc(Ping))).Methods("GET")
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
