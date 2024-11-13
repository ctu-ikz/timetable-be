package routes

import (
	"github.com/ctu-ikz/timetable-be/controllers"
	"github.com/gorilla/mux"
)

func StartSubjectClassRoutes(router *mux.Router) {
	router.HandleFunc("/subjectclass", controllers.GetCurrentSubjectClass).Methods("GET")
}
