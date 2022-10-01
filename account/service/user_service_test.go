package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
	"github.com/yachnytskyi/base-go/account/model/fixture"
	"github.com/yachnytskyi/base-go/account/model/mocks"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUserResponse := &model.User{
			UserID:   userID,
			Email:    "kostya@kostya.com",
			Username: "Kostya Kostyan",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		user := NewUserService(&UserConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", mock.Anything, userID).Return(mockUserResponse, nil)

		ctx := context.TODO()
		getUser, err := user.Get(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, getUser, mockUserResponse)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		user := NewUserService(&UserConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, userID).Return(nil, fmt.Errorf("Some erro down the call chain"))

		ctx := context.TODO()
		getUser, err := user.Get(ctx, userID)

		assert.Nil(t, getUser)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    "kostya@kostya.com",
			Password: "heyman!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		user := NewUserService(&UserConfig{
			UserRepository: mockUserRepository,
		})

		// We can use Run method to modify the user when the Create method is called.
		// We can then chain on a Return method to return no error.
		mockUserRepository.On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User) // arg 0 is context, arg 1 is *User.
				userArg.UserID = userID
			}).Return(nil)

		ctx := context.TODO()
		err := user.SignUp(ctx, mockUser)

		assert.NoError(t, err)

		// assert the user now has a userID.
		assert.Equal(t, userID, mockUser.UserID)

		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "kostya@kostya.com",
			Password: "heyman!",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		user := NewUserService(&UserConfig{
			UserRepository: mockUserRepository,
		})

		mockErr := apperrors.NewConflict("email", mockUser.Email)

		// We can use Run method to modify the user when the Create method is called.
		// We can then chain on a Return method to return no error.
		mockUserRepository.On("Create", mock.AnythingOfType("*context.emptyCtx"), mockUser).Return(mockErr)

		ctx := context.TODO()
		err := user.SignUp(ctx, mockUser)

		// Assert error is error we response with in mock.
		assert.EqualError(t, err, mockErr.Error())

		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignIn(t *testing.T) {
	// Setup valid email/password combo with hashed password to test method
	// response when provided password is invalid.
	email := "kostya@kostya.com"
	validPassword := "somerandomvalidpasssword"
	hashedValidPassword, _ := hashPassword(validPassword)
	invalidPassword := "somerandominvalidpassword"

	mockUserRepository := new(mocks.MockUserRepository)
	user := NewUserService(&UserConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    email,
			Password: validPassword,
		}

		mockUserResponse := &model.User{
			UserID:   userID,
			Email:    email,
			Password: hashedValidPassword,
		}

		mockArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		// We can then chain on a Return method to return no error.
		mockUserRepository.On("FindByEmail", mockArguments...).Return(mockUserResponse, nil)

		ctx := context.TODO()
		err := user.SignIn(ctx, mockUser)

		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArguments...)
	})

	t.Run("Invalid email/password combination", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    email,
			Password: invalidPassword,
		}

		mockUserResponse := &model.User{
			UserID:   userID,
			Email:    email,
			Password: hashedValidPassword,
		}

		mockArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}

		// We can use Run method to modify the user when the Create method is called.
		// We can then chain on a Return method to return no error.
		mockUserRepository.On("FindByEmail", mockArguments...).Return(mockUserResponse, nil)

		ctx := context.TODO()
		err := user.SignIn(ctx, mockUser)

		assert.Error(t, err)
		assert.EqualError(t, err, "Invalid email and password combination")
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArguments...)
	})
}

func TestUpdateDetails(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	user := NewUserService(&UserConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUser := &model.User{
			UserID:   userID,
			Email:    "new@kostya.com",
			Website:  "https://constantine.com",
			Username: "A New Kostya!",
		}

		mockArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockUserRepository.On("Update", mockArguments...).Return(nil)

		ctx := context.TODO()
		err := user.UpdateDetails(ctx, mockUser)

		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "Update", mockArguments...)
	})

	t.Run("Failure", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUser := &model.User{
			UserID: userID,
		}

		mockArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockError := apperrors.NewInternal()

		mockUserRepository.On("Update", mockArguments...).Return(mockError)

		ctx := context.TODO()
		err := user.UpdateDetails(ctx, mockUser)
		assert.Error(t, err)

		appError, ok := err.(*apperrors.Error)
		assert.True(t, ok)
		assert.Equal(t, apperrors.Internal, appError.Type)

		mockUserRepository.AssertCalled(t, "Update", mockArguments...)
	})
}

func TestSetProfileImage(t *testing.T) {
	mockUserRepository := new(mocks.MockUserRepository)
	mockImageRepository := new(mocks.MockImageRepository)

	user := NewUserService(&UserConfig{
		UserRepository:  mockUserRepository,
		ImageRepository: mockImageRepository,
	})

	t.Run("Successful creation of a new image", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		// Does not have an imageURL.
		mockUser := &model.User{
			UserID:   userID,
			Email:    "new@kostya.com",
			Website:  "https://constantine.me",
			Username: " A new Kostya!",
		}

		findByIDArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"), userID,
		}
		mockUserRepository.On("FindByID", findByIDArguments...).Return(mockUser, nil)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		imageFile, _ := imageFileHeader.Open()

		updateProfileArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			imageFile,
		}

		imageURL := "http://imageurl.com/hxkjhfwy25185auctq"

		mockImageRepository.On("UpdateProfile", updateProfileArguments...).Return(imageURL, nil)

		updateImageArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.UserID,
			imageURL,
		}

		mockUpdatedUser := &model.User{
			UserID:   userID,
			Email:    "new@kostya.com",
			Website:  "https://constantine.com",
			Username: "A new Kostya!",
			ImageURL: imageURL,
		}

		mockUserRepository.On("UpdateImage", updateImageArguments...).Return(mockUpdatedUser, nil)

		ctx := context.TODO()

		updatedUser, err := user.SetProfileImage(ctx, mockUser.UserID, imageFileHeader)

		assert.NoError(t, err)
		assert.Equal(t, mockUpdatedUser, updatedUser)
		mockUserRepository.AssertCalled(t, "FindByID", findByIDArguments...)
		mockImageRepository.AssertCalled(t, "UpdateProfile", updateProfileArguments...)
		mockUserRepository.AssertCalled(t, "UpdateImage", updateImageArguments...)
	})

	t.Run("Successful update an image", func(t *testing.T) {
		userID, _ := uuid.NewRandom()
		imageURL := "http://imageurl.com/hxkjhfwy25185auctq"

		// Has an imageURL.
		mockUser := &model.User{
			UserID:   userID,
			Email:    "new@kostya@com",
			Website:  "https://constantine.com",
			Username: "A new Kostya!",
			ImageURL: imageURL,
		}

		findByIDArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"), userID,
		}
		mockUserRepository.On("FindByID", findByIDArguments...).Return(mockUser, nil)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		imageFile, _ := imageFileHeader.Open()

		updateProfileArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			imageFile,
		}

		mockImageRepository.On("UpdateProfile", updateProfileArguments).Return(imageURL, nil)

		updateImageArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.UserID,
			imageURL,
		}

		mockUpdatedUser := &model.User{
			UserID:   userID,
			Email:    "new@kostya.com",
			Website:  "https://constantine.com",
			Username: "A new Kostya!",
			ImageURL: imageURL,
		}

		mockUserRepository.On("UpdateImage", updateImageArguments...).Return(mockUpdatedUser, nil)

		ctx := context.TODO()

		updatedUser, err := user.SetProfileImage(ctx, userID, imageFileHeader)

		assert.NoError(t, err)
		assert.Equal(t, mockUpdatedUser, updatedUser)
		mockUserRepository.AssertCalled(t, "FindByID", findByIDArguments...)
		mockImageRepository.AssertCalled(t, "UpdateProfile", updateProfileArguments...)
		mockUserRepository.AssertCalled(t, "UpdateImage", updateImageArguments...)
	})

	t.Run("UserRepository FindByID Error", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		findByIDArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userID,
		}
		mockError := apperrors.NewInternal()
		mockUserRepository.On("FindByID", findByIDArguments...).Return(nil, mockError)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()

		ctx := context.TODO()

		updatedUser, err := user.SetProfileImage(ctx, userID, imageFileHeader)

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		mockUserRepository.AssertCalled(t, "FindByID", findByIDArguments...)
		mockImageRepository.AssertNotCalled(t, "UpdateProfile")
		mockUserRepository.AssertNotCalled(t, "UpdateImage")
	})

	t.Run("ImageRepository Error", func(t *testing.T) {
		// Need to create a new UserService and repository
		// because testify has no way to overwrtie a mock's
		// "On" call.
		mockUserRepository := new(mocks.MockUserRepository)
		mockImageRepository := new(mocks.MockImageRepository)

		user := NewUserService(&UserConfig{
			UserRepository:  mockUserRepository,
			ImageRepository: mockImageRepository,
		})

		userID, _ := uuid.NewRandom()
		imageURL := "http://imageurl.com/hxkjhfwy25185auctq"

		// Has an imageURL.
		mockUser := &model.User{
			UserID:   userID,
			Email:    "new@kostya.com",
			Website:  "https://constantine.com",
			Username: "A new Kostya!",
			ImageURL: imageURL,
		}

		findByIDArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userID,
		}
		mockUserRepository.On("FindByID", findByIDArguments...).Return(mockUser, nil)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		imageFile, _ := imageFileHeader.Open()

		updateProfileArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			imageFile,
		}

		mockError := apperrors.NewInternal()
		mockImageRepository.On("UpdateProfile", updateProfileArguments...).Return(nil, mockError)

		ctx := context.TODO()
		updatedUser, err := user.SetProfileImage(ctx, userID, imageFileHeader)

		assert.Nil(t, updatedUser)
		assert.Error(t, err)
		mockUserRepository.AssertCalled(t, "FindByID", findByIDArguments...)
		mockImageRepository.AssertCalled(t, "UpdateProfile", updateProfileArguments...)
		mockUserRepository.AssertNotCalled(t, "UpdateImage")
	})

	t.Run("UserRepository UpdateImage Error", func(t *testing.T) {
		userID, _ := uuid.NewRandom()
		imageURL := "http://imageurl.com/hxkjhfwy25185auctq"

		// Has an imageURL.
		mockUser := &model.User{
			UserID:   userID,
			Email:    "new@kostya.com",
			Website:  "https://constantine.com",
			Username: "A new Kostya!",
			ImageURL: imageURL,
		}

		findByIDArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userID,
		}
		mockUserRepository.On("FindByID", findByIDArguments...).Return(mockUser, nil)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		imageFile, _ := imageFileHeader.Open()

		updateProfileArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"),
			imageFile,
		}

		mockImageRepository.On("UpdateProfile", updateProfileArguments...).Return(imageURL, nil)

		updateImageArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.UserID,
			imageURL,
		}

		mockError := apperrors.NewInternal()
		mockUserRepository.On("UpdateImage", updateImageArguments...).Return(nil, mockError)

		ctx := context.TODO()

		updatedUser, err := user.SetProfileImage(ctx, userID, imageFileHeader)

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		mockImageRepository.AssertCalled(t, "UpdateProfile", updateProfileArguments...)
		mockUserRepository.AssertCalled(t, "UpdateImage", updateImageArguments...)

	})
}
