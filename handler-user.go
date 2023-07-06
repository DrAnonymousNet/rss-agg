package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dranonymousnet/rss-agg/internal/database"
	"github.com/google/uuid"
)



func (apicfg *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, fmt.Sprintf("Error Parsing Json %v", err ))
		return
	}
	user, err := apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})
	if err != nil{
		respondWithError(w, 500, fmt.Sprintf("Unable to Create User %v", err))
	}
	respondWithJson(w, 201, dataBaseUserToUser(user))
}
