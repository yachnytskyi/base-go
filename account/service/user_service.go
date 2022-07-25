package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/yachnytskyi/base-go/account/model"
)

// UserService acts as a struct for injecting an implementation of UserRepository
// for use in service methods.
type UserService struct {
	UserRepository model.UserRepository
}

// UserConfig will hold repositories that will eventually be injected into
// this service layer.
type UserConfig struct {
	UserRepository model.UserRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies.
func NewUserService(c *UserConfig) model.UserService {
	return &UserService{
		UserRepository: c.UserRepository,
	}
}

// Get retrieves a user based on their uuid.
func (s *UserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindById(ctx, uid)

	return u, err
}

// SignUp reaches out to a UserRepository to verify the
// email adress is available and signs up the user if this is the case.
func (s *UserService) SignUp(ctx context.Context, u *model.User) error {
	panic("Method not implemented")
}
