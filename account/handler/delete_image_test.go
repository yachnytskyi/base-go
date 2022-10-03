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

func TestDeleteImage(t *testing.T) {
	// Setup.
	gin.SetMode(gin.TestMode)

	// An authorized middleware user.
	userID, _ := uuid.NewRandom()
	contextUser := &model.User{
		UserID: userID,
	}

	router := gin.Default()
	router.Use(func(context *gin.Context) {
		context.Set("user", contextUser)
	})

	// This handler requires UserService.
	mockUserService := new(mocks.MockUserService)

	NewHandler(&Config{
		Router:      router,
		UserService: mockUserService,
	})

	t.Run("Clear profile image error", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		clearProfileImageArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			contextUser.UserID,
		}

		errorResponse := apperrors.NewInternal()
		mockUserService.On("ClearProfileImage", clearProfileImageArguments...).Return(errorResponse)

		request, _ := http.NewRequest(http.MethodDelete, "/image", nil)
		router.ServeHTTP(responseRecorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"error": errorResponse,
		})

		assert.Equal(t, apperrors.Status(errorResponse), responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())
		mockUserService.AssertCalled(t, "ClearProfileImage", clearProfileImageArguments...)
	})

	t.Run("Success", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		// An authorized middleware user - overwriting for unique mock arguments.
		userId, _ := uuid.NewRandom()
		contextUser := &model.User{
			UserID: userId,
		}

		router := gin.Default()
		router.Use(func(context *gin.Context) {
			context.Set("user", contextUser)
		})

		// This handler requires a UserService.
		mockUserService := new(mocks.MockUserService)

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		clearProfileImageArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			contextUser.UserID,
		}

		mockUserService.On("ClearProfileImage", clearProfileImageArguments...).Return(nil)

		request, _ := http.NewRequest(http.MethodDelete, "/image", nil)
		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		mockUserService.AssertCalled(t, "ClearProfileImage", clearProfileImageArguments...)
	})
}
