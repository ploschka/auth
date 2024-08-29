package auth

import (
	"crypto/aes"
	"crypto/cipher"
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

type RefreshToken struct {
	Ip        string `json:"ip"`
	IssuedAt  int64  `json:"iat"`
	GUID      string `json:"guid"`
	Signature string `json:"sign"`
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

// TODO
func EncryptToken(token []byte) ([]byte, error) {
	return nil, nil
}

// TODO
func DecryptToken(token string) (*RefreshToken, error) {
	return nil, nil
}

// TODO
func HashToken(token []byte) ([]byte, error) {
	return nil, nil
}

// TODO
func CheckPair(access string, refresh *RefreshToken) bool {
	return false
}

// TODO
func Validate(tok *RefreshToken, hash []byte) bool {
	return false
}

func GenerateTokens(ip string, user model.User) (access string, refresh *RefreshToken, err error) {
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

	refTok := &RefreshToken{
		Ip:        ip,
		IssuedAt:  currTime.Unix(),
		GUID:      user.Guid,
		Signature: encodedSignature,
	}

	return strtok + "." + encodedSignature, refTok, nil
}
