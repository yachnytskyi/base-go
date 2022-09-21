package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// tokenService used for injecting an implementation
// of TokenRepository for use in service methods
// along with keys and secrets forsigning JWTs.
type tokenService struct {
	TokenRepository          model.TokenRepository
	PrivateKey               *rsa.PrivateKey
	PublicKey                *rsa.PublicKey
	RefreshSecret            string
	IDExpirationSecrets      int64
	RefreshExpirationSecrets int64
}

// TokenServiceConfig will hold repositories
// that will eventually be injected
// into this service layer.
type TokenServiceConfig struct {
	TokenRepository          model.TokenRepository
	PrivateKey               *rsa.PrivateKey
	PublicKey                *rsa.PublicKey
	RefreshSecret            string
	IDExpirationSecrets      int64
	RefreshExpirationSecrets int64
}

// NewTokenService is a factory function
// for initializing a UserService
// with its repository layer dependencies.
func NewTokenService(c *TokenServiceConfig) model.TokenService {
	return &tokenService{
		TokenRepository:          c.TokenRepository,
		PrivateKey:               c.PrivateKey,
		PublicKey:                c.PublicKey,
		RefreshSecret:            c.RefreshSecret,
		IDExpirationSecrets:      c.IDExpirationSecrets,
		RefreshExpirationSecrets: c.RefreshExpirationSecrets,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user.
// If a previous token is included, the previous token
// is removed from the tokens repository.
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

	// Set freshly minted refresh token to valid list.
	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokedID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// Delete user's current refresh token (used when refreshing idToken).
	if previousTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), previousTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.UID.String(), previousTokenID)
		}
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SignedString: idToken},
		RefreshToken: model.RefreshToken{SignedString: refreshToken.SignedString, ID: refreshToken.ID, UID: u.UID},
	}, nil
}

// ValidateIDToken validates the id token jwt string.
// It returns the user extract from the IDTokenCustomClaims.
func (s *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := validateIDToken(tokenString, s.PublicKey) // Uses public RSA key.

	// We will just return unauthorized error in all instances of failing to verify the user.
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify the user from idToken")
	}

	return claims.User, nil
}
