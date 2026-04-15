package main

import (
	"log"
	"net/http"

	"notes/internal/database"
	"notes/internal/handlers"
	"os"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	defer database.CloseDB()

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	db := database.GetDB()
	note := database.NewStore(db)
	handlers := handlers.NewHandler(note)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /notes", handlers.GetAllNotes)
	mux.HandleFunc("POST /notes", handlers.CreateNote)
	mux.HandleFunc("PATCH /notes/{id}", handlers.UpdateNote)
	mux.HandleFunc("DELETE /notes/{id}", handlers.DeleteNote)
	mux.HandleFunc("GET /notes/{id}", handlers.GetNoteById)

	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, mux)

	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
