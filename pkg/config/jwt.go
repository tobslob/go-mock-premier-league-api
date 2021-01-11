package config

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	Email string
	jwt.StandardClaims
}

var (
	// ErrJWTExpired is jwt error format
	ErrJWTExpired = errors.New("token has expired")
	// ErrInvalidToken is jwt error format
	ErrInvalidToken = errors.New("token is an invalid")
	// ErrNoClaims is jwt error format
	ErrNoClaims = errors.New("no claims in token")
)

// Encode generates a jwt token
func (j *JwtWrapper) Encode(email string) (signedToken string, err error) {
	claims := &JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return signedToken, nil
}

// Decode validates the jwt token
func (j *JwtWrapper) Decode(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		err = ErrInvalidToken
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = ErrNoClaims
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Unix() {
		err = ErrJWTExpired
		return
	}
	return claims, nil
}
