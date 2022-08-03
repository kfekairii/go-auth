package auth_handler

import (
	"goauth/user/user_domain"
	"goauth/utils"
	"goauth/utils/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req signupReq

	if ok := utils.BindData(c, &req); !ok {
		return
	}

	// Signup service
	passwordHash, err := h.PasswordService.CreateHash(req.Password)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	u := &user_domain.User{
		Email:    req.Email,
		Password: passwordHash,
	}

	err = h.UserService.Create(u)

	if err != nil {
		log.Printf("Couldn't create a user:\n%v\n", err)
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.TokenService.NewPairFromUser(u, "")

	if err != nil {
		log.Printf("Failed to create tokens for user:%v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
