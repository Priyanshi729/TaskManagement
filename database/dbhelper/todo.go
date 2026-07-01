package dbhelper

import (
	"Task-Management/database"
	"Task-Management/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

func CreateTodo(todo models.TodoRequest, userID string) error {

	SQL := `INSERT INTO todos(title, description,user_id) 
              VALUES(trim($1),trim($2),$3)`

	_, err := database.DB.Exec(
		SQL, todo.Title, todo.Description, userID,
	)

	return err
}

func IsTodoExists(title string, userID string) (bool, error) {
	query := `SELECT COUNT(id) > 0
			  FROM todos
			  WHERE LOWER(title) = LOWER($1)
			 AND user_id = $2     
			  AND archived_at IS NULL`

	var exist bool
	err := database.DB.Get(&exist, query, title, userID)
	return exist, err
}

func GetTodos(userID string) ([]models.Todo, error) {

	query := `SELECT id, user_id, title, description, is_completed
		FROM todos
		WHERE user_id = $1
		  AND archived_at IS NULL`

	todos := make([]models.Todo, 0)

	err := database.DB.Select(
		&todos,
		query,
		userID,
	)

	return todos, err
}

func GetTodoByID(userID, todoId string) (*models.Todo, error) {
	todo := &models.Todo{}

	query := `
		SELECT
			id,
			user_id,
			title,
			description,
			is_completed
		FROM todos
		WHERE id = $1
		AND user_id = $2
		AND archived_at IS NULL
	`

	err := database.DB.Get(todo, query, todoId, userID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return todo, err
}

func UpdateTodo(userID, todoID string, todo models.TodoRequest) error {
	query := `
		UPDATE todos
		SET
			title = $1,
			description = $2
		WHERE id = $3
		AND user_id = $4
		AND archived_at IS NULL
	`

	_, err := database.DB.Exec(
		query,
		todo.Title,
		todo.Description,
		todoID,
		userID,
	)

	return err
}

func MarkTodoAsCompleted(todoID string, userID string) error {
	SQL := `UPDATE todos
            SET is_completed = true
            WHERE id = $1
              AND user_id = $2    
              AND archived_at IS NULL
              `

	_, updErr := database.DB.Exec(SQL, todoID, userID)
	return updErr
}

func DeleteTodo(id string, userID string) error {
	query := `UPDATE todos
			  SET archived_at = NOW()        
			  WHERE id = $1  
			    AND user_id = $2
			    AND archived_at IS NULL
`

	_, err := database.DB.Exec(query, id, userID)
	return err
}

func DeleteAllTodos(db sqlx.Ext, userID string) error {
	SQL := `UPDATE todos
              SET archived_at = NOW()        
              WHERE user_id = $1             
                AND archived_at IS NULL`

	_, delErr := db.Exec(SQL, userID)
	return delErr
}
