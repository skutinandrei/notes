package database

import (
	"database/sql"
	"errors"
	"fmt"
	"notes/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("note not found")

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAll() ([]models.Note, error) {
	notes := make([]models.Note, 0)

	query := `
		SELECT id, title, content, created_at
		FROM note
		ORDER BY created_at DESC
	`
	err := s.db.Select(&notes, query)

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении заметок: %w", err)
	}

	return notes, nil
}

func (s *Store) GetById(id int) (*models.Note, error) {
	var note models.Note

	query := `
		SELECT id, title, content, created_at
		FROM note
		WHERE id = ?
	`
	err := s.db.QueryRowx(query, id).StructScan(&note)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: заметка с id %d не найдена", ErrNotFound, id)
		}
		return nil, fmt.Errorf("ошибка при получении заметки: %w", err)
	}

	return &note, nil
}

func (s *Store) Create(input models.CreateNoteInput) (*models.Note, error) {
	var note models.Note

	query := `
		INSERT INTO note (title, content, created_at)
		VALUES (?, ?, ?)
		RETURNING id, title, content, created_at
	`
	now := time.Now().Format(time.RFC3339)

	err := s.db.QueryRowx(
		query,
		input.Title,
		input.Content,
		now,
	).StructScan(&note)

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("ошибка создания записи в дневнике: %w", err)
	}

	return &note, nil
}

func (s *Store) Update(input models.CreateNoteInput, id int) (*models.Note, error) {
	var note models.Note

	query := `
		UPDATE note
		SET title = ?, content = ?
		WHERE id = ?
		RETURNING id, title, content, created_at
	`

	err := s.db.QueryRowx(query, input.Title, input.Content, id).StructScan(&note)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: заметка с id %d не найдена", ErrNotFound, id)
		}
		return nil, fmt.Errorf("ошибка обновления: %w", err)
	}

	return &note, nil
}

func (s *Store) Delete(id int) (*models.Note, error) {
	var note models.Note

	query := `
		DELETE FROM note		
		WHERE id = ?
		RETURNING id, title, content, created_at
	`

	err := s.db.QueryRowx(query, id).StructScan(&note)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: заметка с id %d не найдена", ErrNotFound, id)
		}
		return nil, fmt.Errorf("ошибка удаления: %w", err)
	}

	return &note, nil
}
