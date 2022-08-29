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

type signInInput struct {
	Username string `json:"username" binding:"required,min=2" example:"mdmitry"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		domain.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(c, input.Username, input.Password)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) logout(c *gin.Context) {

}
