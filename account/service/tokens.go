package service

import (
	"crypto/rsa"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/yachnytskyi/base-go/account/model"
)

// IDTokenCustomClaims holds structure of jwt claims of idToken.
type IDTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

// generateIDToken generates an IDToken which is a jwt with myCustomClaims.
// Could call this GenerateIDTokenString, but the signature makes this fairly clear.
func generateIDToken(u *model.User, key *rsa.PrivateKey) (string, error) {
	unixTime := time.Now().Unix()
	tokenExpiration := unixTime + 60*15 // 15 minutes from current time.

	claims := IDTokenCustomClaims{
		User: u,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExpiration,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedString, err := token.SignedString(key)

	if err != nil {
		log.Println("Failed to sign id token string")
		return "", err
	}

	return signedString, nil
}

// RefreshToken holds the actual signed jwt string along with the ID.
// We return the id so it can be used without re-parsing the JWT from signed string.
type RefreshToken struct {
	SignedString string
	ID           string
	ExpiresIn    time.Duration
}

// RefreshTokenCustomClaims holds the payload of a refresh token.
// This can be used to extract user id for subsequent application operations (IE, fetch user in Redis).
type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

// generateRefreshToken creates a refresh token.
// The refresh token stores only the user's ID, a string.
func generateRefreshToken(uid uuid.UUID, key string) (*RefreshToken, error) {
	currentTime := time.Now()
	tokenExp := currentTime.AddDate(0, 0, 3) // 3 days.
	tokenID, err := uuid.NewRandom()         // v4 uuid in the google uuid lib.

	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}

	claims := RefreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	return &RefreshToken{
		SignedString: signedString,
		ID:           tokenID.String(),
		ExpiresIn:    tokenExp.Sub(currentTime),
	}, nil
}
