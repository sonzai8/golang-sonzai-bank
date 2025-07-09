package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func (maker JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(maker.secretKey))
}

func (maker JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		fmt.Printf("Error parsing token: %+v\n", err)
		return nil, errInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		fmt.Println("payload line 57 :", ok)
		return nil, errInvalidToken
	}

	fmt.Println("payload line60 :", payload)
	err = payload.Valid()
	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, errors.New("secret key too short")
	}
	return &JWTMaker{secretKey: secretKey}, nil
}
