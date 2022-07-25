package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/mocks"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResponse := &model.User{
			UID:      uid,
			Email:    "kostya.com",
			Username: "Kostya Kostyan",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&UserConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindById", mock.Anything, uid).Return(mockUserResponse, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResponse)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&UserConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindById", mock.Anything, uid).Return(nil, fmt.Errorf("Some erro down the call chain"))

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}
