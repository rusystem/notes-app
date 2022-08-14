package service

import (
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository"
)

type NoteService struct {
	repo repository.Note
}

func NewNoteService(repo repository.Note) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(userId int, note domain.Note) (int, error) {
	return s.repo.Create(userId, note)
}

func (s *NoteService) GetByID(userId, id int) (domain.Note, error) {
	return s.repo.GetByID(userId, id)
}

func (s *NoteService) GetAll(userId int) ([]domain.Note, error) {
	return s.repo.GetAll(userId)
}

func (s *NoteService) Delete(userId, id int) error {
	return s.repo.Delete(userId, id)
}

func (s *NoteService) Update(userId, id int, newNote domain.UpdateNote) error {
	if err := newNote.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, id, newNote)
}
