package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/rusystem/cache"
	"github.com/rusystem/notes-app/internal/config"
	"github.com/rusystem/notes-app/internal/domain"
	"github.com/rusystem/notes-app/internal/repository"
	"github.com/rusystem/notes-app/internal/transport"
	"github.com/sirupsen/logrus"
	"time"
)

type NoteService struct {
	cfg   *config.Config
	cache *cache.Cache
	repo  repository.Note
	mq    *transport.Server
}

func NewNoteService(cfg *config.Config, cache *cache.Cache, repo repository.Note, mq *transport.Server) *NoteService {
	return &NoteService{cfg, cache, repo, mq}
}

func (s *NoteService) Create(ctx context.Context, userId int, note domain.Note) (int, error) {
	id, err := s.repo.Create(ctx, userId, note)
	if err != nil {
		return 0, err
	}

	s.cache.Set(fmt.Sprintf("%d.%d", userId, id), domain.Note{ID: id, Title: note.Title, Description: note.Description}, s.cfg.Cache.Ttl)

	if err := s.mq.Publisher(domain.LogItem{
		Entity:    domain.NOTE,
		Action:    domain.CREATE,
		EntityID:  int64(userId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "note.Create",
		}).Error("failed to send log request:", err)
	}

	return id, nil
}

func (s *NoteService) GetByID(ctx context.Context, userId, id int) (domain.Note, error) {
	note, err := s.cache.Get(fmt.Sprintf("%d.%d", userId, id))
	if err == nil {
		return note.(domain.Note), nil
	}

	note, err = s.repo.GetByID(ctx, userId, id)
	if err != nil {
		return domain.Note{}, err
	}
	s.cache.Set(fmt.Sprintf("%d.%d", userId, id), note, s.cfg.Cache.Ttl)

	if err := s.mq.Publisher(domain.LogItem{
		Entity:    domain.NOTE,
		Action:    domain.GET,
		EntityID:  int64(userId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "note.GetById",
		}).Error("failed to send log request:", err)
	}

	return note.(domain.Note), nil
}

func (s *NoteService) GetAll(ctx context.Context, userId int) ([]domain.Note, error) {
	if err := s.mq.Publisher(domain.LogItem{
		Entity:    domain.NOTE,
		Action:    domain.GET,
		EntityID:  int64(userId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "note.GetAll",
		}).Error("failed to send log request:", err)
	}

	return s.repo.GetAll(ctx, userId)
}

func (s *NoteService) Delete(ctx context.Context, userId, id int) error {
	if err := s.mq.Publisher(domain.LogItem{
		Entity:    domain.NOTE,
		Action:    domain.DELETE,
		EntityID:  int64(userId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "note.Delete",
		}).Error("failed to send log request:", err)
	}

	return s.repo.Delete(ctx, userId, id)
}

func (s *NoteService) Update(ctx context.Context, userId, id int, newNote domain.UpdateNote) error {
	if !newNote.IsValid() {
		return errors.New("update structure has no values")
	}
	s.cache.Set(fmt.Sprintf("%d.%d", userId, id), domain.Note{ID: id, Title: *newNote.Title, Description: *newNote.Description}, s.cfg.Cache.Ttl)

	if err := s.mq.Publisher(domain.LogItem{
		Entity:    domain.NOTE,
		Action:    domain.UPDATE,
		EntityID:  int64(userId),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "note.Update",
		}).Error("failed to send log request:", err)
	}

	return s.repo.Update(ctx, userId, id, newNote)
}
