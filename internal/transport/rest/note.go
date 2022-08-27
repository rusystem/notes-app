package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rusystem/notes-app/internal/domain"
	"net/http"
	"strconv"
)

// @Summary Create new note
// @Security ApiKeyAuth
// @Tags note
// @Description Create note
// @ID Create-note
// @Accept json
// @Produce json
// @Param input body domain.UpdateNote true "note info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/note [post]
func (h *Handler) create(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input domain.Note
	if err := c.BindJSON(&input); err != nil {
		domain.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Note.Create(c, userId, input)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get note by id
// @Security ApiKeyAuth
// @Tags note
// @Description Get note by id
// @ID Get-note-by-id
// @Accept json
// @Produce json
// @Param id path integer true "Note ID"
// @Success 200 {object} domain.Note
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/note/{id} [get]
func (h *Handler) getById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		domain.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
	}

	note, err := h.services.GetByID(c, userId, id)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, note)

}

// @Summary Get all notes
// @Security ApiKeyAuth
// @Tags note
// @Description Get all notes
// @ID Get-all-notes
// @Accept json
// @Produce json
// @Success 200 {object} domain.GetAllNoteResponse
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/note [get]
func (h *Handler) getAll(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	notes, err := h.services.Note.GetAll(c, userId)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, domain.GetAllNoteResponse{Data: notes})
}

// @Summary Delete note by id
// @Security ApiKeyAuth
// @Tags note
// @Description Delete note by id
// @ID Delete-note-by-id
// @Accept json
// @Produce json
// @Param id path integer true "Note ID"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/note/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		domain.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Delete(c, userId, id)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}

// @Summary Update note by id
// @Security ApiKeyAuth
// @Tags note
// @Description Update note by id
// @ID Update-note-by-id
// @Accept json
// @Produce json
// @Param id path integer true "Note ID"
// @Param input body domain.UpdateNote true "note info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/note/{id} [put]
func (h *Handler) update(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		domain.NewErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input domain.UpdateNote
	if err := c.BindJSON(&input); err != nil {
		domain.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Update(c, userId, id, input); err != nil {
		domain.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, domain.StatusResponse{Status: "ok"})
}
