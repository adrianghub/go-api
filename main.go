package main

import (
	"fmt"
	"log"
	"net/http"

	"educational_api/auth"
	"educational_api/db"
	"educational_api/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.InitDB()

	router := mux.NewRouter()

	router.HandleFunc("/signup", auth.SignUpHandler).Methods("POST")
	router.HandleFunc("/verify-email", auth.VerifyEmailHandler).Methods("GET")

	handlers.RegisterHandlers(router)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}