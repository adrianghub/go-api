package main

import (
	"fmt"
	"log"
	"net/http"

	"educational_api/db"
	"educational_api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()
	router := mux.NewRouter()

	handlers.RegisterHandlers(router)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}