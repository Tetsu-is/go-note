package repositories

import (
	"database/sql"
	"errors"
	"go-note/domain"
)

// INoteRepositoryの実装

type NoteRepository struct {
	DB *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{
		DB: db,
	}
}

func (nr *NoteRepository) CreateNote(note *domain.Note) error {
	if note.Title == "" {
		return errors.New("title cannot be empty")
	}
	query := "INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id"
	return nr.DB.QueryRow(query, note.Title, note.Content).Scan(&note.ID)
}

func (nr *NoteRepository) ReadNoteAll() ([]*domain.Note, error) {
	rows, err := nr.DB.Query("SELECT id, title, content FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []*domain.Note
	for rows.Next() {
		n := &domain.Note{}
		if err := rows.Scan(&n.ID, &n.Title, &n.Content); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (nr *NoteRepository) ReadNoteByID(id int) (*domain.Note, error) {
	n := &domain.Note{}
	err := nr.DB.QueryRow("SELECT id, title, content FROM notes WHERE id = $1", id).Scan(&n.ID, &n.Title, &n.Content)
	if err == sql.ErrNoRows {
		return nil, errors.New("note not found")
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (nr *NoteRepository) UpdateNote(note *domain.Note) error {
	res, err := nr.DB.Exec("UPDATE notes SET title = $1, content = $2 WHERE id = $3", note.Title, note.Content, note.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("note not found")
	}
	return nil
}

func (nr *NoteRepository) DeleteNote(id int) error {
	res, err := nr.DB.Exec("DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("note not found")
	}
	return nil
}
