package auth_middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goauth/auth/auth_domain"
	"goauth/utils/apperrors"
	"strings"
)

type authHeader struct {
	AccessToken string `header:"Authorization"`
}

type invalidArguments struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser(s auth_domain.ITokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := &authHeader{}
		//	Bind Authorization Header to h
		if err := c.ShouldBindHeader(h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []invalidArguments
				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArguments{
						Field: err.Field(),
						Value: err.Value().(string),
						Tag:   err.Tag(),
						Param: err.Param(),
					})
				}

				err := apperrors.NewBadRequest("Invalid request")
				c.JSON(err.Status(), gin.H{"error": err, "invalidArgs": invalidArgs})
				c.Abort()
				return
			}
		}
		accessTokenHeader := strings.Split(h.AccessToken, "Bearer ")

		if len(accessTokenHeader) < 2 {
			err := apperrors.NewAuthorization("Must provide Authorization header with format 'Bearer {token}'")
			c.JSON(err.Status(), gin.H{"error": err})
			c.Abort()
			return
		}
		//	Validate token
		user, err := s.ValidateAccessToken(accessTokenHeader[1])

		if err != nil {
			err := apperrors.NewAuthorization("invalid token")
			c.JSON(err.Status(), gin.H{"error": err})
			c.Abort()
			return
		}

		fmt.Printf("%v", user)

		c.Set("userID", user)
		c.Next()
	}
}
