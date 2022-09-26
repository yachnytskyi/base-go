package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

type tokensRequest struct {
	RefreshTokenString string `json:"refreshToken" binding:"required"`
}

// Tokens handler.
func (h *Handler) Tokens(c *gin.Context) {
	// Bind JSON to request of type tokensRequest.
	var request tokensRequest

	if ok := bindData(c, &request); !ok {
		return
	}

	ctx := c.Request.Context()

	// Verify refresh JWT.
	refreshToken, err := h.TokenService.ValidateRefreshToken(request.RefreshTokenString)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Get up-to-date user.
	user, err := h.UserService.Get(ctx, refreshToken.UserID)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Create fresh pair of tokens.
	tokens, err := h.TokenService.NewPairFromUser(ctx, user, refreshToken.ID.String())

	if err != nil {
		log.Printf("Failed to create tokens for the user: %+v. Error: %v\n", user, err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
