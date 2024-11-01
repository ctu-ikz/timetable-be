package db

import (
	"fmt"
	"time"

	"github.com/ctu-ikz/timetable-be/models"
)

func GetCurrentSubjectClass(semester models.Semester, classID int, currentTime time.Time) (*models.SubjectClass, error) {
	fmt.Println("DB current subject class requested")

	var subjectClass models.SubjectClass

	weekNumber := int(currentTime.Sub(semester.Start).Hours()/(24*7)) + 1

	err := db.QueryRow(`
		SELECT sc.end_time, sc.start_time, s.code_name, s.name, s.shortcut
		FROM "SubjectClass" sc
		JOIN public."Subject" s on s.id = sc.subject_id
		JOIN public."SubjectWeek" sw on sc.id = sw.subject_class_id
		WHERE start_time <= $1::TIME
			AND sc.end_time >= $1::TIME
			AND sc.semester_id = $2
			AND sc.class_id = $3
			AND sw.week_number = $4
	`, currentTime.Format("15:04"), semester.ID, classID, weekNumber,
	).Scan(&subjectClass.EndTime, &subjectClass.StartTime, &subjectClass.CodeName, &subjectClass.Name, &subjectClass.Shortcut)

	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, subjectClass.EndTime)
	if err != nil {
		return nil, err
	}

	startTime, err := time.Parse(time.RFC3339, subjectClass.StartTime)
	if err != nil {
		return nil, err
	}

	subjectClass.StartTime = startTime.Format("15:04")
	subjectClass.EndTime = endTime.Format("15:04")

	return &subjectClass, nil
}
