package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
)

// MockTokenService is a mock type for model.TokenService.
type MockTokenService struct {
	mock.Mock
}

// NewPairFromUser mocks concrete NewPairFromUser.
func (m *MockTokenService) NewPairFromUser(ctx context.Context, u *model.User, refreshTokenID string) (*model.TokenPair, error) {
	ret := m.Called(ctx, u, refreshTokenID)

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
