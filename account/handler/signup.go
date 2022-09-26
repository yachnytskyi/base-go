package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// signUpRequest is not exported, hence the lowercase name
// is is used for validation and json marshalling.
type signUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// SignUp handler.
func (h *Handler) SignUp(c *gin.Context) {
	// define a variable to which we'll bind incoming
	// json body, {email, password}.
	var jsonRequest signUpRequest

	// Bind incoming json to struct and check for validation errors.
	if ok := bindData(c, &jsonRequest); !ok {
		return
	}

	user := &model.User{
		Email:    jsonRequest.Email,
		Password: jsonRequest.Password,
	}

	ctx := c.Request.Context()
	err := h.UserService.SignUp(ctx, user)

	if err != nil {
		log.Printf("Failed to sign up the user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"erorr": err,
		})
		return
	}

	// Create token pair as strings.
	tokens, err := h.TokenService.NewPairFromUser(ctx, user, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		// May eventually implement rollback logic here
		// meaning, if we fail to create tokens after creating a user,
		// we make sure to cleate/delete the created user in the database.

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
