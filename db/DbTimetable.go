package db

import (
	"fmt"
	"time"

	"github.com/ctu-ikz/timetable-be/models"
)

func GetThisWeekTimetable(time time.Time, class_id string, weeksSinceStart int, semester_id int) (*models.WeeklyTimetable, error) {
	fmt.Println("DB timetable requested")
	var weeklyDBTimetable []models.SubjectClass
	rows, err := db.Query(`select sct.name, su.name, su.shortcut, su.code_name, sc.start_time, sc.end_time, sc.day from "SubjectClass" sc
											join public."Subject" su on sc.subject_id = su.id
											join public."Semester" se on sc.semester_id = se.id
											join public."SubjectClassType" sct on sc.type_id = sct.id
											join public."SubjectWeek" sw on sc.id = SW.subject_class_id
											where se.id = $1 and sc.class_id = $2 and sw.week_number=$3`,
		semester_id, class_id, weeksSinceStart)
	fmt.Println(`select sct.name, su.name, su.shortcut, su.code_name, sc.start_time, sc.end_time, sc.day from "SubjectClass" sc
											join public."Subject" su on sc.subject_id = su.id
											join public."Semester" se on sc.semester_id = se.id
											join public."SubjectClassType" sct on sc.type_id = sct.id
											join public."SubjectWeek" sw on sc.id = SW.subject_class_id
											where se.id = $1 and sc.class_id = $2 and sw.week_number=$3`,
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

	fmt.Println(weeklyDBTimetable)
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
