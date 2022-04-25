package model

import (
	"context"

	"github.com/google/uuid"
)

// Service interface defines methods the handler layer expects
// any service it interacts with to implement

// Repository interface defines methods the service layer expects
// any repository it iteracts with to implement

type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
}

type UserRepository interface {
	FindById(ctx context.Context, uid uuid.UUID) (*User, error)
}
