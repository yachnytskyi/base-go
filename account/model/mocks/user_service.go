package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
)

// MockUserService is a mock type for model.UserService.
type MockUserService struct {
	mock.Mock
}

// Get is mock of UserService Get.
func (m *MockUserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	// Args that will be passed to "Return" in the tests, when function
	// is called with a uid. Hence the name "ret".
	ret := m.Called(ctx, uid)

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

// SignUp is a mock of UserService.SignUp.
func (m *MockUserService) SignUp(ctx context.Context, u *model.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
