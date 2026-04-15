package models

type Note struct {
	Id        int    `db:"id"`
	Title     string `db:"title"`
	Content   string `db:"content"`
	CreatedAt string `db:"created_at"`
}

type CreateNoteInput struct {
	Title   string `db:"title"`
	Content string `db:"content"`
}
