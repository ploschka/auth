package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	_ "crypto/sha512"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ploschka/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
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

func EncryptToken(token []byte) ([]byte, error) {
	nonce := make([]byte, gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, token, nil), nil
}

func DecryptToken(token []byte) ([]byte, error) {
	nonceSize := gcm.NonceSize()
	if len(token) < nonceSize {
		return nil, errors.New("encrypted text invalid length")
	}

	nonce, ciphertext := token[:nonceSize], token[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func HashToken(token []byte) ([]byte, error) {
	hasher := sha512.New()
	_, err := hasher.Write(token)
	if err != nil {
		return nil, err
	}
	sum := hasher.Sum(nil)
	return bcrypt.GenerateFromPassword(sum, hashCost)
}

func CheckPair(access string, refresh *RefreshToken) bool {
	keyfunc := func(_ *jwt.Token) (interface{}, error) {
		return signKey, nil
	}

	token, err := jwt.ParseWithClaims(access, &claims{}, keyfunc, jwt.WithValidMethods(allowedMethods[:]))
	if err != nil {
		return false
	}

	myclaims, ok := token.Claims.(claims)
	if !ok {
		return false
	}

	if myclaims.Ip != refresh.Ip {
		return false
	}

	if base64.RawURLEncoding.EncodeToString(token.Signature) != refresh.Signature {
		return false
	}

	if myclaims.IssuedAt.Unix() != refresh.IssuedAt {
		return false
	}

	return true
}

func Validate(tok []byte, hash []byte) (bool, error) {
	hasher := sha512.New()
	_, err := hasher.Write(tok)
	if err != nil {
		return false, err
	}
	toksum := hasher.Sum(nil)

	err = bcrypt.CompareHashAndPassword(hash, toksum)
	return err == nil, nil
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
