package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type TokenType byte

var (
	ErrExpiredToken = errors.New("token is expired")
	errInvalidToken = errors.New("invalid token")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssueAt   time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (payload Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: payload.ExpiredAt}, nil
}

func (payload Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Now()}, nil
}

func (payload Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Now()}, nil
}

func (payload Payload) GetIssuer() (string, error) {
	return payload.Username, nil
}

func (payload Payload) GetSubject() (string, error) {
	return payload.Username, nil
}

func (payload Payload) GetAudience() (jwt.ClaimStrings, error) {
	//TODO implement me
	panic("implement me")
}

// Valid checks if the token Payload is valid or not
func (payload *Payload) Valid() error {

	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:        tokenID,
		Username:  username,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}
