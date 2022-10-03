package handler

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// used to help extract validation errors.
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// bindData is the helper function, returns false if data is not bound.
func bindData(context *gin.Context, request interface{}) bool {
	if context.ContentType() != "application/json" {
		message := fmt.Sprintf("%s only accepts Content-Type application/json", context.FullPath())

		err := apperrors.NewUnsupportedMediaType(message)

		context.JSON(err.Status(), gin.H{
			"error": err,
		})
		return false
	}

	// bindData incoming json to struct and check for validation errors.
	if err := context.ShouldBind(request); err != nil {
		log.Printf("Error binding data: %+v\n", err)

		if errs, ok := err.(validator.ValidationErrors); ok {
			// Could probably extract this, it is also in middleware_auth_user.
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
			return false
		}

		// Later I will add code for validating max body size here.

		// If a server is not able to properly extract validation errors,
		// it will fallback and return an internal server error.
		fallBack := apperrors.NewInternal()

		context.JSON(fallBack.Status(), gin.H{"error": fallBack})
		return false
	}
	return true
}
