package auth_handler

import (
	"github.com/gin-gonic/gin"
	"goauth/utils"
	"goauth/utils/apperrors"
	"log"
	"net/http"
)

type tokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) Tokens(ctx *gin.Context) {
	var req tokenReq
	if ok := utils.BindData(ctx, &req); !ok {
		return
	}
	//verify jwt
	refreshToken, err := h.TokenService.ValidateRefreshToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	//	 get user
	u, err := h.UserService.Get(refreshToken.UID)
	if err != nil {
		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	//	Create fresh token pairs
	tokens, err := h.TokenService.NewPairFromUser(u, refreshToken.ID)

	if err != nil {
		log.Printf("Failed to create tokens for user:%v\n", err.Error())

		ctx.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
