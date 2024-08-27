package auth

import (
	_ "crypto/sha512"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ploschka/auth/internal/model"
	_ "golang.org/x/crypto/bcrypt"
)

const (
	hashCost             int           = 15
	accessTokenDuration  time.Duration = 1 * time.Hour
	refreshTokenDuration time.Duration = 24 * time.Hour
)

var (
	allowedMethods = [...]string{
		jwt.SigningMethodHS512.Alg(),
	}
)

var (
	signKey       []byte
	encryptionKey []byte
	hashKey       []byte
)

type claims struct {
	jwt.RegisteredClaims
	Ip    string `json:"ip"`
	Admin bool   `json:"admin"`
	GUID  string `json:"guid"`
}

type refreshToken struct {
	Ip       string `json:"ip"`
	IssuedAt int64  `json:"iat"`
	GUID     string `json:"guid"`
}

// TODO
func (tok refreshToken) encrypt() ([]byte, error) {
	return nil, nil
}

// TODO
func decrypt(encrypted []byte) (tok refreshToken, err error) {
	return
}

func init() {
	base64SignKey, ok := os.LookupEnv("SIGN_KEY")
	if !ok || len(base64SignKey) == 0 {
		panic("SIGN_KEY is undefined or emty")
	}

	base64EncryptionKey, ok := os.LookupEnv("ENCRYPTION_KEY")
	if !ok || len(base64EncryptionKey) == 0 {
		panic("ENCRYPTION_KEY is undefined or emty")
	}

	base64HashKey, ok := os.LookupEnv("HASH_KEY")
	if !ok || len(base64EncryptionKey) == 0 {
		panic("HASH_KEY is undefined or emty")
	}

	var err error
	var tempErr error

	signKey, tempErr = base64.StdEncoding.DecodeString(base64SignKey)
	err = errors.Join(err, tempErr)

	encryptionKey, tempErr = base64.StdEncoding.DecodeString(base64EncryptionKey)
	err = errors.Join(err, tempErr)

	hashKey, tempErr = base64.StdEncoding.DecodeString(base64HashKey)
	err = errors.Join(err, tempErr)

	if err != nil {
		panic(err)
	}
}

func generateTokens(ip string, user model.User) (access string, refresh string, err error) {
	var tempErr error = nil

	currTime := time.Now()
	accessExp := currTime.Add(accessTokenDuration)

	claims := claims{
		GUID:  user.Guid,
		Admin: user.Admin,
		Ip:    ip,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currTime),
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}

	access, tempErr = jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(signKey)
	err = errors.Join(err, tempErr)

	refTok := refreshToken{
		Ip:       ip,
		IssuedAt: currTime.Unix(),
		GUID:     user.Guid,
	}

	return
}

// TODO
func RefreshOrGenerate(ip string, user model.User) (string, string, error) {
	return "", "", nil
}
