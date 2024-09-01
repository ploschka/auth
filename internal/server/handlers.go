package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/ploschka/auth/internal/auth"
	log "github.com/ploschka/auth/internal/logger"
	"github.com/ploschka/auth/internal/mailer"
	"github.com/ploschka/auth/internal/model"
	"gorm.io/gorm"
)

var (
	ErrGuidLength         error = errors.New("Invalid guid length")
	ErrInvalidTokenPair   error = errors.New("Invalid token pair")
	ErrInvalidRefresh     error = errors.New("Invalid refresh token")
	ErrNoRefreshAvailable error = errors.New("Refresh is not available for this user")
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(r.RequestURI)
	log.Info(r.RemoteAddr)
	guid := r.URL.Query().Get("guid")
	if len(guid) == 0 {
		badRequest(w, ErrGuidLength)
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
		badRequest(w, result.Error)
		return
	}

	ip := r.RemoteAddr
	slicedIp := strings.Split(ip, ":")
	ip = strings.Join(slicedIp[:len(slicedIp)-1], ":")

	access, refresh, err := auth.GenerateTokens(ip, user)
	if err != nil {
		internalServerError(w, err)
		return
	}

	refJson, err := json.Marshal(refresh)
	if err != nil {
		internalServerError(w, err)
		return
	}

	encrypted, err := auth.EncryptToken(refJson)
	if err != nil {
		internalServerError(w, err)
		return
	}

	hashed, err := auth.HashToken(refJson)
	if err != nil {
		internalServerError(w, err)
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
		internalServerError(w, result.Error)
		return
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		internalServerError(w, err)
		return
	}

	_, err = w.Write(respJson)
	if err != nil {
		internalServerError(w, err)
		return
	}
	log.Info(http.StatusOK)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(r.RequestURI)
	log.Info(r.RemoteAddr)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		internalServerError(w, err)
		return
	}

	req := tokenResponse{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		badRequest(w, err)
		return
	}

	rawRefresh, err := base64.RawURLEncoding.DecodeString(req.Refresh)
	if err != nil {
		badRequest(w, err)
		return
	}

	rawRefresh, err = auth.DecryptToken(rawRefresh)
	if err != nil {
		badRequest(w, err)
		return
	}

	refresh := &auth.RefreshToken{}
	err = json.Unmarshal(rawRefresh, refresh)

	valid := auth.CheckPair(req.Access, refresh)
	if !valid {
		unauthorized(w, ErrInvalidTokenPair)
		return
	}

	ip := r.RemoteAddr
	slicedIp := strings.Split(ip, ":")
	ip = strings.Join(slicedIp[:len(slicedIp)-1], ":")

	user := model.User{
		Guid: refresh.GUID,
	}

	db := model.GetDB()

	q := func(d *gorm.DB) *gorm.DB {
		return d.Where(&user).First(&user)
	}

	result := q(db)
	if result.Error != nil {
		unauthorized(w, err)
		return
	}

	if !user.RefreshKey.Valid {
		unauthorized(w, ErrNoRefreshAvailable)
		return
	}

	rawHash, err := base64.RawURLEncoding.DecodeString(user.RefreshKey.String)
	if err != nil {
		internalServerError(w, err)
		return
	}

	valid, err = auth.Validate(rawRefresh, rawHash)
	if err != nil {
		internalServerError(w, err)
	}
	if !valid {
		unauthorized(w, ErrInvalidRefresh)
		return
	}

	if ip != refresh.Ip {
		err = mailer.SendIpWarning(user.Email, ip)
		if err != nil {
			internalServerError(w, err)
			return
		}
	}

	var access string
	access, refresh, err = auth.GenerateTokens(ip, user)
	if err != nil {
		internalServerError(w, err)
		return
	}

	rawRefresh, err = json.Marshal(refresh)
	if err != nil {
		internalServerError(w, err)
		return
	}

	encrypted, err := auth.EncryptToken(rawRefresh)
	if err != nil {
		internalServerError(w, err)
		return
	}

	hashed, err := auth.HashToken(rawRefresh)
	if err != nil {
		internalServerError(w, err)
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
		internalServerError(w, result.Error)
		return
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		internalServerError(w, err)
		return
	}

	_, err = w.Write(respJson)
	if err != nil {
		internalServerError(w, err)
		return
	}
	log.Info(http.StatusOK)
}
