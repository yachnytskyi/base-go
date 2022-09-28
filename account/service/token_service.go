package service

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/google/uuid"
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
func (s *tokenService) NewPairFromUser(ctx context.Context, user *model.User, previousTokenID string) (*model.TokenPair, error) {
	if previousTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, user.UserID.String(), previousTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for userID: %v, tokenID: %v\n", user.UserID.String(), previousTokenID)

			return nil, err
		}
	}

	// No need to use a repository for idToken as it is unrelated to any data source.
	idToken, err := generateIDToken(user, s.PrivateKey, s.IDExpirationSecrets)

	if err != nil {
		log.Printf("Error generating idToken for userID: %v. Error: %v\n", user.UserID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefreshToken(user.UserID, s.RefreshSecret, s.RefreshExpirationSecrets)

	if err != nil {
		log.Printf("Error generating refreshToken for userID: %v. Error: %v\n", user.UserID, err.Error())
		return nil, apperrors.NewInternal()
	}

	// Set freshly minted refresh token to valid list.
	if err := s.TokenRepository.SetRefreshToken(ctx, user.UserID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for userID: %v. Error: %v\n", user.UserID, err.Error())
		return nil, apperrors.NewInternal()
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SignedString: idToken},
		RefreshToken: model.RefreshToken{SignedString: refreshToken.SignedString, ID: refreshToken.ID, UserID: user.UserID},
	}, nil
}

// SignOut reaches out to the repository layer to delete all valid tokens for a user.
func (s *tokenService) SignOut(ctx context.Context, userID uuid.UUID) error {
	return s.TokenRepository.DeleteUserRefreshTokens(ctx, userID.String())
}

// ValidateIDToken validates the id token jwt string.
// It returns the user extract from the IDTokenCustomClaims.
func (s *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	claims, err := validateIDToken(tokenString, s.PublicKey) // Uses public RSA key.

	// We will just return unauthorized error in all instances of failing to verify the user.
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperrors.NewAuthorization("Unable to verify the user from the idToken")
	}

	return claims.User, nil
}

// ValidateRefreshToken checks to make sure the JWT provided by a string is valid
// and returns a RefreshToken if valid.
func (s *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
	// Validate actual JWT with string a secret.
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

	// We will just return unauthorized error in all instances of failing to verify the user.
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, apperrors.NewAuthorization("Unable to verify the user from the refresh token")
	}

	// Standard claims store ID as a string. I want "model" to be our string
	// is a UUID. So we parse claims.Id as UUID.
	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, apperrors.NewAuthorization("Unable to verify user from refresh token")
	}

	return &model.RefreshToken{
		SignedString: tokenString,
		ID:           tokenUUID,
		UserID:       claims.UserID,
	}, nil
}
