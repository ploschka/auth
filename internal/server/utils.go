package server

import (
	"net/http"

	log "github.com/ploschka/auth/internal/logger"
)

type tokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func badRequest(w http.ResponseWriter, err error) {
	log.Error(err)
	log.Info(http.StatusBadRequest)
	w.WriteHeader(http.StatusBadRequest)
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Error(err)
	log.Info(http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
}

func unauthorized(w http.ResponseWriter, err error) {
	log.Error(err)
	log.Info(http.StatusUnauthorized)
	w.WriteHeader(http.StatusUnauthorized)
}
