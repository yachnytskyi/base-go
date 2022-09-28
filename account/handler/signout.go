package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// SignOut handler.
func (h *Handler) SignOut(c *gin.Context) {
	user := c.MustGet("user")

	ctx := c.Request.Context()
	if err := h.TokenService.SignOut(ctx, user.(*model.User).UserID); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "the user signed out successfully!",
	})
}
