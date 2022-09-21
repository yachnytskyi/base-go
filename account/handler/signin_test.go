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

func TestSignIn(t *testing.T) {
	// Setup.
	gin.SetMode(gin.TestMode)

	// Setup mock services, gin engine/router, handler layer.
	mockUserService := new(mocks.MockUserService)
	mockTokenService := new(mocks.MockTokenService)

	router := gin.Default()

	NewHandler(&Config{
		R:            router,
		UserService:  mockUserService,
		TokenService: mockTokenService,
	})

	t.Run("Bad request data", func(t *testing.T) {
		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Create a request body with invalid fields.
		requestBody, err := json.Marshal(gin.H{
			"email":    "notanemail",
			"password": "shortpassword",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		mockUserService.AssertNotCalled(t, "SignIn")
		mockTokenService.AssertNotCalled(t, "NewTokensFromUser")
	})

	t.Run("Error Returned from UserService.SignIn", func(t *testing.T) {
		email := "kostya@kostya.com"
		password := "passworddoesnotmatch123"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
		}

		// So we can check for a known status code.
		mockError := apperrors.NewAuthorization("invalid email/password combo")

		mockUserService.On("SignIn", mockUSArgs...).Return(mockError)

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Create a request body with valid fields.
		requestBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(responseRecorder, request)

		mockUserService.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTokenService.AssertNotCalled(t, "NewTokensFromUser")
		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		email := "kostya@kostya.com"
		password := "passwordworksgreat123"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
		}

		mockUserService.On("SignIn", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
			"",
		}

		mockTokenPair := &model.TokenPair{
			IDToken:      model.IDToken{SignedString: "idToken"},
			RefreshToken: model.RefreshToken{SignedString: "refreshToken"},
		}

		mockTokenService.On("NewPairFromUser", mockTSArgs...).Return(mockTokenPair, nil)

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Create a request body with valid fields.
		requestBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(responseRecorder, request)

		respBody, err := json.Marshal(gin.H{
			"tokens": mockTokenPair,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, respBody, responseRecorder.Body.Bytes())

		mockUserService.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		email := "cannotproducetoken@kostya.com"
		password := "cannotproducetoken"

		mockUSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
		}

		mockUserService.On("SignIn", mockUSArgs...).Return(nil)

		mockTSArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			&model.User{Email: email, Password: password},
			"",
		}

		mockError := apperrors.NewInternal()
		mockTokenService.On("NewPairFromUser", mockTSArgs...).Return(nil, mockError)
		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Create a request body with valid fields.
		requestBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(responseRecorder, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), responseRecorder.Code)
		assert.Equal(t, respBody, responseRecorder.Body.Bytes())

		mockUserService.AssertCalled(t, "SignIn", mockUSArgs...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", mockTSArgs...)
	})
}
