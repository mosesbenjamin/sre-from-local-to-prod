package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/mosesbenjamin/sre-from-local-to-prod/internal/auth"
	"github.com/mosesbenjamin/sre-from-local-to-prod/internal/database"
)

func (cfg *config) handlerStudentsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	student, err := cfg.DB.CreateStudent(r.Context(), database.CreateStudentParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email:     params.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create student")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseStudentToStudent(student))
}

/*
TODO: Limit to admins only
*/
func (cfg *config) handlerGetStudents(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = specifiedLimit
	}
	/*
		TODO: Allow filtering and pagination
	*/

	students, err := cfg.DB.GetStudents(r.Context(), int32(limit))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get students")
		return
	}

	/*
		TODO: Cache data to improve performance
	*/

	respondWithJSON(w, http.StatusOK, databaseStudentsToStudents(students))
}

func (cfg *config) handlerStudentGet(w http.ResponseWriter, r *http.Request, student database.Student) {
	studentIDStr := chi.URLParam(r, "studentID")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	dbStudent, err := cfg.DB.GetStudentById(r.Context(), studentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get student")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseStudentToStudent(dbStudent))
}

func (cfg *config) handlerStudentDelete(w http.ResponseWriter, r *http.Request, student database.Student) {
	studentIDStr := chi.URLParam(r, "studentID")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	err = cfg.DB.DeleteStudent(r.Context(), studentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete student")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (cfg *config) handlerUpdateStudentPassword(w http.ResponseWriter, r *http.Request, student database.Student) {
	studentIDStr := chi.URLParam(r, "studentID")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	type parameters struct {
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	dbStudent, err := cfg.DB.UpdateStudentPassword(r.Context(), database.UpdateStudentPasswordParams{
		ID:       studentID,
		Password: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update student")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseStudentToStudent(dbStudent))
}
