package routes

import (
	"github.com/ctu-ikz/timetable-be/controllers"
	"github.com/gorilla/mux"
)

func StartTimetableRoutes(router *mux.Router) {
	router.HandleFunc("/timetable", controllers.GetThisWeekTimetable).Methods("GET")
}
