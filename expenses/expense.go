package expenses

import "database/sql"

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

type Handler struct {
	DB *sql.DB
}

func NewApplication(db *sql.DB) *Handler {
	return &Handler{db}
}
