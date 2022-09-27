package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
	"github.com/yachnytskyi/base-go/account/model/mocks"
)

func TestTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTokenService := new(mocks.MockTokenService)
	mockUserService := new(mocks.MockUserService)

	router := gin.Default()

	NewHandler(&Config{

		Router:       router,
		TokenService: mockTokenService,
		UserService:  mockUserService,
	})

	t.Run("Invalid request", func(t *testing.T) {
		// A response recorder for getting written an hhtp response.
		responseRecorder := httptest.NewRecorder()

		// Create a request body witn invalid fields.
		requestBody, _ := json.Marshal(gin.H{
			"notRefreshToken": "this key is not valid for this handler!",
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		mockTokenService.AssertNotCalled(t, "Get")
		mockUserService.AssertNotCalled(t, "Get")
		mockTokenService.AssertNotCalled(t, "NewPairFromUser")
	})

	t.Run("Invalid token", func(t *testing.T) {
		invalidTokenString := "invalid token string"
		mockErrorMessage := "authProblems"
		mockError := apperrors.NewAuthorization(mockErrorMessage)

		mockTokenService.On("ValidateRefreshToken", invalidTokenString).Return(nil, mockError)

		// A response recorder for getting written an http response.
		responseRecorder := httptest.NewRecorder()

		// Create a request body with invalid fields.
		requestBody, _ := json.Marshal(gin.H{
			"refreshToken": invalidTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(responseRecorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), responseRecorder.Code)
		assert.Equal(t, responseBody, responseBody, responseRecorder.Body.Bytes())
		mockTokenService.AssertCalled(t, "ValidateRefreshToken", invalidTokenString)
		mockUserService.AssertNotCalled(t, "Get")
		mockTokenService.AssertNotCalled(t, "NewPairFromUser")
	})

	t.Run("Failure to create new token pair", func(t *testing.T) {
		invalidTokenString := "invalid token"
		mockTokenID, _ := uuid.NewRandom()
		mockUserID, _ := uuid.NewRandom()

		mockRefreshTokenResponse := &model.RefreshToken{
			SignedString: invalidTokenString,
			ID:           mockTokenID,
			UserID:       mockUserID,
		}

		mockTokenService.On("ValidateRefreshToken", invalidTokenString).Return(mockRefreshTokenResponse, nil)

		mockUserResponse := &model.User{
			UserID: mockUserID,
		}
		getArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"), mockRefreshTokenResponse.UserID,
		}

		mockUserService.On("Get", getArguments...).Return(mockUserResponse, nil)

		mockError := apperrors.NewAuthorization("Invalid refresh token")
		newPairArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUserResponse,
			mockRefreshTokenResponse.ID.String(),
		}

		mockTokenService.On("NewPairFromUser", newPairArguments...).Return(nil, mockError)

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Create a request body with invalid fields.
		requestBody, _ := json.Marshal(gin.H{
			"refreshToken": invalidTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(responseRecorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())
		mockTokenService.AssertCalled(t, "ValidateRefreshToken", invalidTokenString)
		mockUserService.AssertCalled(t, "Get", getArguments...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", newPairArguments...)
	})

	t.Run("Success", func(t *testing.T) {
		validTokenString := "valid token string"
		mockTokenID, _ := uuid.NewRandom()
		mockUserID, _ := uuid.NewRandom()

		mockRefreshTokenResponse := &model.RefreshToken{
			SignedString: validTokenString,
			ID:           mockTokenID,
			UserID:       mockUserID,
		}

		mockTokenService.On("ValidateRefreshToken", validTokenString).Return(mockRefreshTokenResponse, nil)

		mockUserResponse := &model.User{
			UserID: mockUserID,
		}
		getArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"), mockRefreshTokenResponse.UserID,
		}

		mockUserService.On("Get", getArguments...).Return(mockUserResponse, nil)

		mockNewTokenID, _ := uuid.NewRandom()
		mockNewUserID, _ := uuid.NewRandom()
		mockTokenPairResponse := &model.TokenPair{
			IDToken: model.IDToken{SignedString: "newIDToken"},
			RefreshToken: model.RefreshToken{
				SignedString: "newRefreshToken",
				ID:           mockNewTokenID,
				UserID:       mockNewUserID,
			},
		}

		newPairArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"), mockUserResponse, mockRefreshTokenResponse.ID.String(),
		}

		mockTokenService.On("NewPairFromUser", newPairArguments...).Return(mockTokenPairResponse, nil)

		// A response recorder for getting written an http response.
		responseRecorder := httptest.NewRecorder()

		// Create a request body with valid fields.
		requestBody, _ := json.Marshal(gin.H{
			"refreshToken": validTokenString,
		})

		request, _ := http.NewRequest(http.MethodPost, "/tokens", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(responseRecorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"tokens": mockTokenPairResponse,
		})

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())

		mockTokenService.AssertCalled(t, "ValidateRefreshToken", validTokenString)
		mockUserService.AssertCalled(t, "Get", getArguments...)
		mockTokenService.AssertCalled(t, "NewPairFromUser", newPairArguments...)
	})

	// TODO - User not found (maybe in the furure).
}
