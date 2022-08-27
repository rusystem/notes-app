package psql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rusystem/notes-app/internal/domain"
	"strings"
)

type NoteRepository struct {
	db *sqlx.DB
}

func NewNoteRepository(db *sqlx.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(ctx context.Context, userId int, note domain.Note) (int, error) {
	var id int

	createNoteQuery := fmt.Sprintf("INSERT INTO %s (uid, title, description) VALUES ($1, $2, $3) RETURNING id",
		notesTable)
	row := r.db.QueryRowContext(ctx, createNoteQuery, userId, note.Title, note.Description)
	err := row.Scan(&id)

	return id, err
}

func (r *NoteRepository) GetByID(ctx context.Context, userId, id int) (domain.Note, error) {
	var note domain.Note

	query := fmt.Sprintf("SELECT id, title, description FROM %s WHERE uid = $1 AND id = $2",
		notesTable)
	err := r.db.GetContext(ctx, &note, query, userId, id)

	return note, err
}

func (r *NoteRepository) GetAll(ctx context.Context, userId int) ([]domain.Note, error) {
	var notes []domain.Note

	query := fmt.Sprintf("SELECT id, title, description FROM %s WHERE uid = $1",
		notesTable)
	err := r.db.SelectContext(ctx, &notes, query, userId)

	return notes, err
}

func (r *NoteRepository) Delete(ctx context.Context, userId, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uid = $1 AND id = $2",
		notesTable)

	_, err := r.db.ExecContext(ctx, query, userId, id)

	return err
}

func (r *NoteRepository) Update(ctx context.Context, userId, id int, newNote domain.UpdateNote) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if newNote.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *newNote.Title)
		argId++
	}

	if newNote.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *newNote.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%d AND uid=$%d", notesTable, setQuery, argId, argId+1)
	args = append(args, id, userId)

	_, err := r.db.ExecContext(ctx, query, args...)

	return err
}
