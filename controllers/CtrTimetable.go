package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/ctu-ikz/timetable-be/models"
)

func GetThisWeekTimetable(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	classID := r.URL.Query().Get("class_id")

	if classID == "" {
		http.Error(w, "Missing class_id parameter", http.StatusBadRequest)
		return
	}

	SemesterCache.Mutex.RLock()

	var semester *models.Semester

	SemesterCache.Mutex.RUnlock()
	for _, semesterL := range SemesterCache.Data {
		if currentTime.After(semesterL.Start) && currentTime.Before(semesterL.End) {
			semester = &semesterL
			break
		}
	}

	if semester == nil {
		semester, err := db.GetSemesterByTime(currentTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if semester == nil {
			http.Error(w, "No semester found", http.StatusInternalServerError)
			return
		}
		SemesterCache.Mutex.Lock()
		if semester.ID != nil {
			SemesterCache.Data[*semester.ID] = *semester
			SemesterCache.Mutex.Unlock()
			TimetableCache.Mutex.RLock()
			if timetable, found := TimetableCache.Data[classID]; found {
				TimetableCache.Mutex.RUnlock()
				json.NewEncoder(w).Encode(timetable)
				return
			}
			TimetableCache.Mutex.RUnlock()

			weeksSinceStart := int(currentTime.Sub(semester.Start).Hours()/(24*7)) + 1
			var semesterID int
			if semester.ID != nil {
				semesterID = *semester.ID
			}
			timetable, err := db.GetThisWeekTimetable(currentTime, classID, weeksSinceStart, semesterID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			TimetableCache.Mutex.Lock()
			TimetableCache.Data[classID] = *timetable
			TimetableCache.Mutex.Unlock()

			json.NewEncoder(w).Encode(timetable)
		} else {
			SemesterCache.Mutex.Unlock()
			http.Error(w, "Semester ID is nil", http.StatusInternalServerError)
			return
		}
	}

}
