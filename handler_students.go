package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mosesbenjamin/sre-from-local-to-prod/internal/database"
)

func (cfg *config) handlerStudentsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	student, err := cfg.DB.CreateStudent(r.Context(), database.CreateStudentParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create student")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseStudentToStudent(student))
}

func (cfg *config) handlerGetStudents(w http.ResponseWriter, r *http.Request) {
	students, err := cfg.DB.GetStudents(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get students")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseStudentsToStudents(students))
}
