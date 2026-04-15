package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"notes/internal/database"
	"notes/internal/models"
	"strconv"
)

type Handler struct {
	store *database.Store
}

func NewHandler(store *database.Store) *Handler {
	return &Handler{
		store: store,
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting notes")

		return
	}

	respondWithJSON(w, http.StatusOK, notes)
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var input models.CreateNoteInput

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	note, err := h.store.Create(input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating note")

		return
	}

	respondWithJSON(w, http.StatusCreated, note)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	var input models.CreateNoteInput

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	note, err := h.store.Update(input, id)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Заметка не найдена")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error updating note")

		return
	}

	respondWithJSON(w, http.StatusOK, note)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	note, err := h.store.Delete(id)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Заметка не найдена")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error deleting note")

		return
	}

	respondWithJSON(w, http.StatusOK, note)
}

func (h *Handler) GetNoteById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	note, err := h.store.GetById(id)

	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "Заметка не найдена")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения заметки")

		return
	}

	respondWithJSON(w, http.StatusOK, note)
}
