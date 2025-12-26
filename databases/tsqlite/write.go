package tsqlite

import (
	"database/sql"
	"fmt"
	"log"
)

func WriteJsonToSqlite(filePath string, table string, json string) {
	// Open database
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Insert books
	books := []struct {
		title  string
		author string
	}{
		{"The Go Programming Language", "Alan Donovan"},
		{"Clean Code", "Robert Martin"},
		{"Design Patterns", "Gang of Four"},
		{"The Pragmatic Programmer", "David Thomas"},
		{"Concurrency in Go", "Katherine Cox-Buday"},
	}

	for _, book := range books {
		result, err := db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", book.title, book.author)
		if err != nil {
			log.Fatal(err)
		}
		id, _ := result.LastInsertId()
		fmt.Printf("Inserted: %s by %s (ID: %d)\n", book.title, book.author, id)
	}

	// Query all books
	fmt.Println("\nAll books:")
	rows, err := db.Query("SELECT id, title, author FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, author string
		rows.Scan(&id, &title, &author)
		fmt.Printf("%d. %s by %s\n", id, title, author)
	}
}
