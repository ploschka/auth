package server

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/ploschka/auth/internal/auth"
	"github.com/ploschka/auth/internal/mailer"
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

	q2 := func(d *gorm.DB) *gorm.DB {
		return d.Save(&user)
	}

	result = q2(db)
	if result.Error != nil {
		internalServerError(w)
		return
	}

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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		internalServerError(w)
		return
	}

	req := tokenResponse{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		badRequest(w)
		return
	}

	rawRefresh, err := base64.RawURLEncoding.DecodeString(req.Refresh)
	if err != nil {
		badRequest(w)
		return
	}

	rawRefresh, err = auth.DecryptToken(rawRefresh)
	if err != nil {
		badRequest(w)
		return
	}

	refresh := &auth.RefreshToken{}
	err = json.Unmarshal(rawRefresh, refresh)

	valid := auth.CheckPair(req.Access, refresh)
	if !valid {
		unauthorized(w)
		return
	}

	ip := r.RemoteAddr
	ip, _, _ = strings.Cut(ip, ":")

	user := model.User{
		Guid: refresh.GUID,
	}

	db := model.GetDB()

	q := func(d *gorm.DB) *gorm.DB {
		return d.Where(&user).First(&user)
	}

	result := q(db)
	if result.Error != nil {
		unauthorized(w)
		return
	}

	if !user.RefreshKey.Valid {
		unauthorized(w)
		return
	}

	rawHash, err := base64.RawURLEncoding.DecodeString(user.RefreshKey.String)
	if err != nil {
		internalServerError(w)
		return
	}

	valid = auth.Validate(rawRefresh, rawHash)
	if !valid {
		unauthorized(w)
		return
	}

	if ip != refresh.Ip {
		err = mailer.SendIpWarning(user.Email, ip)
		if err != nil {
			internalServerError(w)
			return
		}
	}

	var access string
	access, refresh, err = auth.GenerateTokens(ip, user)
	if err != nil {
		internalServerError(w)
		return
	}

	rawRefresh, err = json.Marshal(refresh)
	if err != nil {
		internalServerError(w)
		return
	}

	encrypted, err := auth.EncryptToken(rawRefresh)
	if err != nil {
		internalServerError(w)
		return
	}

	hashed, err := auth.HashToken(rawRefresh)
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

	q2 := func(d *gorm.DB) *gorm.DB {
		return d.Save(&user)
	}

	result = q2(db)
	if result.Error != nil {
		internalServerError(w)
		return
	}

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
