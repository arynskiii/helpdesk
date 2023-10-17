package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorization = "Authorization"
	userCtx       = "userId"
	roleCtx       = "roleCtx"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorization)
	if header == "" {
		c.JSON(401, "header is empty")
		return
	}
	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		c.JSON(http.StatusUnauthorized, "invalid with header")
		return
	}
	user, ok := h.services.GetUserByToken(headerParts[1])

	if user.Name == "admin" && user.Password == "admin" {
		c.Set(roleCtx, "admin")
	} else {
		c.Set(roleCtx, "common")
	}
	if ok != nil {
		c.JSON(400, ok)
	}

	c.Set(userCtx, user.Id)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")

	}
	idInt, ok := id.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "user id is invalid type")
		return 0, errors.New("user id is invalid type")
	}
	return idInt, nil
}
