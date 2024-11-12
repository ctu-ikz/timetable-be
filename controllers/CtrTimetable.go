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
	for _, semesterL := range SemesterCache.Data {
		if currentTime.After(semesterL.Start) && currentTime.Before(semesterL.End) {
			semester = &semesterL
			break
		}
	}
	SemesterCache.Mutex.RUnlock()

	if semester == nil {
		var err error
		semester, err = db.GetSemesterByTime(currentTime)
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
		}
		SemesterCache.Mutex.Unlock()
	}

	TimetableCache.Mutex.RLock()
	if timetable, found := TimetableCache.Data[classID]; found {
		TimetableCache.Mutex.RUnlock()
		if err := json.NewEncoder(w).Encode(timetable); err != nil {
			http.Error(w, "Failed to encode timetable", http.StatusInternalServerError)
		}
		return
	}
	TimetableCache.Mutex.RUnlock()

	weeksSinceStart := int(currentTime.Sub(semester.Start).Hours()/(24*7)) + 1
	var semesterID int64
	if semester.ID != nil {
		semesterID = *semester.ID
	}

	timetable, err := db.GetThisWeekTimetable(currentTime, classID, weeksSinceStart, int(semesterID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	TimetableCache.Mutex.Lock()
	TimetableCache.Data[classID] = *timetable
	TimetableCache.Mutex.Unlock()

	if err := json.NewEncoder(w).Encode(timetable); err != nil {
		http.Error(w, "Failed to encode timetable", http.StatusInternalServerError)
	}
}
