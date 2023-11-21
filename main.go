package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Book struct {
	id        int64
	create_at string
	title     string
	author    string
	price     float64
}

func main() {

	//Configure Database connection properties
	cfg := mysql.Config{
		User:   "root",
		Passwd: "1234",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "book_store",
	}
	fmt.Println("Conecting to database...")
	//Get a database handle
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if pingErr := db.Ping(); pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	fmt.Println("Plese write the author that you want to search in the catalog.")
	var author string
	var askError error
	author, askError = askAuthor()

	if askError != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	books, err := booksByAuthor(author)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Books found whose author name is %v: %d\n", author, len(books))
	for i, book := range books {
		fmt.Println(i, book.title)
	}

	b, _ := bookById(1)
	fmt.Println("The book with id 1 is:", b)

	b2 := Book{
		title:  "La Vuelta al Mundo en 80 Días",
		author: "Julio Verne",
		price:  30,
	}
	r, addError := addBook(b2)
	if addError != nil {
		log.Fatal(addError)
	}
	fmt.Printf("The book with id %d was added \n", r)
}

func askAuthor() (string, error) {
	var author string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		author = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("askAuthor %v", err)
	}
	return author, nil
}

func booksByAuthor(name string) ([]Book, error) {
	var books []Book

	rows, err := db.Query("SELECT * FROM books WHERE author = ?", name)
	if err != nil {
		return nil, fmt.Errorf("booksByAuthor %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.id, &book.create_at, &book.title, &book.author, &book.price); err != nil {
			return nil, fmt.Errorf("booksByAuthor %q: %v", name, err)
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("booksByAuthor %q: %v", name, err)
	}
	return books, nil
}

func bookById(id int64) (Book, error) {
	var b Book
	row := db.QueryRow("SELECT * FROM books WHERE id = ?", id)
	if err := row.Scan(&b.id, &b.create_at, &b.title, &b.author, &b.price); err != nil {
		if err == sql.ErrNoRows {
			return b, fmt.Errorf("bookById : %d - there is no such book", id)
		}
		return b, fmt.Errorf("booksByAuthor %q: %v", id, err)
	}

	return b, nil

}

func addBook(b Book) (int64, error) {
	r, err := db.Exec("INSERT INTO books (id, created_at, title, author, price) VALUES (?,?,?,?,?)", b.id, time.Now(), b.title, b.author, b.price)
	if err != nil {
		return 0, fmt.Errorf("addBook: %v", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addBook: %v", err)
	}
	return id, nil
}
