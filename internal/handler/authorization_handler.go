package handler

import (
	"fmt"
	"net/http"

	"github.com/arynskiii/help_desk/models"
	"github.com/arynskiii/help_desk/pkg/logger"
	"github.com/gin-gonic/gin"
)

type signIn struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *Handler) signUp(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	fmt.Println(user)
	id, err := h.services.Authorization.CreateUser(c, user)
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create user(HANDLER): %w", err))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(ctx *gin.Context) {
	var signIn signIn
	if err := ctx.BindJSON(&signIn); err != nil {

		ctx.JSON(400, err)
		return
	}
	res, err := h.services.GenerateToken(signIn.Name, signIn.Password)
	if err != nil {
		logger.Error(err)
		ctx.JSON(400, err)
		return
	}

	ctx.JSON(http.StatusOK, res.Token)
}
