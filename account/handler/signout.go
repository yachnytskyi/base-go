package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// SignOut handler.
func (h *Handler) SignOut(context *gin.Context) {
	user := context.MustGet("user")

	ctx := context.Request.Context()
	if err := h.TokenService.SignOut(ctx, user.(*model.User).UserID); err != nil {
		context.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "the user signed out successfully!",
	})
}
