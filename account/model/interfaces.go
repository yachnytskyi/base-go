package model

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement.
type UserService interface {
	ClearProfileImage(ctx context.Context, userID uuid.UUID) error
	Get(ctx context.Context, userID uuid.UUID) (*User, error)
	SignUp(ctx context.Context, user *User) error
	SignIn(ctx context.Context, user *User) error
	UpdateDetails(ctx context.Context, user *User) error
	SetProfileImage(ctx context.Context, userID uuid.UUID, imageFileHeader *multipart.FileHeader) (*User, error)
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
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateImage(ctx context.Context, userID uuid.UUID, imageURL string) (*User, error)
}

// TokenRepository defines methids if expects a repository
// it interacts with to implement.
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, previousTokenID string) error
	DeleteUserRefreshTokens(ctx context.Context, userID string) error
}

// ImageRepository defines methods it expects a repository.
// It interacts with to implement.
type ImageRepository interface {
	DeleteProfile(ctx context.Context, objectName string) error
	UpdateProfile(ctx context.Context, objectName string, imageFile multipart.File) (string, error)
}
