package main

import (
	"database/sql"
	"go-note/controllers"
	"go-note/repositories"
	"html/template"

	_ "github.com/lib/pq"
)

func main() {
	// connect to the database
	db, err := sql.Open("postgres", "user=devuser dbname=devdb password=devpassword sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	q := `
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

	// init application modules
	noteRepository := repositories.NewNoteRepository(db)

	noteController := controllers.NewNoteController(noteRepository, template.Must(template.ParseGlob("./assets/html/*.html")))

	appController := controllers.NewApplicationController(noteController)

	appController.StartServer()
	// start the server
}
