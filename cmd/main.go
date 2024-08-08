package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joechea-aupp/go-api/internal/factory"
	"github.com/joechea-aupp/go-api/internal/service"
)

type application struct {
	User *factory.UserService
}

func main() {
	servePort := "8080"

	mongoClient, err := service.ConnectToMongo()
	if err != nil {
		log.Panic(err)
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

	service.New(mongoClient)

	app := &application{
		User: factory.NewUserService(),
	}

	fmt.Printf("Server is running on port %v", servePort)
	srv := &http.Server{
		Addr:    ":" + servePort,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
