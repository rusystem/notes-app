package rest

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rusystem/notes-app/internal/domain"
	"net/http"
)

const (
	cookie  = "session"
	userCtx = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get(cookie)

	if userId == nil {
		domain.NewErrorResponse(c, http.StatusUnauthorized, "user not logged in")
		c.Abort()
		return
	}

	c.Set(userCtx, userId)
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
