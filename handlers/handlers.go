package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"educational_api/db"
	"educational_api/models"

	"github.com/gorilla/mux"
)

func getResources(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resources []models.Resource

	rows, err := db.DB.Query("SELECT * FROM resources")
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	defer rows.Close()

	for rows.Next() {
			var res models.Resource
			if err := rows.Scan(&res.ID, &res.Title, &res.Category, &res.Description, &res.URL, &res.DateAdded, &res.ResourceType, &res.CompletionTime); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
			}
			resources = append(resources, res)
	}

	json.NewEncoder(w).Encode(resources)
}

func getResource(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
			http.Error(w, "Invalid resource ID", http.StatusBadRequest)
			return
	}

	var res models.Resource
	row := db.DB.QueryRow("SELECT * FROM resources WHERE id = ?", id)
	err = row.Scan(&res.ID, &res.Title, &res.Category, &res.Description, &res.URL, &res.DateAdded, &res.ResourceType, &res.CompletionTime)
	if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
	} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	json.NewEncoder(w).Encode(res)
}

func createResource(w http.ResponseWriter, r *http.Request) {
	var res models.Resource
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}

	sqlStatement := `
	INSERT INTO resources (title, category, description, url, resource_type, completion_time)
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err = db.DB.Exec(sqlStatement, res.Title, res.Category, res.Description, res.URL, res.ResourceType, res.CompletionTime)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func updateResource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
			http.Error(w, "Invalid resource ID", http.StatusBadRequest)
			return
	}

	var res models.Resource
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}

	sqlStatement := `
	UPDATE resources
	SET title = ?, category = ?, description = ?, url = ?, resource_type = ?, completion_time = ?
	WHERE id = ?`
	_, err = db.DB.Exec(sqlStatement, res.Title, res.Category, res.Description, res.URL, res.ResourceType, res.CompletionTime, id)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	res.ID = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func deleteResource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
			http.Error(w, "Invalid resource ID", http.StatusBadRequest)
			return
	}

	sqlStatement := `DELETE FROM resources WHERE id = ?`
	_, err = db.DB.Exec(sqlStatement, id)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/resources", getResources).Methods("GET")
	router.HandleFunc("/resources/{id}", getResource).Methods("GET")
	router.HandleFunc("/resources", createResource).Methods("POST")
	router.HandleFunc("/resources/{id}", updateResource).Methods("PUT")
	router.HandleFunc("/resources/{id}", deleteResource).Methods("DELETE")
}