package kiam

import (
	"github.com/dgrijalva/jwt-go"
)

func GetJWTSigningMethod(code string) jwt.SigningMethod {
	switch code {
	case "ES256":
		return jwt.SigningMethodES256

	case "HS256":
		return jwt.SigningMethodHS256
	}
	return jwt.SigningMethodHS256
}
