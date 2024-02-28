package main

import "net/http"

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": http.StatusText(http.StatusOK)})
}
