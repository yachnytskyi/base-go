package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/handler/middleware"
	"github.com/yachnytskyi/base-go/account/model"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// Handler struct holds required services for handler to function.
type Handler struct {
	UserService  model.UserService
	TokenService model.TokenService
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization.
type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	BaseURL         string
	TimeoutDuration time.Duration
}

// NewHandler initializes the handler with required injected services along with http routes.
// Does not return as it deals directly with a reference to the gin Engine.
func NewHandler(c *Config) {
	// Create a handler (with injected services).
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
	} // Currently has no properties.

	// Create an account group.
	g := c.R.Group(c.BaseURL)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
	}

	g.GET("/me", h.Me)
	g.POST("/signup", h.SignUp)
	g.POST("/signin", h.SignIn)
	g.POST("/signout", h.SignOut)
	g.POST("/tokens", h.Tokens)
	g.POST("/image", h.Image)
	g.DELETE("/image", h.DeleteImage)
	g.PUT("/details", h.Details)

}

// SignOut handler.
func (h *Handler) SignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's signout",
	})
}

// SignOut handler.
func (h *Handler) Tokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's tokens",
	})
}

// Image handler.
func (h *Handler) Image(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's image",
	})
}

// DeleteImage handler.
func (h *Handler) DeleteImage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's deleteimage",
	})
}

// Details handler.
func (h *Handler) Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's details",
	})
}
