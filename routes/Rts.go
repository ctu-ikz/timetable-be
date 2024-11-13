package routes

import (
	"github.com/ctu-ikz/timetable-be/controllers"
	"github.com/gorilla/mux"
)

func StartRoutes(router *mux.Router) {
	StartSemesterRoutes(router)
	StartTimetableRoutes(router)
	StartSubjectClassRoutes(router)
	StartAuthRoutes(router)
	router.HandleFunc("/", controllers.GetIndex).Methods("GET")
}
