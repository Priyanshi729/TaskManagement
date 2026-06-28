package dbhelper

import (
	"Task-Management/database"
	"Task-Management/models"
	"Task-Management/utils"
	"time"
)

func IsUserExists(email string) (bool, error) {
	var exists bool
	query := `
		SELECT COUNT(id) > 0
		FROM users
		WHERE email = TRIM($1)
		  AND archived_at IS NULL
	`

	err := database.DB.Get(&exists, query, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CreateUser(name, email, password string) error {
	SQL := `INSERT INTO users (name, email, password)
			  VALUES (TRIM($1), TRIM($2), $3)`

	_, crtErr := database.DB.Exec(SQL, name, email, password)
	return crtErr
}

func CreateUserSession(userID string) (string, error) {
	var sessionID string
	SQL := `INSERT INTO user_session(user_id) 
             VALUES ($1) RETURNING id`
	crtErr := database.DB.Get(&sessionID, SQL, userID)
	return sessionID, crtErr
}

func GetUserID(users models.LoginRequest) (string, error) {
	SQL := `SELECT u.id,
      			   u.password
			  FROM users u
			  WHERE u.email = TRIM($1)
			    AND u.archived_at IS NULL`

	var user models.LoginData
	if getErr := database.DB.Get(&user, SQL, users.Email); getErr != nil {
		return "", getErr
	}
	if passwordErr := utils.CheckPassword(users.Password, user.PasswordHash); passwordErr != nil {
		return "", passwordErr
	}
	return user.ID, nil
}

func GetUser(userID string) (models.User, error) {
	var user models.User
	SQL := `SELECT id, name, email 
              FROM users 
              WHERE id = $1
                AND archived_at IS NULL`

	getErr := database.DB.Get(&user, SQL, userID)
	return user, getErr
}

func GetArchivedAt(sessionID string) (*time.Time, error) {
	var archivedAt *time.Time

	SQL := `SELECT archived_at 
              FROM user_session 
              WHERE id = $1`

	getErr := database.DB.Get(&archivedAt, SQL, sessionID)
	return archivedAt, getErr
}

func DeleteUserSession(sessionID string) error {
	SQL := `UPDATE user_session
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, delErr := database.DB.Exec(SQL, sessionID)
	return delErr
}

func DeleteUser(userID string) error {
	SQL := `UPDATE users
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, delErr := database.DB.Exec(SQL, userID)
	return delErr
}
