package server

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ploschka/auth/internal/auth"
	"github.com/ploschka/auth/internal/model"
	"gorm.io/gorm"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if len(guid) == 0 {
		badRequest(w)
		return
	}

	user := model.User{
		Guid: guid,
	}

	db := model.GetDB()

	q := func(d *gorm.DB) *gorm.DB {
		return d.Where(&user).First(&user)
	}

	result := q(db)
	if result.Error != nil {
		badRequest(w)
		return
	}

	ip := r.RemoteAddr
	ip, _, _ = strings.Cut(ip, ":")

	access, refresh, err := auth.GenerateTokens(ip, user)
	if err != nil {
		internalServerError(w)
		return
	}

	refJson, err := json.Marshal(refresh)
	if err != nil {
		internalServerError(w)
		return
	}

	encrypted, err := auth.EncryptToken(refJson)
	if err != nil {
		internalServerError(w)
		return
	}

	hashed, err := auth.HashToken(refJson)
	if err != nil {
		internalServerError(w)
		return
	}

	resp := tokenResponse{
		Access:  access,
		Refresh: base64.RawURLEncoding.EncodeToString(encrypted),
	}

	user.RefreshKey.String = base64.RawURLEncoding.EncodeToString(hashed)
	user.RefreshKey.Valid = true

	respJson, err := json.Marshal(resp)
	if err != nil {
		internalServerError(w)
		return
	}

	_, err = w.Write(respJson)
	if err != nil {
		internalServerError(w)
		return
	}
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {

}
