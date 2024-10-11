package controllers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ctu-ikz/timetable-be/db"
	"github.com/ctu-ikz/timetable-be/models"
)

var (
	semesterCache      *models.SemesterCache
	semesterCacheMutex sync.RWMutex
)

type TimetableCache struct {
	data  map[string]models.WeeklyTimetable
	mutex sync.RWMutex
}

var timetableCache = models.TimetableCache{
	Data: make(map[string]models.WeeklyTimetable),
}

func GetThisWeekTimetable(w http.ResponseWriter, r *http.Request) {
	curTime := time.Now()
	classID := r.URL.Query().Get("class_id")

	if classID == "" {
		http.Error(w, "Missing class_id parameter", http.StatusBadRequest)
		return
	}

	semesterCacheMutex.RLock()
	if semesterCache != nil && curTime.After(semesterCache.Start) && curTime.Before(semesterCache.End) {
		semester := semesterCache.Semester
		semesterCacheMutex.RUnlock()

		timetableCache.Mutex.RLock()
		if timetable, found := timetableCache.Data[classID]; found {
			timetableCache.Mutex.RUnlock()
			json.NewEncoder(w).Encode(timetable)
			return
		}
		timetableCache.Mutex.RUnlock()

		weeksSinceStart := int(curTime.Sub(semester.Start).Hours()/(24*7)) + 1
		timetable, err := db.GetThisWeekTimetable(curTime, classID, weeksSinceStart, semester.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		timetableCache.Mutex.Lock()
		timetableCache.Data[classID] = *timetable
		timetableCache.Mutex.Unlock()

		json.NewEncoder(w).Encode(timetable)
		return
	}
	semesterCacheMutex.RUnlock()

	semester, err := db.GetSemesterByTime(curTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if semester == nil {
		http.Error(w, "No semester found", http.StatusInternalServerError)
		return
	}

	semesterCacheMutex.Lock()
	semesterCache = &models.SemesterCache{
		Semester: semester,
		Start:    semester.Start,
		End:      semester.End,
	}
	semesterCacheMutex.Unlock()

	weeksSinceStart := int(curTime.Sub(semester.Start).Hours()/(24*7)) + 1
	timetable, err := db.GetThisWeekTimetable(curTime, classID, weeksSinceStart, semester.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timetableCache.Mutex.Lock()
	timetableCache.Data[classID] = *timetable
	timetableCache.Mutex.Unlock()

	json.NewEncoder(w).Encode(timetable)
}
