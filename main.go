package main

import (
	"log"
	"net/http"

	_ "effectivemobile/docs"
	"effectivemobile/server"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	srv := server.NewServer()

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Starting server on :8080")
	if err := srv.ServeHTTP(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
