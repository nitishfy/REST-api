package config 

import (
	"database/sql"
	"log"
	
	"fmt"

)

type Config struct {
	Password string
	Addr string
	DBName string
}

func (c *Config) OpenConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s","root",c.Password,c.Addr,c.DBName))
	if err != nil {
		log.Fatalf("error opening the sql connection: %v", err)
	}

	return db
}

func (c *Config) CloseConnection(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("error closing the connection: %v", err.Error())
	}
}

