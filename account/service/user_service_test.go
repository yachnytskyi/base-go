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
		us := NewUserService(&UserConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindById", mock.Anything, userID).Return(mockUserResponse, nil)

		ctx := context.TODO()
		user, err := us.Get(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, user, mockUserResponse)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&UserConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindById", mock.Anything, userID).Return(nil, fmt.Errorf("Some erro down the call chain"))

		ctx := context.TODO()
		user, err := us.Get(ctx, userID)

		assert.Nil(t, user)
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
