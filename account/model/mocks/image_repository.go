package mocks

import (
	"context"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockImageRepository struct {
	mock.Mock
}

// DeleteProfile is a mock of representantions of ImageRepository Delete Profile.
func (m *MockImageRepository) DeleteProfile(ctx context.Context, objectName string) error {
	ret := m.Called(ctx, objectName)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// UpdateProfile is a mock of representantions of ImageRepository Update Profile.
func (m *MockImageRepository) UpdateProfile(ctx context.Context, objectName string, imageFile multipart.File) (string, error) {
	// Arguments that will be passed to "Return" in the tests, when function
	// is called with a userID. Hence the name "ret".
	ret := m.Called(ctx, objectName, imageFile)

	// First value passed to "Return".
	var r0 string
	if ret.Get(0) != nil {
		// We can just return this if we know we will not be passing function to "Return".
		r0 = ret.Get(0).(string)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
