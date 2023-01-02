package expenses

import (
	"database/sql"
	"log"
	"os"
)

type handler struct {
	DB *sql.DB
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}
func InitDB() *sql.DB {
	db := connectDB()
	// Create table
	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`
	_, err := db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
	return db
}

func connectDB() (db *sql.DB) {
	var err error
	db_url := os.Getenv("DATABASE_URL")
	db_url += "?ssl.mode=disable"
	db, err = sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	return db
}
