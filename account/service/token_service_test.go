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
	secret := "anotsorandomtestsecret"

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
	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "kostya@kostya.com",
		Password: "somerandompassword",
	}

	// Setup mock call responses in setup before t.Run statements.
	uidErrorCase, _ := uuid.NewRandom()
	uErrorCase := &model.User{
		UID:      uidErrorCase,
		Email:    "failed@failed.com",
		Password: "somefailedpassword",
	}
	previousID := "a_previous_tokenID"

	setSuccessArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		u.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		uidErrorCase.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPreviousIDArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		u.UID.String(),
		previousID,
	}

	// Mock call argument/responses.
	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("Error setting refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPreviousIDArguments...).Return(nil)

	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.Background()                                        // Updated from context.TODO()
		tokenPair, err := tokenService.NewPairFromUser(ctx, u, previousID) // Replaced "" with previousID from setup.
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments.
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since previousID is "".
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPreviousIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.IDToken)

		// Decode the Base64URL encoded string
		// simpler to use jwt library which is already imported.
		idTokenClaims := &IDTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.IDToken, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		assert.NoError(t, err)

		// Assert claims on idToken.
		expectedClaims := []interface{}{
			u.UID,
			u.Email,
			u.Username,
			u.ImageURL,
			u.Website,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Username,
			idTokenClaims.User.ImageURL,
			idTokenClaims.User.Website,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // Password must never be encoded to json.

		expiresAt := time.Unix(idTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(time.Duration(idExpiration) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &RefreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken)

		// assert claims on refresh token.
		assert.NoError(t, err)
		assert.Equal(t, u.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshExpiration) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, uErrorCase, "")
		assert.Error(t, err) // Should return an error.

		// SetRefreshToken should be called with setErrorArguments.
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)
		// DeleteRefreshToken should not be since SetRefreshToken causes method to return.
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Empty string provided for previousID", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, u, "")
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments.
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since previousID is "".
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
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

	// Include a password to make sure it is not serialized
	// since json tag is "-".
	uid, _ := uuid.NewRandom()
	user := &model.User{
		UID:      uid,
		Email:    "kostya@kostya.com",
		Password: "somerandompassword",
	}

	t.Run("Valid token", func(t *testing.T) {
		// Maybe not the best approach to defend on utility method.
		// Token will be valid for 15 minutes.
		signedString, _ := generateIDToken(user, privateKey, idExpiration)

		userFromToken, err := tokenService.ValidateIDToken(signedString)
		assert.NoError(t, err)

		assert.ElementsMatch(
			t,
			[]interface{}{user.Email, user.Username, user.UID, user.Website, user.ImageURL},
			[]interface{}{userFromToken.Email, userFromToken.Username, userFromToken.UID, userFromToken.Website, userFromToken.ImageURL},
		)
	})

	t.Run("Expired token", func(t *testing.T) {
		// Maybe not the best approach to defend on utility method.
		// Token will be valid for 15 minutes.
		signedString, _ := generateIDToken(user, privateKey, -1) // Expired one second ago.

		expectedError := apperrors.NewAuthorization("Unable to verify the user from idToken")

		_, err := tokenService.ValidateIDToken(signedString)
		assert.EqualError(t, err, expectedError.Message)
	})

	t.Run("Invalid signature", func(t *testing.T) {
		// Maybe not the best approach to defend on utility method.
		// Token will be valid for 15 minutes.
		signedString, _ := generateIDToken(user, privateKey, -1) // Expired one second ago.

		expectedError := apperrors.NewAuthorization("Unable to verify the user from idToken")

		_, err := tokenService.ValidateIDToken(signedString)
		assert.EqualError(t, err, expectedError.Message)
	})

}
