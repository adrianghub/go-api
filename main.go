package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initDB()
	router := mux.NewRouter()

	registerHandlers(router)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}