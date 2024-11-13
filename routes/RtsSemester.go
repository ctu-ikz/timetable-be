package routes

import (
	"github.com/ctu-ikz/timetable-be/controllers"
	"github.com/gorilla/mux"
)

func StartSemesterRoutes(router *mux.Router) {
	router.HandleFunc("/semester", controllers.GetSemester).Methods("GET")
	router.HandleFunc("/semester", controllers.PostSemester).Methods("POST")
	router.HandleFunc("/semester/{id}", controllers.DeleteSemester).Methods("DELETE")
	router.HandleFunc("/semester/{id}", controllers.PutSemester).Methods("PUT")
	router.HandleFunc("/semester/{id}", controllers.GetSemesterByID).Methods("GET")
}
