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
	semesterCache      *SemesterCache
	semesterCacheMutex sync.RWMutex
)

type TimetableCache struct {
	data  map[string]models.WeeklyTimetable
	mutex sync.RWMutex
}

var timetableCache = TimetableCache{
	data: make(map[string]models.WeeklyTimetable),
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

		timetableCache.mutex.RLock()
		if timetable, found := timetableCache.data[classID]; found {
			timetableCache.mutex.RUnlock()
			json.NewEncoder(w).Encode(timetable)
			return
		}
		timetableCache.mutex.RUnlock()

		weeksSinceStart := int(curTime.Sub(semester.Start).Hours()/(24*7)) + 1
		timetable, err := db.GetThisWeekTimetable(curTime, classID, weeksSinceStart, semester.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		timetableCache.mutex.Lock()
		timetableCache.data[classID] = *timetable
		timetableCache.mutex.Unlock()

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
	semesterCache = &SemesterCache{
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

	timetableCache.mutex.Lock()
	timetableCache.data[classID] = *timetable
	timetableCache.mutex.Unlock()

	json.NewEncoder(w).Encode(timetable)
}
