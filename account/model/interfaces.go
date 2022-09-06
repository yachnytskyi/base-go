package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement.
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	SignUp(ctx context.Context, u *User) error
	SignIn(ctx context.Context, u *User) error
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producting JWTs as string.
type TokenService interface {
	NewPairFromUser(ctx context.Context, u *User, refreshTokenID string) (*TokenPair, error)
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement.
type UserRepository interface {
	FindById(ctx context.Context, uid uuid.UUID) (*User, error)
	Create(ctx context.Context, u *User) error
}

// TokenRepository defines methids if expects a repository
// it interacts with to implement.
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, previousTokenID string) error
}
