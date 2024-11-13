package routes

import (
	"github.com/ctu-ikz/timetable-be/controllers"
	"github.com/gorilla/mux"
)

func StartAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", controllers.PostUser).Methods("POST")
	router.HandleFunc("/auth/user/{id}", controllers.GetUserByID).Methods("GET")
	router.HandleFunc("/auth/login", controllers.LoginUser).Methods("POST")
}
