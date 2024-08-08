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

type application struct{}

func main() {
	/*
	 When you declare var model models.Model, you're creating an instance of Model but without a pointer. To call the method with a pointer receiver, you need a pointer to an instance of Model.
	*/
	var User factory.User
	// --------------------

	servePort := "8080"
	app := &application{}

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

	err = User.CreateUser(factory.User{
		Username: "jack",
		Email:    "jack@aupp.edu.kh",
		Password: "password",
	})
	if err != nil {
		panic(err)
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
