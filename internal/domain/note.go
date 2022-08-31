package domain

type Note struct {
	ID          int    `json:"id"`
	Title       string `json:"title,max=55" example:"Title"`
	Description string `json:"description,max=255" example:"Description"`
}

type UpdateNote struct {
	Title       *string `json:"title" example:"Title"`
	Description *string `json:"description" example:"Description!"`
}

func (un UpdateNote) IsValid() bool {
	return un.Title != nil && un.Description != nil
}
