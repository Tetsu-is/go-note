package controllers

import (
	"fmt"
	"net/http"
)

// application全体の親Controller

type ApplicationController struct {
	nc *NoteController
}

func NewApplicationController(noteController *NoteController) *ApplicationController {
	return &ApplicationController{
		nc: noteController,
	}
}

func (ac *ApplicationController) StartServer() {
	mux := http.NewServeMux()

	// ここにルーティングを追加していく
	mux.HandleFunc("/notes/", ac.nc.ServeHTTP)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", mux)
}
