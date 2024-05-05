package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nitishfy/REST-API/internal/config"
	"github.com/nitishfy/REST-API/internal/models"
)

type ConfigHandler struct {
	Config *config.Config
}

func (c *ConfigHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	db := c.Config.OpenConnection()
	defer c.Config.CloseConnection(db)

	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("error querying the books table %s\n", err.Error())
	}

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Name, &book.Isbn)
		if err != nil {
			log.Fatalf("error while scanning the row %s\n", err.Error())
			return
		}

		books = append(books, book)
	}

	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Fatalf("error encoding the books: %v", err)
		return
	}
}

func (c *ConfigHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	db := c.Config.OpenConnection()
	defer c.Config.CloseConnection(db)

	query := "select * from books where id = ?"

	var book models.Book
	w.Header().Add("Content-Type", "application/json")

	row := db.QueryRow(query, id)
	err := row.Scan(&book.ID, &book.Name, &book.Isbn)

	switch err {
	case sql.ErrNoRows:
		log.Println("No rows were returned")
		return
	case nil:
		json.NewEncoder(w).Encode(&book)
	default:
		panic(err)
	}
}

func (c *ConfigHandler) PostBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Fatalf("error decoding request body: %v", err.Error())
		return
	}

	db := c.Config.OpenConnection()
	defer c.Config.CloseConnection(db)

	query := "insert into books(id, name, isbn) values (?, ?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err.Error())
		return
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, book.ID, book.Name, book.Isbn)
	if err != nil {
		log.Printf("Error %s when inserting row into books table", err.Error())
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return
	}

	log.Printf("%d book created ", rows)
}

func (c *ConfigHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updatedBook models.Book
	err := json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		log.Fatalf("error decoding the request body: %v", err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")

	db := c.Config.OpenConnection()
	defer c.Config.CloseConnection(db)

	var exists bool
	err = db.QueryRowContext(context.Background(), "SELECT EXISTS(SELECT 1 from books WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		log.Printf("error checking if boook with ID %s exists: %v", id, err.Error())
		return
	}

	if !exists {
		log.Fatalf("Book not found")
		return
	}

	query := "UPDATE books SET name = ?, isbn = ? WHERE id = ?"
	_, err = db.ExecContext(context.Background(), query, updatedBook.Name, updatedBook.Isbn, id)
	if err != nil {
		log.Printf("error updating book with ID %s: %v", id, err.Error())
		return
	}

	log.Printf("Book with ID %s updated successfully", id)
}

func (c *ConfigHandler) DeleteBooks(w http.ResponseWriter, r *http.Request) {
	db := c.Config.OpenConnection()
	defer c.Config.CloseConnection(db)

	query := "DELETE FROM books"
	_, err := db.ExecContext(context.Background(), query)
	if err != nil {
		log.Fatalf("error deleting all the rows from the table: %v", err.Error())
	}

	log.Println("Rows deleted successfully!")
}

func (c *ConfigHandler) DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	w.Header().Add("Content-Type", "application/json")

	db := c.Config.OpenConnection()
	defer c.Config.CloseConnection(db)

	query := "DELETE from books where id = ?"
	res, err := db.ExecContext(context.Background(), query, id)
	if err != nil {
		log.Fatalf("error deleting book with the ID %s: %v", id, err)
		return
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return
	}

	if rowAffected == 0 {
		log.Printf("No book found with ID %s", id)
		return
	}

	log.Printf("Book with ID %s deleted successfully!", id)
}
