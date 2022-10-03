package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
	"github.com/yachnytskyi/base-go/account/model/mocks"
)

func TestAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTokenService := new(mocks.MockTokenService)

	userID, _ := uuid.NewRandom()
	user := &model.User{
		UserID: userID,
		Email:  "kostya@kostya.com",
	}

	// Since we mock tokenService, we do not
	// need to create actual JWTs.
	validTokenHeader := "validTokenString"
	invalidTokenHeader := "invalidTokenString"
	invalidTokenError := apperrors.NewAuthorization("Unable to verify the user from idToken")

	mockTokenService.On("ValidateIDToken", validTokenHeader).Return(user, nil)
	mockTokenService.On("ValidateIDToken", invalidTokenHeader).Return(nil, invalidTokenError)

	t.Run("Adds a user to context", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		// Creates a test context and gin engine.
		_, testContext := gin.CreateTestContext(responseRecorder)

		// Will be populated with a user in a handler
		// if AuthUser middleware is successful.
		var contextUser *model.User

		// See this issue - https://github.com/gin-gonic/gin/issues/323
		// https://github.com/gin-gonic/gin/blob/master/auth_test.go#L91-L126
		// We create a handler to return "a user added to context" as this
		// is the only way to test modified context.
		testContext.GET("/me", AuthUser(mockTokenService), func(context *gin.Context) {
			contextKeyValue, _ := context.Get("user")
			contextUser = contextKeyValue.(*model.User)
		})

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validTokenHeader))
		testContext.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Equal(t, user, contextUser)

		mockTokenService.AssertCalled(t, "ValidateIDToken", validTokenHeader)

	})

	t.Run("Invalid Token", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		// Creates a test context and gin engine.
		_, testContext := gin.CreateTestContext(responseRecorder)
		testContext.GET("/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", invalidTokenHeader))
		testContext.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		mockTokenService.AssertCalled(t, "ValidateIDToken", invalidTokenHeader)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()

		// Creates a test context and gin engine.
		_, testContext := gin.CreateTestContext(responseRecorder)

		testContext.GET("/me", AuthUser(mockTokenService))

		request, _ := http.NewRequest(http.MethodGet, "/me", http.NoBody)

		testContext.ServeHTTP(responseRecorder, request)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		mockTokenService.AssertNotCalled(t, "ValidateIDToken")
	})
}
