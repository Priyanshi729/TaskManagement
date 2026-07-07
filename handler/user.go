package handler

import (
	"Task-Management/database"
	"Task-Management/database/dbhelper"
	"Task-Management/middleware"
	"Task-Management/models"
	"Task-Management/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.RegisterRequest

	if parseErr := utils.ParseBody(r, &user); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	v := validator.New()
	if err := v.Struct(user); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, fmt.Sprintf("invalid validation failed"))
		return
	}

	exists, existsErr := dbhelper.IsUserExists(user.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}

	hashedPassword, hasErr := utils.HashPassword(user.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}

	userID, err := dbhelper.CreateUser(user.Name, user.Email, hashedPassword)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
		return
	}

	sessionID, err := dbhelper.CreateUserSession(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create session")
		return
	}

	token, err := utils.GenerateJWT(userID, sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to generate token")
		return
	}

	utils.RespondJSON(w, http.StatusOK,
		struct {
			Message string `json:"message"`
			Token   string `json:"token"`
		}{Message: "user created successfully",
			Token: token})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var users models.LoginRequest

	if parseErr := utils.ParseBody(r, &users); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	v := validator.New()
	if err := v.Struct(users); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}

	userID, userErr := dbhelper.GetUserID(users)
	if userErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, userErr, "failed to find user")
		return
	}

	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "user not found")
		return
	}

	sessionID, crtErr := dbhelper.CreateUserSession(userID)
	if crtErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, crtErr, "failed to create user session")
		return
	}

	token, err := utils.GenerateJWT(userID, sessionID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to generate token")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"session_token"`
	}{Message: "login successful",
		Token: token})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.UserContext(r)
	userID := userCtx.UserID

	user, getErr := dbhelper.GetUser(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get user")
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.UserContext(r)
	sessionID := userCtx.SessionID

	if delErr := dbhelper.DeleteUserSession(sessionID); delErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, delErr, "failed to delete user session")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"logout successful"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)

	err := database.Tx(func(tx *sqlx.Tx) error {
		if err := dbhelper.DeleteUser(tx, user.UserID); err != nil {
			return err
		}

		if err := dbhelper.DeleteAllTodos(tx, user.UserID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete user account")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "account deleted successfully",
	})
}
