package main

import (
	"fmt"
	"net/http"
)

type application struct{}

func main() {
	servePort := "8080"
	app := &application{}

	fmt.Printf("Server is running on port %v", servePort)
	srv := &http.Server{
		Addr:    ":" + servePort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
