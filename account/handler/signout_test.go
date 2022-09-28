package handler

import (
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

func TestSignOut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		contextUser := &model.User{
			UserID: userID,
			Email:  "kostya1@kostya.com",
		}

		// A response recorder for gettings written an http response.
		responseRecorder := httptest.NewRecorder()

		// Creates a test context for setting a user.
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", contextUser)
		})

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.On("SignOut", mock.AnythingOfType("*context.emptyCtx"), contextUser.UserID).Return(nil)

		NewHandler(&Config{
			Router:       router,
			TokenService: mockTokenService,
		})

		request, _ := http.NewRequest(http.MethodPost, "/signout", nil)
		router.ServeHTTP(responseRecorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"message": "the user signed out successfully!",
		})

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())
	})

	t.Run("SignOut Error", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		contextUser := &model.User{
			UserID: userID,
			Email:  "kostya2@kostya.com",
		}

		// A response recorder for getting written an http response.
		responseRecorder := httptest.NewRecorder()

		// Creates a test context for setting a user.
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", contextUser)
		})

		mockTokenService := new(mocks.MockTokenService)
		mockTokenService.On("SignOut", mock.AnythingOfType("*context.emptyCtx"), contextUser.UserID).Return(apperrors.NewInternal())

		NewHandler(&Config{
			Router:       router,
			TokenService: mockTokenService,
		})

		request, _ := http.NewRequest(http.MethodPost, "/signout", nil)
		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
	})
}
