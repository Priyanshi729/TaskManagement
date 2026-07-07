package handler

import (
	"Task-Management/database"
	"Task-Management/database/dbhelper"
	"Task-Management/middleware"
	"Task-Management/models"
	"Task-Management/utils"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.TodoRequest
	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	if err := utils.ParseBody(r, &todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	v := validator.New()
	if err := v.Struct(todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, fmt.Sprintf("invalid validation failed"))
		return
	}

	exists, err := dbhelper.IsTodoExists(todo.Title, userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "database error")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, nil, "todo with title already exists")
		return
	}

	if err := dbhelper.CreateTodo(todo, userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create todo")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"todo created successfully"})
}

func GetTodos(w http.ResponseWriter, r *http.Request) {

	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	todos, err := dbhelper.GetTodos(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Todo []models.Todo `json:"todos"`
	}{todos})
}

func GetTodoById(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")

	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	todo, err := dbhelper.GetTodoByID(userID, todoID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todo")
		return
	}

	if todo == nil {
		utils.RespondJSON(w, http.StatusNotFound, struct {
			Message string `json:"message"`
		}{"todo not found"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")

	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	var todo models.TodoRequest
	if err := utils.ParseBody(r, &todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	err := dbhelper.UpdateTodo(userID, todoID, todo)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update todo")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "todo updated successfully",
	})
}

func MarkTodoAsCompleted(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")

	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	updErr := dbhelper.MarkTodoAsCompleted(todoID, userID)
	if updErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, updErr, "failed to mark todo completed")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "todo marked as completed"})

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")

	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	err := dbhelper.DeleteTodo(todoID, userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete todo")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "todo deleted successfully",
	})
}

func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)

	if err := dbhelper.DeleteAllTodos(database.DB, user.UserID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete all todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "all todos deleted successfully",
	})
}
