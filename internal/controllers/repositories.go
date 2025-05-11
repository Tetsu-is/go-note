package controllers

import "go-note/domain"

type INoteRepository interface {
	CreateNote(*domain.Note) error
	ReadNoteAll() ([]*domain.Note, error)
	ReadNoteByID(id int) (*domain.Note, error)
	UpdateNote(*domain.Note) error
	DeleteNote(id int) error
}
