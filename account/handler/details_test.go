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

func TestDetails(t *testing.T) {
	// Setup.
	gin.SetMode(gin.TestMode)

	userID, _ := uuid.NewRandom()
	contextUser := &model.User{
		UserID: userID,
	}

	router := gin.Default()
	router.Use(func(context *gin.Context) {
		context.Set("user", contextUser)
	})

	mockUserService := new(mocks.MockUserService)

	NewHandler(&Config{
		Router:      router,
		UserService: mockUserService,
	})

	t.Run("Data binding error", func(t *testing.T) {
		responseRecoder := httptest.NewRecorder()

		requestBody, _ := json.Marshal(gin.H{
			"email": "notanemail",
		})
		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(responseRecoder, request)

		assert.Equal(t, http.StatusBadRequest, responseRecoder.Code)
		mockUserService.AssertNotCalled(t, "UpdateDetails")
	})

	t.Run("Update success", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		newUsername := "Constantine"
		newEmail := "constantine@constantin.com"
		newWebsite := "https://constantine.com"

		requestBody, _ := json.Marshal(gin.H{
			"username": newUsername,
			"email":    newEmail,
			"website":  newWebsite,
		})

		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")

		userToUpdate := &model.User{
			UserID:   contextUser.UserID,
			Username: newUsername,
			Email:    newEmail,
			Website:  newWebsite,
		}

		updateArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userToUpdate,
		}

		dbImageURL := "https://constantin.com/static/2185715217/847a431/Image.jpg"

		mockUserService.On("UpdateDetails", updateArguments...).
			Run(func(args mock.Arguments) {
				userArgument := args.Get(1).(*model.User) // Arg 0 is context, arg 1 is *User.
				userArgument.ImageURL = dbImageURL
			}).
			Return(nil)

		router.ServeHTTP(responseRecorder, request)

		userToUpdate.ImageURL = dbImageURL
		responseBody, _ := json.Marshal(gin.H{
			"user": userToUpdate,
		})

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())
		mockUserService.AssertCalled(t, "UpdateDetails", updateArguments...)
	})

	t.Run("Update failure", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		newUsername := "Constantine"
		newEmail := "constantine@constantine.com"
		newWebsite := "https://constantine.com"

		requestBody, _ := json.Marshal(gin.H{
			"username": newUsername,
			"email":    newEmail,
			"website":  newWebsite,
		})

		request, _ := http.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(requestBody))
		request.Header.Set("Content-Type", "application/json")

		userToUpdate := &model.User{
			UserID:   contextUser.UserID,
			Username: newUsername,
			Email:    newEmail,
			Website:  newWebsite,
		}

		updateArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			userToUpdate,
		}

		mockError := apperrors.NewInternal()

		mockUserService.On("UpdateDetails", updateArguments...).Return(mockError)

		router.ServeHTTP(responseRecorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())
		mockUserService.AssertCalled(t, "UpdateDetails", updateArguments...)
	})
}
