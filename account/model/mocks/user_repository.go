package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
)

// MockUserRepository is a mock type for model.UserRepository
type MockUserRepository struct {
	mock.Mock
}

// FindByID is a mock of UserRepository FindByID.
func (m *MockUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	ret := m.Called(ctx, userID)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// Create is a mock for UserRepository Create.
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// FindByEmail is a mock of UserRepository.FindByEmail
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ret := m.Called(ctx, email)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// Update is a mock of UserRepository.Update
func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// UpdateImage is a mock of UserRepository.UpdateImage
func (m *MockUserRepository) UpdateImage(ctx context.Context, userID uuid.UUID, imageURL string) (*model.User, error) {
	ret := m.Called(ctx, userID, imageURL)

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
