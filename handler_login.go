package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mosesbenjamin/sre-from-local-to-prod/internal/auth"
)

func (cfg *config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		Student
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	student, err := cfg.DB.GetStudentByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, student.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	accessToken, err := auth.MakeJWT(student.ID.String(), cfg.JWTSecret, time.Hour, auth.TokenTypeAccess)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	refreshToken, err := auth.MakeJWT(student.ID.String(), cfg.JWTSecret, time.Hour*24*30*6, auth.TokenTypeRefresh)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Student:      databaseStudentToStudent(student),
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

/*
	TODO: Token refresh and revocation
*/
