package dbhelper

import (
	"Task-Management/database"
	"Task-Management/models"
	"database/sql"
	"errors"
)

func CreateTodo(todo models.TodoRequest) error {
	args := []interface{}{
		todo.Title,
		todo.Description,
	}

	SQL := `INSERT INTO todos(title, description)
              VALUES(trim($1),trim($2))`

	_, err := database.DB.Exec(
		SQL,
		args...,
	)

	return err
}

func IsTodoExists(title string) (bool, error) {
	var exist bool

	query := `SELECT COUNT(id) > 0
			  FROM todos
			  WHERE LOWER(title) = LOWER($1)
			  AND archived_at IS NULL`

	err := database.DB.Get(&exist, query, title)
	return exist, err
}

func GetTodos(search, completedStatus string) ([]models.Todo, error) {
	args := []interface{}{
		search,
		completedStatus,
	}

	query := `SELECT id, title, description, is_completed
			FROM todos
			WHERE ($1 = '' OR (title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'))
			  AND ($2 = '' OR is_completed = $2::boolean)
			  AND archived_at IS NULL`

	todos := make([]models.Todo, 0)

	err := database.DB.Select(
		&todos,
		query,
		args...,
	)

	return todos, err
}

func GetTodoByID(todoId string) (*models.Todo, error) {
	todo := &models.Todo{}

	query := `
		SELECT
			id,
			title,
			description,
			is_completed
		FROM todos
		WHERE id = $1
		AND archived_at IS NULL
	`

	err := database.DB.Get(todo, query, todoId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return todo, err
}

func UpdateTodo(todoID string, todo models.TodoRequest) error {
	args := []interface{}{
		todo.Title,
		todo.Description,
		todoID,
	}

	query := `
		UPDATE todos
		SET
			title = $1,
			description = $2,
		    updated_at = NOW()
		WHERE id = $3
	`

	_, err := database.DB.Exec(
		query,
		args...,
	)

	return err
}

func MarkTodoAsCompleted(todoID string) error {
	SQL := `UPDATE todos
            SET is_completed = true
            WHERE id = $1
              AND archived_at IS NULL
              `

	_, updErr := database.DB.Exec(SQL, todoID)
	return updErr
}

func DeleteTodo(id string) error {
	query := `UPDATE todos
			  SET archived_at = NOW()        
			  WHERE id = $1                              
			    AND archived_at IS NULL
`

	_, err := database.DB.Exec(query, id)
	return err
}
