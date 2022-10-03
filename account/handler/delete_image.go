package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// DeleteImage handler.
func (h *Handler) DeleteImage(context *gin.Context) {
	authUser := context.MustGet("user").(*model.User)

	ctx := context.Request.Context()
	err := h.UserService.ClearProfileImage(ctx, authUser.UserID)

	if err != nil {
		log.Printf("Failed to delete the profile image: %v\n", err.Error())

		context.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
