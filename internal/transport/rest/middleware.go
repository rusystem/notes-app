package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/notes-app/internal/domain"
	"net/http"
)

const (
	userCtx = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	token, err := c.Cookie(domain.AuthCookie)
	if err != nil || token == "" {
		domain.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	userId, err := h.services.GetSession(c, token)
	if err != nil || userId == 0 {
		domain.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	c.Set(userCtx, userId)
	c.Next()
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
