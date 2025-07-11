package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	paseto        *paseto.V2
	sysmmetricKey []byte
}

func (maker PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}
	newMaker, err := maker.paseto.Encrypt(maker.sysmmetricKey, payload, nil)
	if err != nil {
		return "", payload, err
	}
	return newMaker, payload, nil
}

func (maker PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.sysmmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// func NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(sysmmetricKey string) (Maker, error) {
	if len(sysmmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalied key size: must be exactly %d charactoers", chacha20poly1305.KeySize)
	}
	return &PasetoMaker{
		paseto:        paseto.NewV2(),
		sysmmetricKey: []byte(sysmmetricKey),
	}, nil
}
