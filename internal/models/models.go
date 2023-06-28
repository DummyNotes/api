package models

type Note struct {
	NoteId    string `json:"NoteId"`
	UserId    string `json:"UserId"`
	Title     string `json:"Title"`
	Text      string `json:"Text"`
	CreatedAt string `json:"CreatedAt"`
	UpdatedAt string `json:"UpdatedAt"`
}

type User struct {
	UserID string
}
