package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/models"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

		err := json.NewDecoder(r.Body).Decode(&repo)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid request body.")
			logger.Log.Error("Error decoding repository object: %s", err)
			return
		}

		validate := validator.New()
		err = validate.Struct(repo)
		if err != nil {
			var errors []string
			for _, fieldErr := range err.(validator.ValidationErrors) {
				field := fieldErr.Field()
				tag := fieldErr.Tag()
				errors = append(errors, fmt.Sprintf("%s failed on %s validation", field, tag))
			}
			utils.WriteError(w, http.StatusUnprocessableEntity, strings.Join(errors, ", "))
			logger.Log.Error("Validation failed for repository object: %s", err)
			return
		}

		repo.ID, err = generateRepoId()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Internal Server Error.")
			logger.Log.Error("Error generating UUID for RepositoryID: %s", err)
			return
		}

		repo.CreationDate = time.Now()

		insertQuery := `INSERT INTO repositories (id, name, owner, creationdate, configurationid) VALUES ($1, $2, $3, $4, $5)`

		_, err = db.Exec(insertQuery, repo.ID, repo.Name, repo.Owner, repo.CreationDate, repo.ConfigurationID)
		if err != nil {
			logger.Log.Error("Error insertind repository in database: %s", err)
		}

		jsonResponse, err := json.Marshal(repo)
		if err != nil {
			// handle error
			fmt.Fprintf(w, "Error: %s", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}).Methods(http.MethodPost)

}

func generateRepoId() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("Failed to generate UUID: %v", err)
		return "", err
	}
	return u.String(), nil
}
