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
	"github.com/yachnytskyi/base-go/account/model/fixture"
	"github.com/yachnytskyi/base-go/account/model/mocks"
)

func TestImage(t *testing.T) {
	// Setup.
	gin.SetMode(gin.TestMode)

	userID, _ := uuid.NewRandom()
	contextUser := model.User{
		UserID: userID,
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user", &contextUser)
	})

	mockUserService := new(mocks.MockUserService)

	NewHandler(&Config{
		Router:       router,
		UserService:  mockUserService,
		MaxBodyBytes: 4 * 1024 * 1024,
	})

	t.Run("Success", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		imageURL := "https://www.imageURL.com/9362"

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()

		setProfileImageArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			contextUser.UserID,
			mock.AnythingOfType("*multipart.FileHeader"),
		}

		updatedUser := contextUser
		updatedUser.ImageURL = imageURL

		mockUserService.On("SetProfileImage", setProfileImageArguments...).Return(&updatedUser, nil)

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", multipartImageFixture.ContentType)

		router.ServeHTTP(responseRecorder, request)

		responseBody, _ := json.Marshal(gin.H{
			"imageURL": imageURL,
			"message":  "success",
		})

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, responseBody, responseRecorder.Body.Bytes())

		mockUserService.AssertCalled(t, "SetProfileImage", setProfileImageArguments...)
	})

	t.Run("Disaloowed mimetype", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		multipartImageFixture := fixture.NewMultipartImage("image.txt", "mage/svg+xml")
		defer multipartImageFixture.Close()

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", "multipart/form-data")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

		mockUserService.AssertNotCalled(t, "SetProfileImage")
	})

	t.Run("No image file provided", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		request, _ := http.NewRequest(http.MethodPost, "/image", nil)
		request.Header.Set("Content-Type", "multipart-form-data")

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

		mockUserService.AssertNotCalled(t, "SetProfileImage")
	})

	t.Run("Error from SetProfileImage", func(t *testing.T) {
		// Create a unique context user for this test.
		userID, _ := uuid.NewRandom()
		contextUser := model.User{
			UserID: userID,
		}

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user", &contextUser)
		})

		mockUserService := new(mocks.MockUserService)

		NewHandler(&Config{
			Router:       router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		responseRecorder := httptest.NewRecorder()

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()

		setProfileImageArguments := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			contextUser.UserID,
			mock.AnythingOfType("*multipart.FileHeader"),
		}

		mockError := apperrors.NewInternal()

		mockUserService.On("SetProfileImage", setProfileImageArguments...).Return(nil, mockError)

		request, _ := http.NewRequest(http.MethodPost, "/image", multipartImageFixture.MultipartBody)
		request.Header.Set("Content-Type", multipartImageFixture.ContentType)

		router.ServeHTTP(responseRecorder, request)

		assert.Equal(t, apperrors.Status(mockError), responseRecorder.Code)

		mockUserService.AssertCalled(t, "SetProfileImage", setProfileImageArguments...)
	})

	// TODO - how to handle large files? Creating large files is very slow
	// maybe create a byte slice and dupe Go into thinking it's an image...?
}
