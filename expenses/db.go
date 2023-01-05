package expenses

import (
	"database/sql"
	"log"
	"os"
)

func InitDB(db *sql.DB) *sql.DB {
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

func ConnectDB() *sql.DB {
	var err error
	db_url := os.Getenv("DATABASE_URL")
	db_url += "?ssl.mode=disable"
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	return db
}
