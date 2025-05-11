package repositories_test

import (
	"database/sql"
	"go-note/domain"
	"go-note/repositories"

	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("postgres", "user=testuser password=testpassword dbname=testdb host=localhost port=5433 sslmode=disable")
	if err != nil {
		panic(err)
	}

	q := `
	DROP TABLE IF EXISTS notes;

	CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL
	);

	INSERT INTO notes (title, content) VALUES
		('Test Note 1', 'This is the content of test note 1'),
		('Test Note 2', 'This is the content of test note 2'),
		('Test Note 3', 'This is the content of test note 3');
	`

	_, err = db.Exec(q)
	if err != nil {
		panic(err)
	}

	return db
}

var db *sql.DB

func init() {
	db = setupTestDB()
}

func TestCreateNote(t *testing.T) {
	tests := []struct {
		name    string
		note    *domain.Note
		wantErr bool
	}{
		{"valid note", &domain.Note{Title: "Title", Content: "Content"}, false},
		{"empty title", &domain.Note{Title: "", Content: "Content"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositories.NewNoteRepository(db)
			err := repo.CreateNote(tt.note)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadNoteAll(t *testing.T) {
	t.Run("empty repo", func(t *testing.T) {
		repo := repositories.NewNoteRepository(db)
		notes, err := repo.ReadNoteAll()
		assert.NoError(t, err)
		assert.Len(t, notes, 3)
	})
	t.Run("with notes", func(t *testing.T) {
		repo := repositories.NewNoteRepository(db)
		_ = repo.CreateNote(&domain.Note{Title: "A", Content: "B"})
		notes, err := repo.ReadNoteAll()
		assert.NoError(t, err)
		assert.NotEmpty(t, notes)
	})
}

func TestReadNoteByID(t *testing.T) {
	repo := repositories.NewNoteRepository(db)
	note := &domain.Note{Title: "A", Content: "B"}
	_ = repo.CreateNote(note)
	t.Run("existing", func(t *testing.T) {
		got, err := repo.ReadNoteByID(1)
		assert.NoError(t, err)
		assert.Equal(t, "Test Note 1", got.Title)
	})
	t.Run("not found", func(t *testing.T) {
		_, err := repo.ReadNoteByID(999)
		assert.Error(t, err)
	})
}

func TestUpdateNote(t *testing.T) {
	repo := repositories.NewNoteRepository(db)
	note := &domain.Note{Title: "A", Content: "B"}
	_ = repo.CreateNote(note)
	t.Run("update existing", func(t *testing.T) {
		updated := &domain.Note{ID: 1, Title: "C", Content: "D"}
		err := repo.UpdateNote(updated)
		assert.NoError(t, err)
		got, _ := repo.ReadNoteByID(1)
		assert.Equal(t, "C", got.Title)
	})
	t.Run("update not found", func(t *testing.T) {
		updated := &domain.Note{ID: 999, Title: "X", Content: "Y"}
		err := repo.UpdateNote(updated)
		assert.Error(t, err)
	})
}

func TestDeleteNote(t *testing.T) {
	repo := repositories.NewNoteRepository(db)
	note := &domain.Note{Title: "A", Content: "B"}
	_ = repo.CreateNote(note)
	t.Run("delete existing", func(t *testing.T) {
		err := repo.DeleteNote(1)
		assert.NoError(t, err)
		_, err = repo.ReadNoteByID(1)
		assert.Error(t, err)
	})
	t.Run("delete not found", func(t *testing.T) {
		err := repo.DeleteNote(999)
		assert.Error(t, err)
	})
}
