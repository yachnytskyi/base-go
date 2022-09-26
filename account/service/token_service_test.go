package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
	"github.com/yachnytskyi/base-go/account/model/mocks"

	"github.com/dgrijalva/jwt-go"
)

func TestNewPairFromUser(t *testing.T) {
	var idExpiration int64 = 15 * 60
	var refreshExpiration int64 = 3 * 24 * 2600
	private, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(private)
	public, _ := ioutil.ReadFile("../rsa_public_test.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(public)
	secret := "anothersomerandomtestsecret"

	mockTokenRepository := new(mocks.MockTokenRepository)

	// Instantiate a common token service to be used by all tests.
	tokenService := NewTokenService(&TokenServiceConfig{
		TokenRepository:          mockTokenRepository,
		PrivateKey:               privateKey,
		PublicKey:                publicKey,
		RefreshSecret:            secret,
		IDExpirationSecrets:      idExpiration,
		RefreshExpirationSecrets: refreshExpiration,
	})

	// Include password to make sure it is not serialized
	// since json tag is "-".
	userID, _ := uuid.NewRandom()
	user := &model.User{
		UserID:   userID,
		Email:    "kostya@kostya.com",
		Password: "somerandompassword",
	}

	// Setup mock call responses in setup before t.Run statements.
	userIDErrorCase, _ := uuid.NewRandom()
	userErrorCase := &model.User{
		UserID:   userIDErrorCase,
		Email:    "failed@failed.com",
		Password: "somerfailedpassword",
	}
	previousID := "a_previous_tokenID"

	setSuccessArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		user.UserID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		userIDErrorCase.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPreviousIDArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		user.UserID.String(),
		previousID,
	}

	// Mock call argument/responses.
	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("Error setting refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPreviousIDArguments...).Return(nil)

	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.Background()                                           // Updated from context.TODO().
		tokenPair, err := tokenService.NewPairFromUser(ctx, user, previousID) // Replaced "" with previousID from setup.
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments.
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should be called since previousID is not empty.
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPreviousIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.IDToken.SignedString)

		// Decode the Base64URL encoded string
		// simpler to use jwt library which is already imported.
		idTokenClaims := &idTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.IDToken.SignedString, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		assert.NoError(t, err)

		// Assert claims on idToken.
		expectedClaims := []interface{}{
			user.UserID,
			user.Email,
			user.Username,
			user.ImageURL,
			user.Website,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UserID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Username,
			idTokenClaims.User.ImageURL,
			idTokenClaims.User.Website,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // Password should never be encoded to json.

		expiresAt := time.Unix(idTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(time.Duration(idExpiration) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &refreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken.SignedString, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken.SignedString)

		// Assert claims on a refresh token.
		assert.NoError(t, err)
		assert.Equal(t, user.UserID, refreshTokenClaims.UserID)

		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshExpiration) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, userErrorCase, "")
		assert.Error(t, err) // Should return an error.

		// SetRefreshToken should be called with setErrorArguments.
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)
		// DeleteRefreshToken should not be since SetRefreshToken causes method to return.
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Empty string provided for previousID", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, user, "")
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments.
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is "".
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Previous token not in a repository", func(t *testing.T) {
		ctx := context.Background()
		userID, _ := uuid.NewRandom()
		user := &model.User{
			UserID: userID,
		}

		tokenIDNotInRepo := "not_in_token_repo"

		deleteArgs := mock.Arguments{
			ctx,
			user.UserID.String(),
			tokenIDNotInRepo,
		}

		mockError := apperrors.NewAuthorization("Invalid refresh token")
		mockTokenRepository.On("DeleteRefreshToken", deleteArgs...).Return(mockError)

		_, err := tokenService.NewPairFromUser(ctx, user, tokenIDNotInRepo)
		assert.Error(t, err)

		appError, ok := err.(*apperrors.Error)

		assert.True(t, ok)
		assert.Equal(t, apperrors.Authorization, appError.Type)
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteArgs...) // Should be called with invalid arguments.
		mockTokenRepository.AssertNotCalled(t, "SetRefreshToken")                // Should not be called, because we don't delete the previous refresh token.
	})
}

func TestValidateIDToken(t *testing.T) {
	var idExpiration int64 = 15 * 60

	private, _ := ioutil.ReadFile("../rsa_private_test.pem")
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(private)
	public, _ := ioutil.ReadFile("../rsa_public_test.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(public)

	// Instantiate a common token service to be used by all tests.
	tokenService := NewTokenService(&TokenServiceConfig{
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
		IDExpirationSecrets: idExpiration,
	})

	// Include password to make sure it is not serialized
	// since json tag is "-".
	userID, _ := uuid.NewRandom()
	user := &model.User{
		UserID:   userID,
		Email:    "kostya@kostya.com",
		Password: "somerandompassword",
	}

	t.Run("Valid token", func(t *testing.T) {
		// Maybe not the best approach to depend on utility method.
		// Token will be valid for 15 minutes.
		signedString, _ := generateIDToken(user, privateKey, idExpiration)

		userFromToken, err := tokenService.ValidateIDToken(signedString)
		assert.NoError(t, err)

		assert.ElementsMatch(
			t,
			[]interface{}{user.Email, user.Username, user.UserID, user.Website, user.ImageURL},
			[]interface{}{userFromToken.Email, userFromToken.Username, userFromToken.UserID, userFromToken.Website, userFromToken.ImageURL},
		)
	})

	t.Run("Expired token", func(t *testing.T) {
		// Maybe not the best approach to depend on utility method.
		// Token will be valid for 15 minutes.
		signedString, _ := generateIDToken(user, privateKey, -1) // Expired one second ago.

		expectedError := apperrors.NewAuthorization("Unable to verify the user from the idToken")

		_, err := tokenService.ValidateIDToken(signedString)
		assert.EqualError(t, err, expectedError.Message)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		// Maybe not the best approach to depend on utility method.
		// Token won't be valid.
		signedString, _ := generateIDToken(user, privateKey, -1) // Expired one second ago.

		expectedError := apperrors.NewAuthorization("Unable to verify the user from the idToken")

		_, err := tokenService.ValidateIDToken(signedString)
		assert.EqualError(t, err, expectedError.Message)
	})

	// TODO - Add other invalid token types (maybe in the future).
}

func TestValidateRefreshToken(t *testing.T) {
	var refreshExpiration int64 = 3 * 24 * 2600
	secret := "anothersomerandomtestsecret"

	tokenService := NewTokenService(&TokenServiceConfig{
		RefreshSecret:            secret,
		RefreshExpirationSecrets: refreshExpiration,
	})

	userID, _ := uuid.NewRandom()
	user := &model.User{
		UserID:   userID,
		Email:    "kostya@kostya.com",
		Password: "somerandomsecret",
	}

	t.Run("Valid token", func(t *testing.T) {
		testRefreshToken, _ := generateRefreshToken(user.UserID, secret, refreshExpiration)

		validatedRefreshToken, err := tokenService.ValidateRefreshToken(testRefreshToken.SignedString)
		assert.NoError(t, err)

		assert.Equal(t, user.UserID, validatedRefreshToken.UserID)
		assert.Equal(t, testRefreshToken.SignedString, validatedRefreshToken.SignedString)
	})
	t.Run("Expired token", func(t *testing.T) {
		testRefreshToken, _ := generateRefreshToken(user.UserID, secret, -1)

		expectedError := apperrors.NewAuthorization("Unable to verify the user from the refresh token")

		_, err := tokenService.ValidateRefreshToken(testRefreshToken.SignedString)
		assert.EqualError(t, err, expectedError.Message)
	})
}
