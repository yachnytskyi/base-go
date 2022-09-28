package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
)

// MockTokenService is a mock type for model.TokenService
type MockTokenService struct {
	mock.Mock
}

// NewPairFromUser mocks concrete NewPairFromUser.
func (m *MockTokenService) NewPairFromUser(ctx context.Context, user *model.User, refreshTokenID string) (*model.TokenPair, error) {
	ret := m.Called(ctx, user, refreshTokenID)

	// First value passed to "Return".
	var r0 *model.TokenPair
	if ret.Get(0) != nil {
		// We can just return this if we know we won't be passing the function to "Return".
		r0 = ret.Get(0).(*model.TokenPair)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// SignOut mocks concrete SignOut.
func (m *MockTokenService) SignOut(ctx context.Context, userID uuid.UUID) error {
	ret := m.Called(ctx, userID)
	var r0 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// ValidateIDToken mocks concrete ValidateIDToken.
func (m *MockTokenService) ValidateIDToken(tokenString string) (*model.User, error) {
	ret := m.Called(tokenString)

	// First value passed to "Return".
	var r0 *model.User
	if ret.Get(0) != nil {
		// We can just return this if we know we won't be passing the function to "Return".
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// ValidateRefreshToken mocks concrete ValidateRefreshToken.
func (m *MockTokenService) ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error) {
	ret := m.Called(refreshTokenString)

	var r0 *model.RefreshToken
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.RefreshToken)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
