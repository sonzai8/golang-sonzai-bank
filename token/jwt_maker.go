package token

import (
	"encoding/base64"
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
		fmt.Println("token.Method: ", token.Method)
		dume, ok := token.Method.(*jwt.SigningMethodHMAC)
		fmt.Println("dume.SigningMethodHMAC:", dume)
		if !ok {
			fmt.Printf("err in function: : %+v\n", ok)
			return nil, errInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	encoded := base64.StdEncoding.EncodeToString([]byte(token))
	fmt.Printf("token base64 in validate : %s\n", encoded)
	fmt.Printf("token base64 in validate : %+v\n", err)

	if err != nil {
		fmt.Println("payload line 46 :", err)
		fmt.Printf("payload line 49 :%+v \n", err)
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrInvalidKey) {
			return nil, jwt.ErrInvalidKey
		}
		return nil, err
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
