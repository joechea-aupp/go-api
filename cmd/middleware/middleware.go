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
	Web        map[string]string
}

var Feed = &Middleware{
	ServerFeed: helper.NewServerFeed(),
	Web:        map[string]string{},
}

func (mid *Middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Feed.ServerFeed.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (mid *Middleware) LogURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Feed.Web["path"] = r.URL.Path
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

		claims := &helper.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			helper.ResponseWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// get data from claim
		// log.Println("user is: ", claims.Username)
		// log.Println("expired at: ", claims.ExpiresAt)

		next.ServeHTTP(w, r)
	})
}
