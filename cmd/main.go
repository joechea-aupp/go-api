package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/joechea-aupp/go-api/cmd/api"
	"github.com/joechea-aupp/go-api/cmd/helper"
	"github.com/joechea-aupp/go-api/cmd/middleware"
	"github.com/joechea-aupp/go-api/cmd/web"
	"github.com/joechea-aupp/go-api/internal/db"
	"github.com/joho/godotenv"
)

type application struct {
	User       *db.UserService
	Api        *api.Api
	Web        *web.Web
	Middleware *middleware.Middleware
}

func main() {
	godotenv.Load(".env")
	servePort := os.Getenv("PORT")
	feed := helper.NewServerFeed()

	mongoClient, err := db.ConnectToMongo()
	if err != nil {
		feed.ErrorLog.Println(err)
	}

	// create a context to timeout mongodb connection
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		// the timeout is time on the disconnect function, it has 15 seconds to disconnect
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db.New(mongoClient)

	app := &application{
		User: db.NewUserService(),
	}

	feed.InfoLog.Printf("Server is running on port %v", servePort)
	srv := &http.Server{
		Addr:    ":" + servePort,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		feed.ErrorLog.Println(err)
	}
}
