package auth

import (
	"crypto/aes"
	"crypto/cipher"
	_ "crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ploschka/auth/internal/model"
	_ "golang.org/x/crypto/bcrypt"
)

const (
	hashCost            int           = 15
	accessTokenDuration time.Duration = 1 * time.Hour
)

var (
	allowedMethods = [...]string{
		jwt.SigningMethodHS512.Alg(),
	}
)

var (
	signKey       []byte
	encryptionKey []byte

	gcm cipher.AEAD
)

type claims struct {
	jwt.RegisteredClaims
	Ip    string `json:"ip"`
	Admin bool   `json:"admin"`
}

type refreshToken struct {
	Ip        string `json:"ip"`
	IssuedAt  int64  `json:"iat"`
	GUID      string `json:"guid"`
	Signature string `json:"sign"`
}

// TODO
func EncryptToken(token []byte) (string, error) {
	return "", nil
}

// TODO
func decryptToken(token string) ([]byte, error) {
	return nil, nil
}

// TODO
func HashToken(token []byte) (string, error) {
	return "", nil
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

	var err error

	signKey, err = base64.StdEncoding.DecodeString(base64SignKey)
	if err != nil {
		panic(err)
	}

	encryptionKey, err = base64.StdEncoding.DecodeString(base64EncryptionKey)
	if err != nil {
		panic(err)
	}

	if len(encryptionKey) != 32 {
		panic(errors.New("encryption key length is not 32 bytes"))
	}

	var block cipher.Block
	block, err = aes.NewCipher(encryptionKey)
	if err != nil {
		panic(err)
	}

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
}

func GenerateTokens(ip string, user model.User) (access string, refresh []byte, err error) {
	currTime := time.Now()
	accessExp := currTime.Add(accessTokenDuration)

	claims := claims{
		Admin: user.Admin,
		Ip:    ip,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currTime),
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}

	var aToken *jwt.Token
	var strtok string

	aToken = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	strtok, err = aToken.SigningString()
	if err != nil {
		return
	}

	var signature []byte
	signature, err = aToken.Method.Sign(strtok, signKey)
	if err != nil {
		return
	}

	encodedSignature := aToken.EncodeSegment(signature)

	refTok := refreshToken{
		Ip:        ip,
		IssuedAt:  currTime.Unix(),
		GUID:      user.Guid,
		Signature: encodedSignature,
	}

	refJson, err := json.Marshal(refTok)
	if err != nil {
		return
	}

	return strtok + "." + encodedSignature, refJson, nil
}

// TODO
func RefreshTokens(ip string, user model.User) (string, []byte, bool, error) {
	return "", nil, false, nil
}
