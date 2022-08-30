package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rusystem/notes-app/internal/domain"
	"net/http"
)

// @Summary SignUp
// @Tags auth
// @Description Create new account
// @ID Create-account
// @Accept json
// @Produce json
// @Param input body domain.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input domain.User
	if err := c.BindJSON(&input); err != nil {
		domain.NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(c, input)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "credentials"
// @Success 200 {string} string "message"
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		domain.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cookie, err := h.services.SignIn(c, input)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(cookie.Name, cookie.Token, cookie.MaxAge, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "authorization was successful",
	})
}

// @Summary Logout
// @Tags auth
// @Description logout
// @ID logout
// @Accept json
// @Produce json
// @Success 200 {string} string "message"
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/logout [get]
func (h *Handler) logout(c *gin.Context) {
	token, err := c.Cookie(domain.AuthCookie)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if token == "" {
		domain.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if err = h.services.Logout(c, token); err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(domain.AuthCookie, token, -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "you have logged out",
	})
}
