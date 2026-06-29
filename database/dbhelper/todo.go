package dbhelper

import (
	"Task-Management/database"
	"Task-Management/models"
	"database/sql"
	"errors"
)

func CreateTodo(todo models.TodoRequest, userID string) error {
	args := []interface{}{
		todo.Title,
		todo.Description,
		userID,
	}

	SQL := `INSERT INTO todos(title, description,user_id) 
              VALUES(trim($1),trim($2),$3)`

	_, err := database.DB.Exec(
		SQL,
		args...,
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

func GetTodos(userID, search, completedStatus string) ([]models.Todo, error) {
	args := []interface{}{
		userID,
		search,
		completedStatus,
	}

	query := `SELECT id,user_id,title, description, is_completed
			FROM todos
			WHERE user_id = $1
			  AND ( $2 = '' OR (title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'))
			  AND ($3 = '' OR is_completed = $2::boolean)
			  AND archived_at IS NULL`

	todos := make([]models.Todo, 0)

	err := database.DB.Select(
		&todos,
		query,
		args...,
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
	args := []interface{}{
		todo.Title,
		todo.Description,
		todoID,
		userID,
	}

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
		args...,
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
