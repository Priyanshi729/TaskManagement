package utils

import (
	"Task-Management/models"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ParseBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			log.Printf("Failed to respond JSON with error: %+v", err)
		}
	}
}

func newClientError(err error, statusCode int, messageToUser string) *models.ClientError {
	clientErr := ""
	if err != nil {
		clientErr = err.Error()
	}

	return &models.ClientError{
		MessageToUser: messageToUser,
		Err:           clientErr,
		StatusCode:    statusCode,
	}
}

func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string, additionalInfoForDevs ...string) {
	log.Printf("status: %d, message: %s, err: %+v ", statusCode, messageToUser, err)
	clientError := newClientError(err, statusCode, messageToUser)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		log.Printf("status: %d, message: %s, err: %+v ", statusCode, messageToUser, err)
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(userID, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"userId":    userID,
		"sessionId": sessionID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
