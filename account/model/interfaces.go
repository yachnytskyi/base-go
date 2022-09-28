package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement.
type UserService interface {
	Get(ctx context.Context, userID uuid.UUID) (*User, error)
	SignUp(ctx context.Context, user *User) error
	SignIn(ctx context.Context, user *User) error
	UpdateDetails(ctx context.Context, user *User) error
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producting JWTs as string.
type TokenService interface {
	NewPairFromUser(ctx context.Context, user *User, refreshTokenID string) (*TokenPair, error)
	SignOut(ctx context.Context, userID uuid.UUID) error
	ValidateIDToken(tokenString string) (*User, error)
	ValidateRefreshToken(refreshTokenString string) (*RefreshToken, error)
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement.
type UserRepository interface {
	FindById(ctx context.Context, userID uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

// TokenRepository defines methids if expects a repository
// it interacts with to implement.
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, previousTokenID string) error
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
}
