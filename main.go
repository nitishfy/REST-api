package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)


type config struct {
	Password string
	Addr string
	DBName string
}

type Book struct {
	ID string
	Name string
	Isbn string
}

var (
	API_PATH = "/apis/v1/books"
	books = []Book{}
)

func main() {
	// host, db-name, password, apiPath
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "nitish"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	router := mux.NewRouter()
	c := config {
		Password: dbPass,
		Addr: dbHost,
		DBName: dbName,
	}

	router.HandleFunc(apiPath,c.getBooks).Methods("GET")
	router.HandleFunc(apiPath,c.postBook).Methods("POST")
	router.HandleFunc(apiPath,c.deleteBooks).Methods("DELETE")
	if err := http.ListenAndServe(":8081", router); err != nil  {
		log.Fatalf("error while listening: %v", err)
		return 
	}
}

func (c *config) getBooks(w http.ResponseWriter,  r *http.Request) {
	books = []Book{}
	// open the connection
	db := c.OpenConnection()
	// read the boooks
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("error querying the books table %s\n", err.Error())
	}
	
	
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("error while scanning the row %s\n", err.Error())
		}

		book := Book {
			ID: id,
			Name: name,
			Isbn: isbn,
		}

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)

	// close the connection
	c.CloseConnection(db)
}

func (c *config) deleteBooks(w http.ResponseWriter, r *http.Request) {
	// open the connection
	db := c.OpenConnection()
	// delete all the books
	query := "DELETE FROM books"
	_, err := db.ExecContext(context.Background(),query)
	if err != nil {
		log.Fatalf("error deleting all the rows from the table: %v", err.Error())
	}

	log.Println("Rows deleted successfully!")
	books = nil
	// close the connection
	c.CloseConnection(db)
}

func (c *config) postBook(w http.ResponseWriter, r *http.Request) {
	// read the request into an instance of a book
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)
	// open connection
	db := c.OpenConnection()
	// write the data
	query := "insert into books(id, name, isbn) values (?, ?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(),5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err.Error())
		return
	}
    defer stmt.Close()
	res, err := stmt.ExecContext(ctx, book.ID,book.Name,book.Isbn)
	if err != nil {
		log.Printf("Error %s when inserting row into books table", err.Error())
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return
	}
	log.Printf("%d books created ", rows)
	// close the connection
	c.CloseConnection(db)
}

func (c *config) OpenConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s","root",c.Password,c.Addr,c.DBName))
	if err != nil {
		log.Fatalf("error opening the sql connection: %v", err)
	}

	return db
}

func (c *config) CloseConnection(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("error closing the connection: %v", err.Error())
	}
}