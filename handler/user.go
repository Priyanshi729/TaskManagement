package handler

import (
	"Task-Management/database"
	"Task-Management/database/dbhelper"
	"Task-Management/middleware"
	"Task-Management/models"
	"Task-Management/utils"
	"fmt"
	"time"

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

	sessionToken := utils.HashString(user.Email + time.Now().String())
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbhelper.CreateUser(tx, user.Name, user.Email, hashedPassword)
		if saveErr != nil {
			return saveErr
		}

		sessionErr := dbhelper.CreateUserSession(tx, userID, sessionToken)
		if sessionErr != nil {
			return sessionErr
		}
		return nil
	})

	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}

	utils.RespondJSON(w, http.StatusOK,
		struct {
			Message string `json:"message"`
			Token   string `json:"token"`
		}{Message: "user created successfully",
			Token: sessionToken})
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

	sessionToken := utils.HashString(users.Email + time.Now().String())
	if sessionErr := dbhelper.CreateUserSession(database.DB, userID, sessionToken); sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "failed to create user session")
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message      string `json:"message"`
		SessionToken string `json:"session_token"`
	}{Message: "login successful",
		SessionToken: sessionToken})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.UserContext(r)
	if userCtx == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}
	userID := userCtx.ID

	user, getErr := dbhelper.GetUser(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get user")
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-api-key")
	userCtx := middleware.UserContext(r)

	if err := dbhelper.DeleteSessionToken(userCtx.ID, token); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "logout successful",
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	token := r.Header.Get("x-api-key")

	if err := dbhelper.DeleteUser(user.ID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete user account")
		return
	}

	if err := dbhelper.DeleteSessionToken(user.ID, token); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete user session")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "account deleted successfully",
	})
}
