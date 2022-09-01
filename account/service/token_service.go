package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// TokenService used for injecting an implementation of TokenRepository for use in service methods
// along with keys and secrets forsigning JWTs.
type TokenService struct {
	// TokenRepository model.TokenRepository.
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

// TokenServiceConfig will hold repositories that will eventually be injected into this service layer.
type TokenServiceConfig struct {
	// TokenRepository model.TokenRepository
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

// NewTokenService is a factory function for initializing a UserService with its repository layer dependencies.
func NewTokenService(c *TokenServiceConfig) model.TokenService {
	return &TokenService{
		PrivateKey:    c.PrivateKey,
		PublicKey:     c.PublicKey,
		RefreshSecret: c.RefreshSecret,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user.
// If a previous token is included, the previous token is removed from the tokens repository.
func (s *TokenService) NewPairFromUser(ctx context.Context, u *model.User, previousTokenID string) (*model.TokenPair, error) {
	// No need to use a repository for idToken as it is unrelated to any data source.
	idToken, err := generateIDToken(u, s.PrivateKey)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// TODO: store refresh tokens by calling TokenRepository methods.

	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SignedString,
	}, nil
}
