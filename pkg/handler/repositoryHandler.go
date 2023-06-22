package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/config"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/models"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetAllRepositories(r *mux.Router, db *sqlx.DB, path string) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var repos []models.Repository
		err := db.Select(&repos, "SELECT * FROM repositories")
		if err != nil {
			// handle error
			utils.WriteError(w, http.StatusInternalServerError, "Error fetching repositories.")
			return
		}

		jsonResponse, err := json.Marshal(repos)
		if err != nil {
			// handle error
			utils.WriteError(w, http.StatusUnprocessableEntity, "Error marshalling repository into struct.")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}).Methods(http.MethodGet) // set the HTTP method to GET
}

func GetRepositoryById(r *mux.Router, db *sqlx.DB, path string) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var repo models.Repository
		vars := mux.Vars(r)
		repositoryId := vars["repository_id"]
		ctx := r.Context()

		query := "SELECT * FROM repositories WHERE id = $1"
		err := db.Get(&repo, query, repositoryId)

		if err != nil {
			if err == sql.ErrNoRows {
				logger.Log.WithField("traceId", utils.GetTraceId(ctx)).Error(fmt.Sprintf("Repository with id %s not found", repositoryId))
				utils.WriteError(w, http.StatusNotFound, fmt.Sprintf("Repository with id %s not found", repositoryId))
			} else {
				logger.Log.WithField("traceId", utils.GetTraceId(ctx)).Error(fmt.Sprintf("Error fetching repository with id %s. Error: %s", repositoryId, err))
				utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching repository with id %s.", repositoryId))
			}
			return
		}
		jsonResponse, err := json.Marshal(repo)
		if err != nil {
			// handle error
			utils.WriteError(w, http.StatusUnprocessableEntity, "Error marshalling repository into struct.")
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

		configId := repo.ConfigurationID

		repo.ConfigurationID, err = fetchConfigurationId(r.Context(), configId)
		if err != nil {
			utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Sprintf("Error getting configurationId %s from the configuration service", configId))
			logger.Log.Error(err)
			return
		}

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

func DeleteRepository(r *mux.Router, db *sqlx.DB, path string) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		repositoryId, ok := vars["repository_id"]

		ctx := r.Context()

		if !ok {
			logger.Log.WithField("traceId", utils.GetTraceId(ctx)).Error(fmt.Sprintf("RepositoryId %s not sent properly.", repositoryId))
			utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("Error receiving repository_id variable. Variable value: %s", repositoryId))
			return
		}

		query := `DELETE FROM repositories WHERE id = $1`

		result, err := db.ExecContext(ctx, query, repositoryId)
		if err != nil {
			logger.Log.WithField("traceId", utils.GetTraceId(ctx)).Error(fmt.Sprintf("Error deleting repository %s. Error: %s", repositoryId, err))
			utils.WriteError(w, http.StatusNotFound, fmt.Sprintf("Error deleting repository with ID %s.", repositoryId))
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting repository with ID %s.", repositoryId))
			logger.Log.WithField("traceId", utils.GetTraceId(ctx)).Error(fmt.Sprintf("Error deleting repository %s. Error: %s", repositoryId, err))
			return
		}

		if rowsAffected == 0 {
			logger.Log.WithField("traceId", utils.GetTraceId(ctx)).Error(fmt.Sprintf("Repository with id %s not found. Error: %s", repositoryId, err))
			utils.WriteError(w, http.StatusNotFound, fmt.Sprintf("Repository with id %s not found.", repositoryId))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

	}).Methods(http.MethodDelete)
}

func fetchConfigurationId(ctx context.Context, configID string) (string, error) {
	// Generate the URL for service B
	url := fmt.Sprintf("%s/%s/%s", config.Config.ConfigurationServiceURL, "configuration", configID)
	logger.Log.WithField("traceId", utils.GetTraceId(ctx)).Info("Requesting configuration service...")
	// Create a new request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request for configuration service: %w", err)
	}

	// Inject the current span context into the headers of the outbound request
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Send GET request to service B
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error communicating with configuration service. Error: %s", err)
	}
	if response.StatusCode != 200 {
		// handle error
		return "", fmt.Errorf("error communicating with configuration service. Status Code: %s", strconv.Itoa(response.StatusCode))
	}
	defer response.Body.Close()

	// Decode the response
	var result map[string]string
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		// handle error
		return "", fmt.Errorf("error decoding response from service B: %w", err)
	}

	logger.Log.WithField("traceId", utils.GetTraceId(ctx)).Info("Configuration service called successfully.")

	// Extract the configuration value
	configuration := result["configuration"]

	return configuration, nil
}

func generateRepoId() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("Failed to generate UUID: %v", err)
		return "", err
	}
	return u.String(), nil
}
