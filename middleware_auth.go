package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mosesbenjamin/sre-from-local-to-prod/internal/auth"
	"github.com/mosesbenjamin/sre-from-local-to-prod/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.Student)

func (cfg *config) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
			return
		}

		subject, err := auth.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't vaidate JWT")
			return
		}

		studentID, err := uuid.Parse(subject)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse student ID")
			return
		}

		student, err := cfg.DB.GetStudentById(r.Context(), studentID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't get student")
			return
		}

		handler(w, r, student)
	}
}
