package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// userService acts as a struct for injecting an implementation of UserRepository
// for use in service methods.
type userService struct {
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
	return &userService{
		UserRepository: c.UserRepository,
	}
}

// Get retrieves a user based on their uuid.
func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindById(ctx, uid)

	return u, err
}

// SignUp reaches out to a UserRepository to verify the
// email adress is available and signs up the user if this is the case.
func (s *userService) SignUp(ctx context.Context, u *model.User) error {
	password, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", u.Email)
		return apperrors.NewInternal()
	}

	// Assign the hashPassword to the User.
	u.Password = password

	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	return nil
}

// SignIn reaches our to a UserRepository check if the user exists
// and when compares the supplied password with the provided password
// if a valid email/password combo is provided, u will hold all
// available user fields.
func (s *userService) SignIn(ctx context.Context, u *model.User) error {
	panic("Not implemented")
}
