package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// Image handler.
func (h *Handler) Image(c *gin.Context) {
	authUser := c.MustGet("user").(*model.User)

	// Limit overly large requests bodies.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

	imageFileHeader, err := c.FormFile("imageFile")

	// Check for an error before checking non-nil header.
	if err != nil {
		// Should be a validation error.
		log.Printf("Unable to parse multipart/form-data: %+v", err)

		if err.Error() == "http: the request body is too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"erorr": fmt.Sprintf("Max request body size is %v bytes\n", h.MaxBodyBytes),
			})
			return
		}
		err := apperrors.NewBadRequest("Unable to parse multipart/form-data")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	if imageFileHeader == nil {
		err := apperrors.NewBadRequest("Must include an imageFile")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	mimeType := imageFileHeader.Header.Get("Content-Type")

	// Validate an image mime-type is allowable.
	if valid := isAllowedImageType(mimeType); !valid {
		log.Println("The image is not an allowable mime-type")
		err := apperrors.NewBadRequest("imageFile must be 'image/jpeg' or 'image/png'")
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	ctx := c.Request.Context()

	updatedUser, err := h.UserService.SetProfileImage(ctx, authUser.UserID, imageFileHeader)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imageURL": updatedUser.ImageURL,
		"message":  "success",
	})

}
