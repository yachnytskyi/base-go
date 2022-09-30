package mocks

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
)

// MockUserService is a mock type for model.UserService.
type MockUserService struct {
	mock.Mock
}

// Get is mock of UserService Get.
func (m *MockUserService) Get(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	// Args that will be passed to "Return" in the tests, when function
	// is called with a userID. Hence the name "ret".
	ret := m.Called(ctx, userID)

	// First value passed to "Return".
	var r0 *model.User
	if ret.Get(0) != nil {
		// We can just return this if we know we won't be passing function to "Return".
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// SignUp is a mock of UserService.SignUp
func (m *MockUserService) SignUp(ctx context.Context, user *model.User) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// SignIn is a mock for UserService.SignIn
func (m *MockUserService) SignIn(ctx context.Context, user *model.User) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// UpdateDetails is a mock of UserService.UpdateDetails
func (m *MockUserService) UpdateDetails(ctx context.Context, user *model.User) error {
	ret := m.Called(ctx, user)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// SetProfileImage is a mock of UserService.SetProfileImage
func (m *MockUserService) SetProfileImage(ctx context.Context, userID uuid.UUID, imageFileHeader *multipart.FileHeader) (*model.User, error) {
	ret := m.Called(ctx, userID, imageFileHeader)

	// First value passed to "Return"
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
