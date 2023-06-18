package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func GetAllRepositories(r *mux.Router, db *sqlx.DB, path string) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var repos []models.Repository
		err := db.Select(&repos, "SELECT * FROM repositories")
		if err != nil {
			// handle error
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		jsonResponse, err := json.Marshal(repos)
		if err != nil {
			// handle error
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}).Methods(http.MethodGet) // set the HTTP method to GET
}

func CreateRepository(r *mux.Router, db *sqlx.DB, path string) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var repo models.Repository
	})

}
