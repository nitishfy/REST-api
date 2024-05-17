package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/nitishfy/REST-API/internal/config"
	"github.com/nitishfy/REST-API/internal/handlers"
)

var (
	API_PATH = "/apis/v1/books"
)

func main() {
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
		dbPass = "my-deault-password"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	router := mux.NewRouter()

	ch := handlers.ConfigHandler {
		Config: &config.Config{
			Password: dbPass,
			Addr:     dbHost,
			DBName:   dbName,
		},
	}

	router.HandleFunc(apiPath, ch.GetBooks).Methods("GET")
	router.HandleFunc(apiPath+"/{id}", ch.GetBookByID).Methods("GET")
	router.HandleFunc(apiPath, ch.PostBook).Methods("POST")
	router.HandleFunc(apiPath+"/{id}", ch.UpdateBook).Methods("PUT")
	router.HandleFunc(apiPath, ch.DeleteBooks).Methods("DELETE")
	router.HandleFunc(apiPath+"/{id}", ch.DeleteBookByID).Methods("DELETE")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("error while listening: %v", err)
		return
	}
}
