package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// tokenService used for injecting an implementation of TokenRepository for use in service methods
// along with keys and secrets forsigning JWTs.
type tokenService struct {
	// TokenRepository model.TokenRepository.
	PrivateKey               *rsa.PrivateKey
	PublicKey                *rsa.PublicKey
	RefreshSecret            string
	IDExpirationSecrets      int64
	RefreshExpirationSecrets int64
}

// TokenServiceConfig will hold repositories that will eventually be injected into this service layer.
type TokenServiceConfig struct {
	// TokenRepository model.TokenRepository
	PrivateKey               *rsa.PrivateKey
	PublicKey                *rsa.PublicKey
	RefreshSecret            string
	IDExpirationSecrets      int64
	RefreshExpirationSecrets int64
}

// NewTokenService is a factory function for initializing a UserService with its repository layer dependencies.
func NewTokenService(c *TokenServiceConfig) model.TokenService {
	return &tokenService{
		PrivateKey:               c.PrivateKey,
		PublicKey:                c.PublicKey,
		RefreshSecret:            c.RefreshSecret,
		IDExpirationSecrets:      c.IDExpirationSecrets,
		RefreshExpirationSecrets: c.RefreshExpirationSecrets,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user.
// If a previous token is included, the previous token is removed from the tokens repository.
func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, previousTokenID string) (*model.TokenPair, error) {
	// No need to use a repository for idToken as it is unrelated to any data source.
	idToken, err := generateIDToken(u, s.PrivateKey, s.IDExpirationSecrets)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecrets)

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
