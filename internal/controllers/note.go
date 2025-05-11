package controllers

import (
	"go-note/domain"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Note struct {
	ID      int
	Title   string
	Content string
}

type NoteController struct {
	NoteRepository INoteRepository
	Templates      *template.Template
}

func NewNoteController(nr INoteRepository, tmpl *template.Template) *NoteController {
	return &NoteController{
		NoteRepository: nr,
		Templates:      tmpl,
	}
}

func (nc *NoteController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ルーティング判定用にパスを正規化
	trimmed := strings.TrimPrefix(r.URL.Path, "/notes")
	trimmed = strings.Trim(trimmed, "/")

	// /notes または /notes/ の場合は一覧
	if trimmed == "" {
		notes, err := nc.NoteRepository.ReadNoteAll()
		if err != nil {
			http.Error(w, "Failed to get notes", http.StatusInternalServerError)
			return
		}
		nc.Templates.ExecuteTemplate(w, "notes_list.html", notes)
		return
	}

	// /notes/new
	if trimmed == "new" && r.Method == http.MethodGet {
		nc.Templates.ExecuteTemplate(w, "note_new.html", nil)
		return
	}
	if trimmed == "new" && r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		note := &domain.Note{Title: title, Content: content}
		if err := nc.NoteRepository.CreateNote(note); err != nil {
			http.Error(w, "Failed to create note", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	}

	// /notes/{id}/edit
	if strings.HasSuffix(trimmed, "edit") {
		idStr := strings.TrimSuffix(trimmed, "/edit")
		idStr = strings.TrimSuffix(idStr, "/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			note, err := nc.NoteRepository.ReadNoteByID(id)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			nc.Templates.ExecuteTemplate(w, "note_edit.html", note)
			return
		}
		if r.Method == http.MethodPost {
			title := r.FormValue("title")
			content := r.FormValue("content")
			note := &domain.Note{ID: id, Title: title, Content: content}
			if err := nc.NoteRepository.UpdateNote(note); err != nil {
				http.Error(w, "Failed to update note", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/notes/"+idStr, http.StatusSeeOther)
			return
		}
	}

	// /notes/{id}/delete
	if strings.HasSuffix(trimmed, "delete") && r.Method == http.MethodPost {
		idStr := strings.TrimSuffix(trimmed, "/delete")
		idStr = strings.TrimSuffix(idStr, "/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if err := nc.NoteRepository.DeleteNote(id); err != nil {
			http.Error(w, "Failed to delete note", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	}

	// /notes/{id}（ノート詳細）
	id, err := strconv.Atoi(trimmed)
	if err == nil {
		note, err := nc.NoteRepository.ReadNoteByID(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		nc.Templates.ExecuteTemplate(w, "note_detail.html", note)
		return
	}

	http.NotFound(w, r)
}
