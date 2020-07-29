package mssql

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

// Golang MSSQL driver

type Client struct {
	db *sql.DB
}

func NewClient(connectionString string) (Client, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return Client{}, err
	}

	return Client{db: db}, nil
}
