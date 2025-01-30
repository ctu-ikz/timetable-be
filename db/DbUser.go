package db

import (
	"database/sql"
	"errors"

	"github.com/ctu-ikz/timetable-be/helpers"
	"github.com/ctu-ikz/timetable-be/models"
)

func PostUser(user *models.UserDB) (*models.UserDB, error) {
	if user.Password == nil || *user.Password == "" {
		return nil, errors.New("Password cannot be nil or empty")
	}

	id, err := helpers.GenerateSnowflakeID()
	if err != nil {
		return nil, err
	}

	user.ID = &id

	hash, err := helpers.HashPassword(*user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = &hash

	_, err = db.Exec(`INSERT INTO "users" ("id", "username", "password")
					VALUES ($1, $2, $3);`,
		user.ID, user.Username, user.Password)

	if err != nil {
		return nil, err
	}

	user.Password = nil

	return user, nil
}

func GetUserByID(id int64) (*models.UserDB, error) {
	var user models.UserDB
	err := db.QueryRow(`SELECT id,username FROM "users" WHERE id = $1`, id).Scan(&user.ID, &user.Username)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*models.UserDB, error) {
	var user models.UserDB
	err := db.QueryRow(`SELECT id, username, password FROM "users" WHERE username = $1`, username).
		Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserNameTaken(username string) (bool, error) {
	var user models.UserDB
	err := db.QueryRow(`SELECT id,username FROM "users" WHERE username = $1`, username).Scan(&user.ID, &user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
