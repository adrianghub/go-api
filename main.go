package main

import (
	"fmt"
	"log"
	"net/http"

	"educational_api/auth"
	"educational_api/db"
	"educational_api/resources"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.InitDB()

	router := mux.NewRouter()

	auth.RegisterHandlers(router)
	resources.RegisterHandlers(router)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}