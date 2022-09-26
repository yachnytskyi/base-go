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
	// Setup.
	gin.SetMode(gin.TestMode)

	t.Run("Email and Password Required", func(t *testing.T) {
		// Show that it is not called in this case
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user.
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with empty email and password.
		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader.
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 400, responseRecorder.Code)
		mockUserService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Invalid email", func(t *testing.T) {
		// Show that it is not called in this case.
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user.
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with the wrong email and password.
		reqBody, err := json.Marshal(gin.H{
			"email":    "kostya@kostya",
			"password": "secretpassword",
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader.
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 400, responseRecorder.Code)
		mockUserService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Password is too short", func(t *testing.T) {
		// Show that it is not called in this case.
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user.
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with the wrong email and password.
		reqBody, err := json.Marshal(gin.H{
			"email":    "kostya@gmail.com",
			"password": "secr",
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader.
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 400, responseRecorder.Code)
		mockUserService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Password is too long", func(t *testing.T) {
		// Show that it is not called in this case.
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user.
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with the wrong email and password.
		reqBody, err := json.Marshal(gin.H{
			"email":    "kostya@gmail.com",
			"password": "secretpassworddasuhd89ydhuiasdajkdh792dyhuaksjdhnajdb78w",
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader.
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 400, responseRecorder.Code)
		mockUserService.AssertNotCalled(t, "SignUp")
	})

	t.Run("Error returned from UserService", func(t *testing.T) {
		user := &model.User{
			Email:    "kostya@kostya.com",
			Password: "secretpassword",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), user).Return(apperrors.NewConflict("User Already Exists", user.Email))

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Do not need a middleware as we don't yet have the authorized user.
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// Create a request body with the filed email and password.
		reqBody, err := json.Marshal(gin.H{
			"email":    user.Email,
			"password": user.Password,
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader.
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 409, responseRecorder.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		user := &model.User{
			Email:    "kostya@kostya.com",
			Password: "secretpassword",
		}

		mockTokenResponse := &model.TokenPair{
			IDToken:      model.IDToken{SignedString: "idToken"},
			RefreshToken: model.RefreshToken{SignedString: "refreshToken"},
		}

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), user).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*context.emptyCtx"), user, "").Return(mockTokenResponse, nil)

		// A response recorder for getting a written http response.
		responseRecorder := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have an authorized user.
		router := gin.Default()

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		// Create a request body with an empty email and password.
		requestBody, err := json.Marshal(gin.H{
			"email":    user.Email,
			"password": user.Password,
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader.
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		expectedResponseBody, err := json.Marshal(gin.H{
			"tokens": mockTokenResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, responseRecorder.Code)
		assert.Equal(t, expectedResponseBody, responseRecorder.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		user := &model.User{
			Email:    "kostya@kostya.com",
			Password: "secretpassword",
		}

		mockErrorResponse := apperrors.NewInternal()

		mockUserService := new(mocks.MockUserService)
		mockTokenService := new(mocks.MockTokenService)

		mockUserService.On("SignUp", mock.AnythingOfType("*context.emptyCtx"), user).Return(nil)
		mockTokenService.On("NewPairFromUser", mock.AnythingOfType("*context.emptyCtx"), user, "").Return(nil, mockErrorResponse)

		// A response recorder for getting a written http response.
		responseRecorder := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have an authorized user.
		router := gin.Default()

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		// Create a request body with an empty email and password.
		requestBody, err := json.Marshal(gin.H{
			"email":    user.Email,
			"password": user.Password,
		})
		assert.NoError(t, err)

		// Use bytes.NewBuffer to create a reader.
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecorder, request)

		expectedResponseBody, err := json.Marshal(gin.H{
			"error": mockErrorResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockErrorResponse.Status(), responseRecorder.Code)
		assert.Equal(t, expectedResponseBody, responseRecorder.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})
}
