package auth

import (
	_ "crypto/sha512"

	"github.com/golang-jwt/jwt/v5"
	_ "golang.org/x/crypto/bcrypt"
)

var (
	allowedMethods = [...]string{
		jwt.SigningMethodHS512.Alg(),
	}
)
