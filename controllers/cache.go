package controllers

import "github.com/ctu-ikz/timetable-be/models"

var TimetableCache = models.TimetableCache{
	Data: make(map[string]models.WeeklyTimetable),
}

var SemesterCache = models.SemesterCache{
	Data: make(map[int]models.Semester),
}
