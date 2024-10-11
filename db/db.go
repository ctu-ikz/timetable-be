package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ctu-ikz/timetable-be/models"
	"github.com/joho/godotenv"
)

var db *sql.DB

func ConnectToDB() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbSSLMode)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")
	return db, nil
}

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

func GetThisWeekTimetable(time time.Time, class_id string, weeksSinceStart int, semester_id int) (*models.WeeklyTimetable, error) {
	fmt.Println("DB timetable requested")
	var weeklyDBTimetable []models.SubjectClass
	rows, err := db.Query(`select sct.name, su.name, su.shortcut, su.code_name, sc.start_time, sc.end_time, sc.day from "SubjectClass" sc
											join public."Subject" su on sc.subject_id = su.id
											join public."Semester" se on sc.semester_id = se.id
											join public."SubjectClassType" sct on sc.type = sct.id
											join public."SubjectWeek" sw on sc.id = SW.subject_class_id
											where se.id = $1 and sc.class = $2 and sw.week_number=$3`,
		semester_id, class_id, weeksSinceStart)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var subjectClass models.SubjectClass
		err = rows.Scan(&subjectClass.Type, &subjectClass.Name, &subjectClass.Shortcut, &subjectClass.CodeName, &subjectClass.StartTime, &subjectClass.EndTime, &subjectClass.Day)
		if err != nil {
			return nil, err
		}

		weeklyDBTimetable = append(weeklyDBTimetable, subjectClass)
	}

	var WeeklyTimetable models.WeeklyTimetable

	for _, value := range weeklyDBTimetable {
		valueWithoutDay := value
		valueWithoutDay.Day = nil
		if *value.Day == 1 {
			WeeklyTimetable.Monday = append(WeeklyTimetable.Monday, valueWithoutDay)
		} else if *value.Day == 2 {
			WeeklyTimetable.Tuesday = append(WeeklyTimetable.Tuesday, valueWithoutDay)
		} else if *value.Day == 3 {
			WeeklyTimetable.Wednesday = append(WeeklyTimetable.Wednesday, valueWithoutDay)
		} else if *value.Day == 4 {
			WeeklyTimetable.Thursday = append(WeeklyTimetable.Thursday, valueWithoutDay)
		} else if *value.Day == 5 {
			WeeklyTimetable.Friday = append(WeeklyTimetable.Friday, valueWithoutDay)
		}
	}

	return &WeeklyTimetable, nil

}
