package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ctu-ikz/timetable-be/models"
)

func GetSemesterByTime(time time.Time) (*models.Semester, error) {
	fmt.Println("DB semester requested")

	timeString := time.Format("2006-01-02")

	var semester models.Semester
	err := db.QueryRow(`SELECT id,start,"end",codename FROM "Semester" WHERE $1 BETWEEN "start" AND "end";`, timeString).Scan(&semester.ID, &semester.Start, &semester.End, &semester.Codename)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no semester found")
		}
		return nil, err
	}

	return &semester, nil
}
