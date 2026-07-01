package dbhelper

import (
	"Task-Management/database"
	"Task-Management/models"
	"Task-Management/utils"

	"github.com/jmoiron/sqlx"
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

func CreateUser(db sqlx.Ext, name, email, password string) (string, error) {
	var userID string

	SQL := `
		INSERT INTO users (name, email, password)
		VALUES (TRIM($1), TRIM($2), $3)
		RETURNING id
	`

	err := db.QueryRowx(SQL, name, email, password).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func CreateUserSession(db sqlx.Ext, userID, sessionToken string) error {
	SQL := `INSERT INTO user_session(user_id,session_Token) 
             VALUES ($1,$2) RETURNING id`
	_, crtErr := db.Exec(SQL, userID, sessionToken)
	return crtErr
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

func GetUserBySession(token string) (*models.User, error) {
	var user models.User

	SQL := `
		SELECT u.id,
		       u.name,
		       u.email
		FROM users u
		JOIN user_session us
		  ON us.user_id = u.id
		WHERE us.session_token = $1
		  AND us.archived_at IS NULL
		  AND u.archived_at IS NULL
	`

	err := database.DB.Get(&user, SQL, token)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteSessionToken(userID, token string) error {
	SQL := `
		UPDATE user_session
		SET archived_at = NOW()
		WHERE user_id = $1
		  AND session_token = $2
		  AND archived_at IS NULL
	`

	_, err := database.DB.Exec(SQL, userID, token)
	return err
}

func DeleteUser(userID string) error {
	SQL := `UPDATE users
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, delErr := database.DB.Exec(SQL, userID)
	return delErr
}
