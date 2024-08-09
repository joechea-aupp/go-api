package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joechea-aupp/go-api/internal/db"
)

type application struct {
	User     *db.UserService
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	servePort := "8080"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mongoClient, err := db.ConnectToMongo()
	if err != nil {
		errorLog.Println(err)
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
		User:     db.NewUserService(),
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	infoLog.Printf("Server is running on port %v", servePort)
	srv := &http.Server{
		Addr:    ":" + servePort,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Println(err)
	}
}
