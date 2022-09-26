package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// Me handler calls services for getting
// a user's details.
func (h *Handler) Me(c *gin.Context) {
	// A *model.User will eventually be added to context in middleware.
	user, exists := c.Get("user")

	// This shouldn't happen, as our middleware ought to throw an error
	// This is an extra safety measure
	// We'll extract this logic later as it will be common to all handler
	// methods which require a valid user.
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	userID := user.(*model.User).UserID

	// Use the Request context.
	ctx := c.Request.Context()
	user, err := h.UserService.Get(ctx, userID)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", userID, err)
		e := apperrors.NewNotFound("user", userID.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
