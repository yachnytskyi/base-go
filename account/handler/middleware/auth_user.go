package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

// Used to help extract validation errors.
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// AuthUser extracts a user from the Authorization header
// which is of the form "Bearer token".
// It sets the user to the context if the user exists.
func AuthUser(s model.TokenService) gin.HandlerFunc {
	return func(context *gin.Context) {
		h := authHeader{}

		// Bind authorization Header to h and check for validation errors.
		if err := context.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				// We used this type in bind_data to extract desired fields from errs.
				// You might consider extracting it.
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

				context.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				context.Abort()
				return
			}

			// Otherwise erorr type is uknown.
			err := apperrors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": err,
			})

			context.Abort()
			return
		}

		idTokenHeader := strings.Split(h.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

			context.JSON(err.Status(), gin.H{
				"eror": err,
			})
			context.Abort()
			return
		}

		// Validate ID token here.
		user, err := s.ValidateIDToken(idTokenHeader[1])

		if err != nil {
			err := apperrors.NewAuthorization("Provided token is invalid")
			context.JSON(err.Status(), gin.H{
				"error": err,
			})
			context.Abort()
			return
		}

		context.Set("user", user)

		context.Next()
	}
}
