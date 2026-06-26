package handler

import (
	"Task-Management/database/dbhelper"
	"Task-Management/models"
	"Task-Management/utils"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.TodoRequest

	if err := utils.ParseBody(r, &todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	exists, err := dbhelper.IsTodoExists(todo.Title)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "database error")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, nil, "todo with title already exists")
		return
	}

	if err := dbhelper.CreateTodo(todo); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create todo")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"todo created successfully"})
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	completedStatus := r.URL.Query().Get("completedStatus")

	todos, err := dbhelper.GetTodos(search, completedStatus)
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
	fmt.Printf("todoID = %q\n", todoID)

	// todo add created
	todo, err := dbhelper.GetTodoByID(todoID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todo")
		return
	}

	if todo == nil {
		utils.RespondJSON(w, http.StatusNotFound, struct {
			Message string `json:"message"`
		}{"todo not found"})
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")

	var todo models.TodoRequest
	if err := utils.ParseBody(r, &todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	err := dbhelper.UpdateTodo(todoID, todo)
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

	updErr := dbhelper.MarkTodoAsCompleted(todoID)
	if updErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, updErr, "failed to mark todo completed")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		"todo marked as completed"})

}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")

	err := dbhelper.DeleteTodo(todoID)
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
