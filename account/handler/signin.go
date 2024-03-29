package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// signInRequest is not exported.
type signInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// SignIn used to authenticate extant user.
func (h *Handler) SignIn(context *gin.Context) {
	var req signInRequest

	if ok := bindData(context, &req); !ok {
		return
	}

	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := context.Request.Context()
	err := h.UserService.SignIn(ctx, user)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())
		context.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.TokenService.NewPairFromUser(ctx, user, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		context.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
