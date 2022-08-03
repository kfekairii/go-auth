package auth_handler

import (
	"github.com/gin-gonic/gin"
	"goauth/auth/auth_domain"
	"goauth/auth/auth_middleware"
	"goauth/user/user_domain"

	"net/http"
)

// Handler Holds required services
type Handler struct {
	UserService     user_domain.IUserService
	TokenService    auth_domain.ITokenService
	PasswordService auth_domain.IPasswordService
}

// Config Hold services that will be injected at the initialization of the handler
type Config struct {
	RouterGroup     *gin.RouterGroup
	UserService     user_domain.IUserService
	PasswordService auth_domain.IPasswordService
	TokenService    auth_domain.ITokenService
}

func NewHandler(c *Config) {
	h := &Handler{
		UserService:     c.UserService,
		TokenService:    c.TokenService,
		PasswordService: c.PasswordService,
	}

	// create Auth group
	g := c.RouterGroup.Group("/auth")

	g.GET("/current", auth_middleware.AuthUser(c.TokenService), h.Current)
	g.POST("/signup", h.Signup)
	g.POST("/signin", h.SignIn)
	g.POST("/signout", h.SignOut)
	g.POST("/tokens", h.Tokens)
}

func (h *Handler) SignOut(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"signout": "user",
	})
}
