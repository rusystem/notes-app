package psql

import (
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

func (r *NoteRepository) Create(userId int, note domain.Note) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createNoteQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", notesTable)
	row := tx.QueryRow(createNoteQuery, note.Title, note.Description)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}

		return 0, err
	}

	createUsersNotesQuery := fmt.Sprintf("INSERT INTO %s (user_id, note_id) VALUES ($1, $2)", usersNotesTable)
	_, err = tx.Exec(createUsersNotesQuery, userId, id)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, err
		}

		return 0, err
	}

	return id, tx.Commit()
}

func (r *NoteRepository) GetByID(userId, id int) (domain.Note, error) {
	var note domain.Note

	query := fmt.Sprintf("SELECT n.id, n.title, n.description FROM %s n INNER JOIN %s un on n.id = un.note_id WHERE un.user_id = $1 AND un.note_id = $2",
		notesTable, usersNotesTable)
	err := r.db.Get(&note, query, userId, id)

	return note, err
}

func (r *NoteRepository) GetAll(userId int) ([]domain.Note, error) {
	var notes []domain.Note

	query := fmt.Sprintf("SELECT n.id, n.title, n.description FROM %s n INNER JOIN %s un on n.id = un.note_id WHERE un.user_id = $1",
		notesTable, usersNotesTable)
	err := r.db.Select(&notes, query, userId)

	return notes, err
}

func (r *NoteRepository) Delete(userId, id int) error {
	query := fmt.Sprintf("DELETE FROM %s n USING %s un WHERE n.id = un.note_id AND un.user_id=$1 AND un.note_id=$2",
		notesTable, usersNotesTable)

	_, err := r.db.Exec(query, userId, id)

	return err
}

func (r *NoteRepository) Update(userId, id int, newNote domain.UpdateNote) error {
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

	query := fmt.Sprintf("UPDATE %s n SET %s FROM %s un WHERE n.id = un.note_id AND un.note_id=$%d AND un.user_id=$%d",
		notesTable, setQuery, usersNotesTable, argId, argId+1)
	args = append(args, id, userId)

	_, err := r.db.Exec(query, args...)

	return err
}
