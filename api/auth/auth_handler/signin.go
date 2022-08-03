package auth_handler

import (
	"github.com/gin-gonic/gin"
	"goauth/user/user_domain"
	"goauth/utils"
	"goauth/utils/apperrors"
	"log"
	"net/http"
)

type signInReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignIn(ctx *gin.Context) {
	var req signInReq
	if ok := utils.BindData(ctx, &req); !ok {
		return
	}

	u := &user_domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	// SignIn Service
	user, err := h.UserService.GetByEmail(u.Email)

	if user == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": apperrors.NewBadRequest("check email or password"),
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": apperrors.NewInternal(),
		})
		return
	}

	match, err := h.PasswordService.ComparePasswordAndHash(u.Password, user.Password)

	if !match {
		ctx.JSON(http.StatusOK, gin.H{
			"error": apperrors.NewBadRequest("check email or password"),
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": apperrors.NewInternal(),
		})
		return
	}

	tokens, err := h.TokenService.NewPairFromUser(user, "")
	if err != nil {
		log.Printf("Failed to create tokens for user:%v\n", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"error": apperrors.NewInternal(),
		})
		return
	}
	log.Printf("tokensss %v", tokens)
	ctx.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
