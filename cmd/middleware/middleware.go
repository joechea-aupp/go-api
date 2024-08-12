package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joechea-aupp/go-api/cmd/helper"
)

type Middleware struct {
	ServerFeed *helper.ServerFeed
}

var feed = &Middleware{
	ServerFeed: helper.NewServerFeed(),
}

func (mid *Middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed.ServerFeed.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (mid *Middleware) VerifyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			helper.ResponseWithError(w, http.StatusUnauthorized, "missing token")
			return
		}

		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			helper.ResponseWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
