package main

import (
	"encoding/json"
	"fmt"
	"gobox/rssagg/internal/auth"
	"gobox/rssagg/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Using this as a method to pass db configs
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing the json: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create the user: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUser(user))
}

// Authenticated endpoint, gives user details when provided with apikey
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}
