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

func PostSemester(semester *models.Semester) (*models.Semester, error) {
	fmt.Println("DB semester put")

	var id int
	err := db.QueryRow(`INSERT INTO "Semester" ("codename", "start", "end")
						VALUES ($1, $2, $3) RETURNING id;`,
		semester.Codename, semester.Start, semester.End).Scan(&id)

	if err != nil {
		return nil, err
	}

	semester.ID = &id

	return semester, nil
}

func DeleteSemester(id int) error {
	fmt.Println("DB semester delete")

	result, err := db.Exec(`DELETE FROM "Semester" WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no semester found with id %d", id)
	}

	return nil
}

func PutSemester(id int, semester *models.Semester) error {
	result, err := db.Exec(`UPDATE "Semester" SET "codename" = $1, "end" = $2, "start" = $3 WHERE id = $4`, semester.Codename, semester.End, semester.Start, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no semester updated with id %d", id)
	}

	return nil
}

func GetSemesterByID(id int) (*models.Semester, error) {
	fmt.Println("DB semester get")
	var semester models.Semester
	err := db.QueryRow(`SELECT id,start,"end",codename FROM "Semester" WHERE id = $1`, id).Scan(&semester.ID, &semester.Start, &semester.End, &semester.Codename)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no semester found with id %d", id)
		}
		return nil, err
	}

	return &semester, nil
}
