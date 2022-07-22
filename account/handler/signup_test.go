package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
	"github.com/yachnytskyi/base-go/account/model/mocks"
)

func TestSignUp(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Email and Password Required", func(t *testing.T) {
		// Show that it is not called in this case
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Invalid email", func(t *testing.T) {
		// Show that it is not called in this case
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with the wrong email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "kostya@kostya",
			"password": "secretpassword",
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Password is too short", func(t *testing.T) {
		// Show that it is not called in this case
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with the wrong email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "kostya@gmail.com",
			"password": "secr",
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Password is too long", func(t *testing.T) {
		// Show that it is not called in this case
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with the wrong email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "kostya@gmail.com",
			"password": "secretpassworddasuhd89ydhuiasdajkdh792dyhuaksjdhnajdb78w",
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Error returned from UserService", func(t *testing.T) {
		u := &model.User{
			Email:    "kostya@kostya.com",
			Password: "password",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*gin.Context"), u).Return(apperrors.NewConflict("User Already Exists", u.Email))

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with the filed email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 409, rr.Code)
		mockUserService.AssertExpectations(t)
	})
}
