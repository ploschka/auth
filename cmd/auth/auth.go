package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
)

func main() {
	secret := []byte("secret_key")
	k := hmac.New(sha512.New, secret)
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg": "HS512", "typ": "JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"guid": "123-456-789", "ip": "192.168.0.1", "role": "admin"}`))
	token := header + "." + payload
	k.Write([]byte(token))
	signature := base64.RawURLEncoding.EncodeToString(k.Sum(nil))
	fmt.Println(token + "." + signature)
}
