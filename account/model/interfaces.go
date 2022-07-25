package model

import (
	"context"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
	SignUp(ctx context.Context, u *User) error
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
}
