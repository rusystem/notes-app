package domain

import "errors"

type Note struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateNote struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (un *UpdateNote) Validate() error {
	if un.Title == nil && un.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
