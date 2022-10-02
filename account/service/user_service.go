package service

import (
	"context"
	"log"
	"mime/multipart"
	"net/url"
	"path"

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

// ClearProfileImage finds a user by its ID,
// than removes an image from Google Cloud,
// and finally removes an imageURL
// string in postgres.
func (s *userService) ClearProfileImage(ctx context.Context, userID uuid.UUID) error {
	user, err := s.UserRepository.FindByID(ctx, userID)

	if err != nil {
		return err
	}

	if user.ImageURL == "" {
		return nil
	}

	objectName, err := objectNameFromUrl(user.ImageURL)
	if err != nil {
		return err
	}

	err = s.ImageRepository.DeleteProfile(ctx, objectName)
	if err != nil {
		return err
	}

	_, err = s.UserRepository.UpdateImage(ctx, userID, "")

	if err != nil {
		return err
	}

	return nil
}

// Get retrieves a user based on their uuid.
func (s *userService) Get(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := s.UserRepository.FindByID(ctx, userID)

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

func (s *userService) SetProfileImage(ctx context.Context, userID uuid.UUID, imageFileHeader *multipart.FileHeader) (*model.User, error) {
	user, err := s.UserRepository.FindByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	objectName, err := objectNameFromUrl(user.ImageURL)

	if err != nil {
		return nil, err
	}

	imageFile, err := imageFileHeader.Open()

	if err != nil {
		log.Printf("Failed to ope the image file: %v\n", err)
		return nil, apperrors.NewInternal()
	}

	// Upload a user's image to ImageRepository.
	// Possibly received an updated imageURL.
	imageURL, err := s.ImageRepository.UpdateProfile(ctx, objectName, imageFile)

	if err != nil {
		log.Printf("Unable to upload the image to the cloud provider: %v\n", err)
		return nil, err
	}

	updatedUser, err := s.UserRepository.UpdateImage(ctx, user.UserID, imageURL)

	if err != nil {
		log.Printf("Unable to update the imageURL: %v\n", err)
		return nil, err
	}

	return updatedUser, nil
}

func objectNameFromUrl(imageURL string) (string, error) {
	// If a user does not have an imageURL - create one.
	// Otherwise, extract the last part of the URL to get a cloud storage object name.
	if imageURL == "" {
		objectID, _ := uuid.NewRandom()
		return objectID.String(), nil
	}

	// Split off the last part of the URL, which is the image's storage object ID.
	urlPath, err := url.Parse(imageURL)

	if err != nil {
		log.Printf("Failed to parse objectName from the imageURL: %v\n", imageURL)
		return "", apperrors.NewInternal()
	}
	// Get "path" of an url (everything is after a domain).
	// Then get "base", the last part.
	return path.Base(urlPath.Path), nil
}
