package auth_handler

import (
	"fmt"
	"goauth/utils/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Current(ctx *gin.Context) {

	// A *user.User will be added to gin context using a middleware
	u, exists := ctx.Get("userID")

	if !exists {
		log.Printf("unable to extract user from request context")
		err := apperrors.NewInternal()
		ctx.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	fetchedUser, err := h.UserService.Get(u.(uint))

	if err != nil {
		log.Printf("Unable to find user: %v\n%v\n", u, err)
		e := apperrors.NewNotFound("user", fmt.Sprintf("ID(%v)", u))

		ctx.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": fetchedUser,
	})
}
