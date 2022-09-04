package service

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yachnytskyi/base-go/account/model"

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

	// Instantiate a common token service to be used by all tests.
	tokenService := NewTokenService(&TokenServiceConfig{
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

	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.TODO()
		tokenPair, err := tokenService.NewPairFromUser(ctx, u, "")
		assert.NoError(t, err)

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
}
