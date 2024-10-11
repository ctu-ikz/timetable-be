package models

import (
	"sync"
	"time"
)

type SemesterCache struct {
	Semester *Semester
	Start    time.Time
	End      time.Time
}

type TimetableCache struct {
	Data  map[string]WeeklyTimetable
	Mutex sync.RWMutex
}
