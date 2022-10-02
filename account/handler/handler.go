package handler

import (
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
	MaxBodyBytes int64
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization.
type Config struct {
	Router          *gin.Engine
	UserService     model.UserService
	TokenService    model.TokenService
	BaseURL         string
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

// NewHandler initializes the handler with required injected services along with http routes.
// Does not return as it deals directly with a reference to the gin Engine.
func NewHandler(c *Config) {
	// Create a handler (with injected services).
	h := &Handler{
		UserService:  c.UserService,
		TokenService: c.TokenService,
		MaxBodyBytes: c.MaxBodyBytes,
	} // Currently has no properties.

	// Create an account group.
	g := c.Router.Group(c.BaseURL)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
		g.GET("/me", middleware.AuthUser(h.TokenService), h.Me)
		g.POST("/signout", middleware.AuthUser(h.TokenService), h.SignOut)
		g.PUT("/details", middleware.AuthUser(h.TokenService), h.Details)
		g.POST("/image", middleware.AuthUser(h.TokenService), h.Image)
		g.DELETE("/image", middleware.AuthUser(h.TokenService), h.DeleteImage)

	} else {
		g.GET("/me", h.Me)
		g.POST("/signout", h.SignOut)
		g.PUT("/details", h.Details)
		g.POST("/image", h.Image)
		g.DELETE("/image", h.DeleteImage)

	}

	g.POST("/signup", h.SignUp)
	g.POST("/signin", h.SignIn)
	g.POST("/tokens", h.Tokens)
}
