package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// Omitempty must be listed first (tags evaluated sequentially, it seems).
type detailsRequest struct {
	Username string `json:"username" binding:"omitempty,max=40"`
	Email    string `json:"email" binding:"required,email"`
	Website  string `json:"website" binding:"omitempty,url"`
}

// Details handler.
func (h *Handler) Details(c *gin.Context) {
	authUser := c.MustGet("user").(*model.User)

	var request detailsRequest

	if ok := bindData(c, &request); !ok {
		return
	}

	// Should be returned with current imageURL.
	user := &model.User{
		UserID:   authUser.UserID,
		Username: request.Username,
		Email:    request.Email,
		Website:  request.Website,
	}

	ctx := c.Request.Context()
	err := h.UserService.UpdateDetails(ctx, user)

	if err != nil {
		log.Printf("Failed to update the user: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
