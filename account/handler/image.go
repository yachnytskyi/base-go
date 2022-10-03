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
func (h *Handler) Image(context *gin.Context) {
	authUser := context.MustGet("user").(*model.User)

	// Limit overly large requests bodies.
	context.Request.Body = http.MaxBytesReader(context.Writer, context.Request.Body, h.MaxBodyBytes)

	imageFileHeader, err := context.FormFile("imageFile")

	// Check for an error before checking non-nil header.
	if err != nil {
		// Should be a validation error.
		log.Printf("Unable to parse multipart/form-data: %+v", err)

		if err.Error() == "http: the request body is too large" {
			context.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"erorr": fmt.Sprintf("Max request body size is %v bytes\n", h.MaxBodyBytes),
			})
			return
		}
		err := apperrors.NewBadRequest("Unable to parse multipart/form-data")
		context.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	if imageFileHeader == nil {
		err := apperrors.NewBadRequest("Must include an imageFile")
		context.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	mimeType := imageFileHeader.Header.Get("Content-Type")

	// Validate an image mime-type is allowable.
	if valid := isAllowedImageType(mimeType); !valid {
		log.Println("The image is not an allowable mime-type")
		err := apperrors.NewBadRequest("imageFile must be 'image/jpeg' or 'image/png'")
		context.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	ctx := context.Request.Context()

	updatedUser, err := h.UserService.SetProfileImage(ctx, authUser.UserID, imageFileHeader)
	if err != nil {
		context.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"imageURL": updatedUser.ImageURL,
		"message":  "success",
	})

}
