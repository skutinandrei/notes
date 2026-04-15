package models

type Note struct {
	Id        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type CreateNoteInput struct {
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}
