package db

import (
	"fmt"
	"time"

	"github.com/ctu-ikz/timetable-be/models"
)

func CheckRefreshTokenValidity(hashedToken string, userID int64) (bool, error) {
	var refreshToken models.RefreshTokenDB

	rows, err := db.Query(`SELECT id, token, user_id FROM refresh_tokens WHERE is_valid = true`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&refreshToken.ID, &refreshToken.Token, &refreshToken.UserID); err != nil {
			return false, err
		}

		if hashedToken == refreshToken.Token {
			if refreshToken.UserID == userID {
				return true, nil
			}
		}
	}

	if err := rows.Err(); err != nil {
		return false, err
	}

	return false, fmt.Errorf("no valid refresh token found for user")
}

func CreateRefreshToken(userID int64, token string, validUntil time.Time) error {
	_, err := db.Exec(`INSERT INTO refresh_tokens (token, user_id, expires_at, created_at, is_valid) VALUES ($1, $2, $3, $4, $5);`, token, userID, validUntil, time.Now(), true)
	if err != nil {
		return err
	}

	return nil
}

func InvalidateRefreshToken(hashedToken string) error {
	_, err := db.Exec(`UPDATE refresh_tokens SET is_valid = false WHERE token = $1;`, hashedToken)
	if err != nil {
		return err
	}

	return nil
}
