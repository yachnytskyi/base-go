package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// userService acts as a struct for injecting
// an implementation of UserRepository
// for use in service methods.
type userService struct {
	UserRepository  model.UserRepository
	ImageRepository model.ImageRepository
}

// UserConfig will hold repositories that
// will eventually be injected into
// this service layer.
type UserConfig struct {
	UserRepository  model.UserRepository
	ImageRepository model.ImageRepository
}

// NewUserService is a factory function for
// initializing a UserService with its
// repository layer dependencies.
func NewUserService(c *UserConfig) model.UserService {
	return &userService{
		UserRepository:  c.UserRepository,
		ImageRepository: c.ImageRepository,
	}
}

// Get retrieves a user based on their uuid.
func (s *userService) Get(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.UserRepository.FindById(ctx, userID)

	return user, err
}

// SignUp reaches out to a UserRepository to verify the
// email adress is available and signs up the user
// if this is the case.
func (s *userService) SignUp(ctx context.Context, user *model.User) error {
	password, err := hashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return apperrors.NewInternal()
	}

	// Assign the hashPassword to the User.
	user.Password = password

	if err := s.UserRepository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

// SignIn reaches our to a UserRepository check if the user exists
// and when compares the supplied password with the provided password
// if a valid email/password combo is provided, u will hold all
// available user fields.
func (s *userService) SignIn(ctx context.Context, user *model.User) error {
	userFetched, err := s.UserRepository.FindByEmail(ctx, user.Email)

	// Will return NotAuthorized to client to omit details of why.
	if err != nil {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	// verify password - we previously created this method.
	match, err := comparePasswords(userFetched.Password, user.Password)

	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	*user = *userFetched
	return nil
}

func (s *userService) UpdateDetails(ctx context.Context, user *model.User) error {
	// Update a user in UserRepository.
	err := s.UserRepository.Update(ctx, user)

	if err != nil {
		return err
	}

	return nil
}
