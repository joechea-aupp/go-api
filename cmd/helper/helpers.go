package helper

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ServerFeed struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewServerFeed() *ServerFeed {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &ServerFeed{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
}

func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ResponseWithError(w http.ResponseWriter, code int, message string) {
	errorMsg := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}

	ResponseWithJSON(w, code, errorMsg)
}

func ResponseWithHyperMedia(w http.ResponseWriter, code int, htmlString string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	w.Write([]byte(htmlString))
}

type Claims struct {
	Authorized bool   `json:"authorized"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(username, email string) (string, error) {
	claims := &Claims{
		Authorized: true,
		Username:   username,
		Email:      email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
