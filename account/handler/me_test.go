package handler

import (
	"encoding/json"
	"fmt"
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

func TestMe(t *testing.T) {
	// Setup.
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		userID, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UserID:   userID,
			Email:    "kostya@kostya.com",
			Username: "Kostya Kostyan",
		}

		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.AnythingOfType("*context.emptyCtx"), userID).Return(mockUserResp, nil)

		// A response recorder for etting written http response.
		responseRecorder := httptest.NewRecorder()

		// Use a middleware to set context for test
		// the only claims we care about in this test
		// is the UserID.
		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", &model.User{
				UserID: userID,
			},
			)
		})
		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(responseRecorder, request)

		expectedResponseBody, err := json.Marshal(gin.H{
			"user": mockUserResp,
		})
		assert.NoError(t, err)

		assert.Equal(t, 200, responseRecorder.Code)
		assert.Equal(t, expectedResponseBody, responseRecorder.Body.Bytes())
		mockUserService.AssertExpectations(t) // Assert that UserService.Get was called.
	})

	t.Run("NoContextUser", func(t *testing.T) {
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.Anything, mock.Anything).Return(nil, nil)

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		// Do not append user to context.
		router := gin.Default()
		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, 500, responseRecorder.Code)
		mockUserService.AssertNotCalled(t, "Get", mock.Anything)
	})

	t.Run("NotFound", func(t *testing.T) {
		userID, _ := uuid.NewRandom()
		mockUserService := new(mocks.MockUserService)
		mockUserService.On("Get", mock.Anything, userID).Return(nil, fmt.Errorf("Some error down call chain"))

		// A response recorder for getting written http response.
		responseRecorder := httptest.NewRecorder()

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", &model.User{
				UserID: userID,
			},
			)
		})

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(responseRecorder, request)

		expectedResponseError := apperrors.NewNotFound("user", userID.String())

		expectedrResponseBody, err := json.Marshal(gin.H{
			"error": expectedResponseError,
		})
		assert.NoError(t, err)

		assert.Equal(t, expectedResponseError.Status(), responseRecorder.Code)
		assert.Equal(t, expectedrResponseBody, responseRecorder.Body.Bytes())
		mockUserService.AssertExpectations(t) // Assert that UserService.Get was called.
	})
}
